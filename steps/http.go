package steps

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/glopal/go-yp/yplib"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

var wrapValuesInSeq *yqlib.ExpressionNode

func init() {
	gob.Register(Http{})
	yplib.RegisterStep("http", NewHttp)
}

type Http struct {
	client      http.Client
	Method      string
	Url         yplib.Dval
	Headers     yplib.Dval
	PayloadForm yplib.Dval
	PayloadJson yplib.Dval
}

type httpOut struct {
	StatusCode int          `yaml:"statusCode"`
	Response   httpResponse `yaml:"response"`
}

type httpResponse struct {
	Body    string              `yaml:"body"`
	Headers map[string][]string `yaml:"headers"`
}

func NewHttp(tag string, node yplib.Node, ech yplib.ExecContextHooks) (yplib.Step, error) {
	if wrapValuesInSeq == nil {
		en, _ := yqlib.ExpressionParser.ParseExpression(`with((.[] | select(kind != "seq")); . = [.])`)
		wrapValuesInSeq = en
	}

	method := "GET"
	if tag != "" {
		method = strings.ToUpper(strings.ReplaceAll(tag, "!", ""))
	}

	return Http{
		client: http.Client{
			Timeout: time.Second * 10,
		},
		Method:      method,
		Url:         node.ValueResolver(".url"),
		Headers:     node.Resolver(".headers"),
		PayloadForm: node.Resolver(".payload.form"),
		PayloadJson: node.Resolver(".payload.json"),
	}, nil
}

func (h Http) Run(ion yplib.IoNode) (yplib.IoNode, error) {
	URL, err := getUrl(h.Url, ion)
	if err != nil {
		return nil, err
	}

	headers, err := getHeaders(h.Headers, ion)
	if err != nil {
		return nil, err
	}

	var payload io.Reader
	if h.Method == "POST" || h.Method == "PUT" {
		if reader, contentType, err := getPayload(h, ion); err == nil {
			payload = reader
			headers["Content-Type"] = []string{contentType}
		}
	}

	req, err := http.NewRequest(h.Method, URL, payload)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	body := ""
	if bodyText, err := io.ReadAll(resp.Body); err == nil {
		body = string(bodyText)
	}

	out := httpOut{
		StatusCode: resp.StatusCode,
		Response: httpResponse{
			Body:    string(body),
			Headers: resp.Header,
		},
	}

	outNode := &yaml.Node{}

	err = outNode.Encode(&out)
	if err != nil {
		return nil, err
	}

	return ion.Out(outNode), nil
}

func getUrl(urlResolver yplib.Dval, input yplib.IoNode) (string, error) {
	urlIo, err := urlResolver(input)
	if err != nil {
		return "", err
	}

	return urlIo.GetNode().Value, nil
}

func getHeaders(headersResolver yplib.Dval, input yplib.IoNode) (map[string][]string, error) {
	headers := map[string][]string{}
	headersIo, err := headersResolver(input)
	if err != nil {
		return headers, nil
	}

	return decodeNormalizedHeaderMap(headersIo)
}

func getPayload(h Http, input yplib.IoNode) (io.Reader, string, error) {
	if reader, err := getPayloadJson(h.PayloadJson, input); err == nil {
		return reader, "application/json", nil
	} else if reader, err := getPayloadJson(h.PayloadJson, input); err == nil {
		return reader, "application/x-www-form-urlencoded", nil
	}

	return nil, "", errors.New("payload not found")
}

func getPayloadForm(resolver yplib.Dval, input yplib.IoNode) (io.Reader, error) {
	formIo, err := resolver(input)
	if err != nil {
		return nil, err
	}

	form, err := decodeNormalizedHeaderMap(formIo)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(url.Values(form).Encode()), nil
}

func getPayloadJson(resolver yplib.Dval, input yplib.IoNode) (io.Reader, error) {
	jsonIo, err := resolver(input)
	if err != nil {
		return nil, err
	}

	fmt.Println(jsonIo.GetNode().Content[0].Value)

	jsonData, err := jsonIo.ToJson()
	if err != nil {
		return nil, err
	}

	fmt.Println("JSON: ", string(jsonData))
	return bytes.NewBuffer(jsonData), nil
}

func decodeNormalizedHeaderMap(input yplib.IoNode) (map[string][]string, error) {
	headers := map[string][]string{}
	headersIo := input.Yq(wrapValuesInSeq)

	err := headersIo.GetNode().Decode(&headers)

	return headers, err
}

package steps

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/glopal/go-ya/yalib"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

var wrapValuesInSeq *yqlib.ExpressionNode

func init() {
	gob.Register(Http{})
	yalib.RegisterStep("http", NewHttp)
}

type Http struct {
	client      http.Client
	Method      string
	Url         yalib.Dval
	Headers     yalib.Dval
	PayloadForm yalib.Dval
	PayloadJson yalib.Dval
}

type httpOut struct {
	StatusCode int          `yaml:"statusCode"`
	Response   httpResponse `yaml:"response"`
}

type httpResponse struct {
	Body    string              `yaml:"body"`
	Headers map[string][]string `yaml:"headers"`
}

func NewHttp(tag string, node yalib.Node, ech yalib.ExecContextHooks) (yalib.Step, error) {
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

func (h Http) Run(ion yalib.IoNode) (yalib.IoNode, error) {
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

func getUrl(urlResolver yalib.Dval, input yalib.IoNode) (string, error) {
	urlIo, err := urlResolver(input)
	if err != nil {
		return "", err
	}

	return urlIo.GetNode().Value, nil
}

func getHeaders(headersResolver yalib.Dval, input yalib.IoNode) (map[string][]string, error) {
	headers := map[string][]string{}
	headersIo, err := headersResolver(input)
	if err != nil {
		return headers, nil
	}

	return decodeNormalizedHeaderMap(headersIo)
}

func getPayload(h Http, input yalib.IoNode) (io.Reader, string, error) {
	if reader, err := getPayloadJson(h.PayloadJson, input); err == nil {
		return reader, "application/json", nil
	} else if reader, err := getPayloadForm(h.PayloadForm, input); err == nil {
		return reader, "application/x-www-form-urlencoded", nil
	}

	return nil, "", errors.New("payload not found")
}

func getPayloadForm(resolver yalib.Dval, input yalib.IoNode) (io.Reader, error) {
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

func getPayloadJson(resolver yalib.Dval, input yalib.IoNode) (io.Reader, error) {
	jsonIo, err := resolver(input)
	if err != nil {
		return nil, err
	}

	jsonData, err := jsonIo.ToJson()
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}

func decodeNormalizedHeaderMap(input yalib.IoNode) (map[string][]string, error) {
	headers := map[string][]string{}
	headersIo := input.Yq(wrapValuesInSeq)

	err := headersIo.GetNode().Decode(&headers)

	return headers, err
}

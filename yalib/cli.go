package yalib

type Cli struct {
	Cli Command `yaml:"cli"`
}
type Command struct {
	Use   string `yaml:"use"`
	Short string `yaml:"short"`
}

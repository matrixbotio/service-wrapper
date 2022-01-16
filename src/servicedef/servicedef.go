package servicedef

import (
	"os"

	"gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type ProcessDef struct {
	ServiceName string `yaml:"service_name"`
	LinesToPreserve int `yaml:"lines_to_preserve"`
	SeparateStdoutStderr bool `yaml:"separate_stdout_stderr"`
	Command string `yaml:"command"`
	Args []string `yaml:"args"`
}

func GetProcessDefinition() ProcessDef {
	pdef := ProcessDef{}
	yamlData, err := os.ReadFile(os.Args[1])
	check(err)
	err = yaml.Unmarshal([]byte(yamlData), &pdef)
	check(err)
	return pdef
}

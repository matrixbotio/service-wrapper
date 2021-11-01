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

type processDef struct {
	ServiceName string `yaml:"service_name"`
	LinesToPreserve int `yaml:"lines_to_preserve"`
	SeparateStdoutStderr bool `yaml:"separate_stdout_stderr"`
	Command string `yaml:"command"`
	Args []string `yaml:"args"`
}

func getProcessDefinition() processDef {
	pdef := processDef{}
	yamlData, err := os.ReadFile(os.Args[1])
	check(err)
	err = yaml.Unmarshal([]byte(yamlData), &pdef)
	check(err)
	return pdef
}

var def = getProcessDefinition()

// exported
var ServiceName = def.ServiceName
var LinesToPreserve = def.LinesToPreserve
var SeparateStdoutStderr = def.SeparateStdoutStderr
var Command = def.Command
var Args = def.Args

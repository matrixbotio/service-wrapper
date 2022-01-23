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

func getEnvName(str string) string {
	closeIdx := 0
	for i := 0; i < len(str); i++ {
		if str[i:i + 1] == "}" {
			closeIdx = i
			break
		}
	}
	if closeIdx == 0 {
		return ""
	}
	return str[:closeIdx]
}

func getArg(rawArg string) string {
	res := ""
	rawArgLen := len(rawArg)
	if rawArgLen < 4 {
		return rawArg
	}
	i := 0
	for ; i < rawArgLen - 4; i++ {
		if rawArg[i:i + 2] == "${" {
			name := getEnvName(rawArg[i + 2:])
			res += os.Getenv(name)
			i += len(name) + 2
		} else {
			res += string(rawArg[i])
		}
	}
	return res + rawArg[i:rawArgLen]
}

func GetProcessDefinition() ProcessDef {
	pdef := ProcessDef{}
	yamlData, err := os.ReadFile(os.Args[1])
	check(err)
	err = yaml.Unmarshal([]byte(yamlData), &pdef)
	check(err)
	for i, v := range pdef.Args {
		pdef.Args[i] = getArg(v)
	}
	return pdef
}

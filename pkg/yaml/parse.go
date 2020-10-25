package yaml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigFile struct {
	MasterN MasterNodes `yaml:"master-nodes"`
	WorkerN WorkerNodes `yaml:"worker-nodes"`
}

type MasterNodes struct {
	H []Hosts `yaml:"hosts"`
}

type WorkerNodes struct {
	W []Hosts `yaml:"hosts"`
}

type Config struct {
	CD ConfigData `yaml:"configuration"`
}

type ConfigData struct {
	Username string `yaml:"username"`
}

type Hosts struct {
	Hostname  string `yaml:"hostname"`
	IPAddress string `yaml:"ip-address"`
}

func ParseMain(filePath string) (infrastructure *ConfigFile, err error) {
	data := readFile(filePath)

	err = yaml.Unmarshal(data, &infrastructure)
	if err != nil {
		return
	}

	return
}

func readFile(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil
	}

	return data
}

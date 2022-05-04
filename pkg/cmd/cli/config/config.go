package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MusicPath string `yaml:"musicPath"`
	DBPath    string `yaml:"dbPath"`
}

func LoadConfig(confPath string) (Config, error) {
	conf := Config{}
	// rename current dir
	usr, _ := user.Current()
	if strings.HasPrefix(confPath, "~") {
		confPath = strings.Replace(confPath, "~", usr.HomeDir, 1)
	}

	// check existance dir & file
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		if err := os.MkdirAll(confPath, 0755); err != nil {
			return conf, err
		}
	}
	if _, err := os.Stat(confPath + "/conf.yaml"); os.IsNotExist(err) {
		CreateDefaultFile(confPath)
	}

	// load conf
	bytes, err := ioutil.ReadFile(confPath + "/conf.yaml")
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal([]byte(bytes), &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

func CreateDefaultFile(confPath string) error {
	fp, err := os.Create(confPath + "/conf.yaml")
	if err != nil {
		return err
	}
	defer fp.Close()

	// create template if need
	raw := `musicPath: .
dbPath: ~/.punos/punos.db`

	fp.WriteString(raw)
	return nil
}

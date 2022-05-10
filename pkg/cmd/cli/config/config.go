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

func LoadConfig(confPath string) (*Config, error) {
	conf := Config{}
	// rename current dir
	usr, _ := user.Current()
	if strings.HasPrefix(confPath, "~") {
		confPath = strings.Replace(confPath, "~", usr.HomeDir, 1)
	}

	// check existance dir & files
	_, dirErr := os.Stat(confPath)
	_, fileErr := os.Stat(confPath + "/conf.yaml")

	if os.IsNotExist(dirErr) {
		if err := os.MkdirAll(confPath, 0755); err != nil {
			return &conf, err
		}
	}
	if os.IsNotExist(fileErr) {
		CreateDefaultFile(confPath)
	}

	// load conf
	bytes, err := ioutil.ReadFile(confPath + "/conf.yaml")
	if err != nil {
		return &conf, err
	}
	err = yaml.Unmarshal([]byte(bytes), &conf)
	if err != nil {
		return &conf, err
	}

	// update path starting from HOME dir
	if strings.HasPrefix(conf.DBPath, "~") {
		usr, _ := user.Current()
		conf.DBPath = strings.Replace(conf.DBPath, "~", usr.HomeDir, 1)
	}

	return &conf, nil
}

func CreateDefaultFile(confPath string) error {
	fp, err := os.Create(confPath + "/conf.yaml")
	if err != nil {
		return err
	}
	defer fp.Close()

	// create template if need
	raw := `dbPath: ~/.punos/punos.db`

	fp.WriteString(raw)
	return nil
}

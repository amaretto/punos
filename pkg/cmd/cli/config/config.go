package config

import (
	"fmt"
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

	if os.IsNotExist(dirErr) || os.IsNotExist(fileErr) {
		fmt.Printf("Thank you for playing punos! \n")
		fmt.Printf("It seems that there is no conf/db files. Is it ok to create these files in %s? [y/N]:", confPath)
		var choise string
		for {
			fmt.Scanf("%s", &choise)
			if choise == "y" {
				break
			} else if choise == "N" {
				os.Exit(0)
			} else {
				fmt.Printf("Is it ok to create : %s/conf.yaml?[y/N]\n", confPath)
			}
		}
	}

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

	//ToDo: create db if it doesn't exist
	return &conf, nil
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

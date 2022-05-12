package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/fatih/color"
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
	confFilePath := confPath + "/conf.yaml"

	// check existance dir & files
	_, dirErr := os.Stat(confPath)
	_, fileErr := os.Stat(confFilePath)

	if os.IsNotExist(dirErr) {
		if err := os.MkdirAll(confPath, 0755); err != nil {
			return &conf, err
		}
	}
	if os.IsNotExist(fileErr) {
		CreateDefaultFile(confPath)
	}

	// load conf
	bytes, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return &conf, err
	}
	err = yaml.Unmarshal([]byte(bytes), &conf)
	if err != nil {
		return &conf, err
	}

	conf.updateHomePath()

	// If there is no MusicPath setting, ask user to set current directory as MusicPath
	if conf.MusicPath == "" {
		cd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		fmt.Printf("Thank you for playing with punos!\n\n")
		fmt.Printf("It seems you set music directory for punos yet. \nIs it ok to set the below current dir[%s] as music path?[y/N]:\n", color.GreenString(cd))

		var choise string
		fmt.Scanf("%s", &choise)
		for {
			if choise == "y" {
				conf.MusicPath = cd
				f, err := os.OpenFile(confFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					return nil, err
				}
				enc := yaml.NewEncoder(f)
				enc.Encode(conf)
				break
			} else if choise == "N" {
				fmt.Printf("OK. Please move to your music path or write your music path to %s\n", confPath+"/conf.yaml")
				os.Exit(0)
			} else {
				fmt.Printf("Is it ok to set the below current dir[%s] as music path?[y/N]:\n", cd)
				fmt.Scanf("%s", &choise)
			}
		}
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

// update path starting from HOME dir
func (conf *Config) updateHomePath() {
	usr, _ := user.Current()
	if strings.HasPrefix(conf.MusicPath, "~") {
		conf.DBPath = strings.Replace(conf.MusicPath, "~", usr.HomeDir, 1)
	}
	if strings.HasPrefix(conf.DBPath, "~") {
		conf.DBPath = strings.Replace(conf.DBPath, "~", usr.HomeDir, 1)
	}
}

package config

import (
	"os"
	"os/user"
	"strings"
	"testing"
)

var configTestCases = []struct {
	Path    string
	Cleanup bool
	Result  error
}{
	{"~/.punos", false, nil}, // use homedir, file exist
	{"~/.test", true, nil},   // use homedir, file doesn't exist
	{".current", true, nil},  // use homedir, file doesn't exist
}

func cleanUp(confPath string) error {
	usr, _ := user.Current()
	if strings.HasPrefix(confPath, "~") {
		confPath = strings.Replace(confPath, "~", usr.HomeDir, 1)
	}
	if err := os.RemoveAll(confPath); err != nil {
		return err
	}
	return nil
}

func TestLoadConfig(t *testing.T) {
	// ToDo: create files for testing
	for _, testCase := range configTestCases {
		result, err := loadConfig(testCase.Path)
		if err != testCase.Result {
			t.Errorf("invalid result. testCase:%v, actual:%v", testCase, result)
		}
		if testCase.Cleanup {
			err := cleanUp(testCase.Path)
			if err != nil {
				t.Errorf("cleanup failed with configPath(%s)", testCase.Path)
			}
		}
	}
}

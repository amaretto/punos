package player

import "testing"

var configTestCases = []struct {
	Path   string
	Result error
}{
	{"~/.punos/conf", nil},  // use homedir, file exist
	{"~/.punos/conf2", nil}, // use homedir, file doesn't exist
	// ToDo: fix the below paths
	//{"/.punos/conf", nil},  // use absolute path, file exist
	//{"/.punos/conf", nil},  // use absolute path, file doesn't exist
}

func TestLoadConfig(t *testing.T) {
	// ToDo: create files for testing
	for _, testCase := range configTestCases {
		result := loadConfig(testCase.Path)
		if result != testCase.Result {
			t.Errorf("invalid result. testCase:%v, actual:%v", testCase, result)
		}
	}
}

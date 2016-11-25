package cli

import (
	"io/ioutil"
	"strings"
)

func contentFromFileOrCLI(s string) (string, error) {
	if strings.HasPrefix(s, "file://") {
		file := strings.TrimPrefix(s, "file://")
		bs, err := ioutil.ReadFile(file)
		return strings.Trim(string(bs), "\n"), err
	}
	return s, nil
}

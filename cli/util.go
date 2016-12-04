package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/mgo.v2/bson"
)

func contentFromFileOrCLI(s string) (string, error) {
	if strings.HasPrefix(s, "file://") {
		file := strings.TrimPrefix(s, "file://")
		bs, err := ioutil.ReadFile(file)
		return strings.Trim(string(bs), "\n"), err
	}
	return s, nil
}

func bsonID(c *cli.Context) (bson.ObjectId, error) {
	id := c.String("id")
	if !bson.IsObjectIdHex(id) {
		return "", fmt.Errorf("invalid bson id %s", id)
	}
	return bson.ObjectIdHex(id), nil
}

func pretty(data interface{}) error {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(append(b, '\n'))
	return err
}

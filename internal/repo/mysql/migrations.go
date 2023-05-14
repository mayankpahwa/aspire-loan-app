package mysql

import (
	"io/ioutil"
	"strings"
)

func LoadSQLFile(path string) error {
	// Read file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	for _, q := range strings.Split(string(file), ";") {
		q := strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := GetConnection().Exec(q); err != nil {
			return err
		}
	}
	return nil
}

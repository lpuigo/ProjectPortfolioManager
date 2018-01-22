package Config

import (
	"encoding/json"
	"os"
)

type ConfigStruct interface {
}

// SetFromFile loads conf values from given file if exists, or create file with given Conf (default) values if not exists
func SetFromFile(file string, conf ConfigStruct) error {
	// test if file exists
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			err := saveConfigFile(file, conf)
			if err != nil {
				return err
			}
			return nil
		}
	}
	defer f.Close()
	// file exist, lets load it and replace conf values
	err = json.NewDecoder(f).Decode(conf)
	if err != nil {
		return err
	}

	return nil
}

func saveConfigFile(file string, conf ConfigStruct) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	je := json.NewEncoder(f)
	je.SetIndent("", "\t")
	return je.Encode(conf)
}

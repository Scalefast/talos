package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func ParseConfig(config string) (*viper.Viper, error) {
	v := viper.New()
	format, err := CheckFile(config)
	if err != nil {
		return v, err
	}
	v.SetConfigType(format)
	v.SetConfigFile(config)
	return v, v.ReadInConfig()
}

func CheckFile(config string) (viperFormat string, err error) {
	fileContent, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return "", err
	}

	var checkfile interface{}
	err = json.Unmarshal(fileContent, &checkfile)
	if err == nil {
		return "json", nil
	}
	err = yaml.Unmarshal(fileContent, &checkfile)
	if err == nil {
		return "yaml", nil
	}
	return "", fmt.Errorf("The file is not in JSON or YAML format")
}

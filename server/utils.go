package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func LoadRouterConfig(fileName string) map[string]map[string]string {

	var config map[string]map[string]string
	mockPath := GetMockPath()
	jsonFile := GetFile(mockPath, fileName)
	if err := json.Unmarshal(jsonFile, &config); err != nil {
		panic(err)
	}

	return config
}

func GetFile(pathFile string, filename string) []byte {

	//filename := path.Join(path.Join(os.Getenv("PWD"), "resource"), "layout.html.mustache")
	filename = path.Join(pathFile, filename)
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	return file
}

func GetMockPath() string {
	path, _ := os.Getwd()
	path = fmt.Sprintf("%s%s", path, "/configuration/router")
	return path
}

func GetMapValueFromJsonRawMessage[T any](properties []byte, key string) (any, error) {
	j := make(map[string]any)
	err := json.Unmarshal(properties, &j)
	return j[key], err
}

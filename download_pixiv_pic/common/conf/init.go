package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var ConfigData map[string]map[string]interface{}

func init() {

	data, err := os.ReadFile("./config.yaml")
	cwd, _ := os.Getwd()
	fmt.Printf("Reading %s/config.yaml\n", cwd)
	if err != nil {
		log.Panicf("当前目录 %s 未找到 config.yaml", cwd)
	}
	m := make(map[string]map[string]interface{})
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatal(err)
	}
	ConfigData = m
	fmt.Println("Configuration reload.")

}

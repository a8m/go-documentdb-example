package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Url       string `json:"url"`
	MasterKey string `json:"masterKey"`
}

func main() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	config := &Config{}
	if err = json.Unmarshal(data, config); err != nil {
		panic(err)
	}
	fmt.Println(config)
}

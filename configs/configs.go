package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type AppConfigs struct {
	App App `json:"app"`
}

type MongoConfigs struct {
	Mongo Mongo `json:"mongo"`
}

type App struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

type Mongo struct {
	URI         string   `json:"uri"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Database    string   `json:"db_name"`
	Collections []string `json:"collections"`
}

func ReadConfigs() (*App, *Mongo, error) {
	// read file
	data, err := ioutil.ReadFile("./config.json")

	if err != nil {
		return nil, nil, err
	}

	var appConfigs AppConfigs
	var mongoConfigs MongoConfigs

	// unmarshaling
	err = json.Unmarshal(data, &appConfigs)
	err = json.Unmarshal(data, &mongoConfigs)

	if err != nil {
		return nil, nil, err
	}
	fmt.Println(appConfigs.App)
	fmt.Println(mongoConfigs.Mongo)

	return &appConfigs.App, &mongoConfigs.Mongo, nil
}

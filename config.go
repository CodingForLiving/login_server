package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	RedisAddr string `json:"RedisAddr"`
	RedisAuth string `json:"RedisAuth"`
	RedisIndex int   `json:"RedisIndex"`
	MysqlStr string  `json:"MysqlStr"`
	HttpAddr string  `json:"HttpAddr"`
}

var config *Config = &Config {
	RedisAddr: "127.0.0.1:6379",
	RedisAuth: "auth",
	RedisIndex: 1,
	MysqlStr: "root:system@tcp(127.0.0.1:3306)/game",
	HttpAddr: ":80",
}

func LoadConfig(path string){
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic("read config file err:" + err.Error())
		return
	}

	err = json.Unmarshal(bytes, config)
	if err != nil {
		log.Panic("parse config file err:"+err.Error())
		return
	}

	log.Println("load config success")
}
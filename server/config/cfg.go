package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	e "taskFive/server/lib/err"
)

type Config struct {
	Database struct {
		Addr       string `yaml:"host"`
		Port       string `yaml:"port"`
		Username   string `yaml:"user"`
		Password   string `yaml:"pass"`
		DBname     string `yaml:"dbname"`
		DriverName string `yaml:"driverName"`
	} `yaml:"database"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	Nats struct {
		ClusterID string `yaml:"cluster_id"`
		ClientID  string `yaml:"client_id"`
		Url       string `yaml:"url"`
	} `yaml:"nats"`
	CHouse struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"cHouse"`
}

func (cfg *Config) InitFile() {
	f, err := os.Open("config/config.yml")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer func() { err = e.WrapIfErr("can`t init config-file", err) }()
}

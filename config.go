package main

import (
	"github.com/BurntSushi/toml"
)

type Sql struct {
	Host          string `toml:"host"`
	User          string `toml:"user"`
	Pass          string `toml:"pass"`
	DBName        string `toml:"database_name"`
	DBType        string `toml:"database_type"`
	ConnectionMax int    `toml:"connection_max"`
	IsEnabled     bool   `toml:"enabled"`
}

type Config struct {
	UsersAmount int    `toml:"users_amount"`
	StartVkId   int    `toml:"start_vk_id"`
	UsersPass   string `toml:"users_password"`
	Sql         Sql    `toml:"database"`
}

func GetConfig() (*Config, error) {
	var conf = &Config{}
	_, err := toml.DecodeFile("config.toml", conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

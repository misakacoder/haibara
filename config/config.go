package config

import (
	"github.com/jinzhu/configor"
	"github.com/misakacoder/logger"
)

type configuration struct {
	Log struct {
		Filename string
		Level    string
	}
	Server struct {
		Port int
	}
	Database struct {
		Host            string
		Port            int
		Username        string
		Password        string
		Name            string
		MaxIdleConn     int    `yaml:"maxIdleConn"`
		MaxOpenConn     int    `yaml:"maxOpenConn"`
		ConnMaxLifeTime string `yaml:"connMaxLifeTime"`
		SlowSqlTime     string `yaml:"slowSqlTime"`
		PrintSql        bool   `yaml:"printSql"`
	}
}

var Configuration = configuration{}

func init() {
	err := configor.Load(&Configuration, "haibara.yml")
	if err != nil {
		logger.Panic(err.Error())
	}
}

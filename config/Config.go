package config

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	AppName     string         `json:"app_name"`
	AppMode     string         `json:"app_mode"`
	AppHost     string         `json:"app_host"`
	AppPort     string         `json:"app_port"`
	Sms         SmsConfig      `json:"sms"`
	Database    DatabaseConfig `json:"database"`
	RedisConfig RedisConfig    `json:"redis"`
}

type SmsConfig struct {
	SignName     string `json:"sign_name"`
	TemplateCode string `json:"template_code"`
	AppKey       string `json:"app_key"`
	AppSecret    string `json:"app_secret"`
	EndPoint     string `json:"endpoint"`
}

type DatabaseConfig struct {
	Driver    string `json:"driver,omitempty"`
	User      string `json:"user,omitempty"`
	Password  string `json:"password,omitempty"`
	Host      string `json:"host,omitempty"`
	Port      string `json:"port,omitempty"`
	DbName    string `json:"db_name,omitempty"`
	Charset   string `json:"charset,omitempty"`
	ParseTime string `json:"parseTime"`
	Loc       string `json:"loc"`
}

type RedisConfig struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
	Db       int    `json:"db,omitempty"`
}

var cfg *Config

func GetConfig() *Config {
	return cfg
}

// ParseConfig 解析app.json 中的文件
func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取json文件，利用方法将json文件中信息，映射到结构体上
	read := bufio.NewReader(file)
	decoder := json.NewDecoder(read)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

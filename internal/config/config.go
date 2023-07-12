package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Dbname   string `yaml:"dbname"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"postgres"`
	Grpc struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"grpc"`
}

func GetConfig() (Config, error) {
	c := Config{}

	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return c, errors.New("cant read config")
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, errors.New("bad config")
	}

	dbname, exists := os.LookupEnv("POSTGRES_DB")
	if exists {
		c.Postgres.Dbname = dbname
	}

	user, exists := os.LookupEnv("POSTGRES_USER")
	if exists {
		c.Postgres.User = user
	}

	password, exists := os.LookupEnv("POSTGRES_PASSWORD")
	if exists {
		c.Postgres.Password = password
	}

	if c.Grpc.Port == "" || c.Grpc.Host == "" {
		return c, errors.New("missing grpc data in config")
	}

	if c.Postgres.User == "" || c.Postgres.Password == "" {
		return c, errors.New("missing postgres login or password in config")
	}

	if c.Postgres.Dbname == "" {
		return c, errors.New("missing postgres dbname in config")
	}

	if c.Postgres.Host == "" || c.Postgres.Port == "" {
		return c, errors.New("missing postgres host or port in config")
	}

	return c, nil
}

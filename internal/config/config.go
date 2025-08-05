package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	AppEnv string `yaml:"app_env" default:"development"`
	Path   string `yaml:"dav_path" required:"true"`
	Prefix string `yaml:"dav_prefix" default:"/"`
	Port   int    `yaml:"port" default:"5656"`
	TLS    bool   `yaml:"tls" default:"false"`
	User   User   `yaml:"user"`
	Argon  Argon  `yaml:"argon"`
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Argon struct {
	ArgonMemory      uint32 `yaml:"argon_memory" env-default:"1024"`
	ArgonIterations  uint32 `yaml:"argon_iterations" env-default:"3"`
	ArgonParallelism uint8  `yaml:"argon_parallelism" env-default:"2"`
	ArgonSaltLength  uint32 `yaml:"argon_salt_length" env-default:"16"`
	ArgonKeyLength   uint32 `yaml:"argon_key_length" env-default:"32"`
}

func MustLoad(path string) Config {
	cfg := Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", path)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to process .env: " + err.Error())
	}

	return cfg
}

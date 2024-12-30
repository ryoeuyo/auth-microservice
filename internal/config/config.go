package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env          string       `yaml:"environment" env-required:"true"`
	GRPCServer   GRPCServer   `yaml:"grpc_server"`
	MetricServer MetricServer `yaml:"metric_server"`
	Database     Database     `yaml:"database"`
	JWTSecretKey string       `env:"JWT_SECRET_KEY,required"`
}

type GRPCServer struct {
	Address     string        `yaml:"address"`
	Port        uint16        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	TokenTTL    time.Duration `yaml:"token_ttl"`
}

type MetricServer struct {
	Port    uint16 `yaml:"port"`
	Address string `yaml:"address"`
}

type Database struct {
	Engine       string `yaml:"engine"`
	Host         string `yaml:"host"`
	Port         uint16 `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MigrationDir string `yaml:"migration_dir"`
}

func MustLoad(path ...string) *AppConfig {
	godotenv.Load()

	var configPath string

	if path != nil && len(path) > 0 {
		configPath = path[0]
		if configPath == "" {
			panic("config file path is empty")
		}
	} else {
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			panic("CONFIG_PATH is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists")
	}

	var cfg AppConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config")
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read env")
	}

	return &cfg
}

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Postgresql struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Url      string `mapstructure:"url"`
}

func (p *Postgresql) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.Database,
	)
}

type Config struct {
	Test       string      `mapstructure:"test"`
	AppId      string      `mapstructure:"appId"`
	Mode       string      `mapstructure:"mode"`
	Bust       int64       `mapstructure:"bust"`
	Env        string      `mapstructure:"env"`
	Postgresql *Postgresql `mapstructure:"postgresql"`
	JwtKey     string      `mapstructure:"jwtKey"`
	Port       string      `mapstructure:"port"`
	TestMode   bool        `mapstructure:"testMode"`
}

func (c *Config) Valid() error {
	if c.Mode == "" {
		return errors.New("mode is empty")
	}

	// check port
	if c.Port == "" {
		return errors.New("port is not set")
	}

	return nil
}

func findGitRoot() (string, error) {
	// Start from the current directory.
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, ".gitignore")); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", errors.New("no .gitignore found in any parent directory")
}

var instance *Config
var once sync.Once

// New creates and returns a new Config instance
func New(mode string, absConfigPath string) *Config {
	once.Do(func() {
		if mode == "" {
			mode = os.Getenv("Mode")
			if mode == "" {
				mode = "local"
			}
		}

		instance = &Config{}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		var configPath string
		baseDir, err := findGitRoot()

		if err != nil {
			fmt.Printf("Failed to find Git root: %v", err)
			os.Exit(1)
		}

		if absConfigPath == "" {
			switch mode {
			case "test":
				configPath = "config/test/"
			case "pre":
				configPath = "config/pre/"
			case "pro":
				configPath = "config/pro/"
			default:
				configPath = "config/local/"
			}
		} else {
			configPath = absConfigPath
		}

		// fmt.Println("✅✅✅✅✅configPath: ", configPath)

		viper.AddConfigPath(filepath.Join(baseDir, configPath))
		viper.AddConfigPath(filepath.Join("/app/", configPath))

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file: %s\n", err)
			os.Exit(1)
		}

		viper.SetEnvPrefix("CONFIG_ENV")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv() // todo: check why this is not working

		viper.BindEnv("jwtKey", "CONFIG_ENV_JWT_KEY")
		viper.BindEnv("postgresql.host", "CONFIG_ENV_PG_HOST")
		viper.BindEnv("postgresql.port", "CONFIG_ENV_PG_PORT")
		viper.BindEnv("postgresql.user", "CONFIG_ENV_PG_USER")
		viper.BindEnv("postgresql.password", "CONFIG_ENV_PG_PASSWORD")
		viper.BindEnv("postgresql.database", "CONFIG_ENV_PG_DATABASE")
		viper.BindEnv("postgresql.url", "CONFIG_ENV_PG_URL")

		if err := viper.Unmarshal(&instance); err != nil {
			fmt.Printf("Error unmarshalling config: %s\n", err)
			os.Exit(1)
		}

		if instance == nil {
			fmt.Println("Config instance is nil after unmarshalling")
			os.Exit(1)
		}

	})
	return instance
}

package config

import (
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	SSLMode         string        `yaml:"ssl-mode"`
	DDLMode         string        `yaml:"ddl-mode"`
	Reset           bool          `yaml:"reset"`
	GCInterval      time.Duration `yaml:"GCInterval"`
	MaxIdleConns    int           `yaml:"max-idle-conns"`
	MaxOpenConns    int           `yaml:"max-open-conns"`
	ConnMaxLifetime time.Duration `yaml:"conn-max-lifetime"`
}

type AppSettings struct {
	Name      string `yaml:"name"`
	Env       string `yaml:"env"`
	Port      string `yaml:"port"`
	LogLevel  string `yaml:"log-level"`
	JWTSecret string `yaml:"jwt-secret"`
}

type CorsOriginConfig struct {
	Origin       string `yaml:"origin"`
	AllowMethods string `yaml:"allowMethods"`
	AllowHeaders string `yaml:"allowHeaders"`
}

type CorsConfig struct {
	Origins []CorsOriginConfig `yaml:"origins"`
}

type Config struct {
	Database DBConfig    `yaml:"database"`
	App      AppSettings `yaml:"app"`
	Cors     CorsConfig  `yaml:"cors"`
}

var GlobalConfig Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		return err
	}

	if GlobalConfig.Database.Host == "" || GlobalConfig.App.Port == "" {
		log.Fatal("Missing required configuration")
	}

	return nil
}

func GetAllowOrigins() string {
	var origins []string
	for _, origin := range GlobalConfig.Cors.Origins {
		if origin.Origin != "" {
			origins = append(origins, origin.Origin)
		}
	}
	unique := splitAndUnique(origins)
	return strings.Join(unique, ",")
}

func uniqueStrings(input []string) []string {
	m := make(map[string]struct{})
	var result []string
	for _, str := range input {
		if _, exists := m[str]; !exists {
			m[str] = struct{}{}
			result = append(result, str)
		}
	}
	return result
}

func splitAndUnique(values []string) []string {
	var all []string
	for _, v := range values {
		parts := strings.Split(v, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				all = append(all, p)
			}
		}
	}
	return uniqueStrings(all)
}

func GetAllowMethods() string {
	var methods []string
	for _, origin := range GlobalConfig.Cors.Origins {
		if origin.AllowMethods != "" {
			methods = append(methods, origin.AllowMethods)
		}
	}
	unique := splitAndUnique(methods)
	return strings.Join(unique, ",")
}

func GetAllowHeaders() string {
	var headers []string
	for _, origin := range GlobalConfig.Cors.Origins {
		if origin.AllowHeaders != "" {
			headers = append(headers, origin.AllowHeaders)
		}
	}
	unique := splitAndUnique(headers)
	return strings.Join(unique, ",")
}

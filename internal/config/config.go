package config

import "os"

type Config struct {
	Port    string
	APIPath string
}

func LoadConfig() Config {
	return Config{
		Port:    getEnv("PORT", "8080"),
		APIPath: getEnv("API_PATH", "/"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

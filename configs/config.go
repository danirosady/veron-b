package configs

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Env     string
	Port    string
	Host    string
	Origins []string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}

	cfg.App.Env = viper.GetString("APP_ENV")
	cfg.App.Port = viper.GetString("APP_PORT")
	cfg.App.Host = viper.GetString("APP_HOST")

	origins := viper.GetString("APP_ALLOWED_ORIGINS")
	if origins != "" {
		cfg.App.Origins = strings.Split(origins, ",")
	} else {
		cfg.App.Origins = []string{"http://localhost:5173"}
	}

	cfg.Database.Host = viper.GetString("DB_HOST")
	cfg.Database.Port = viper.GetString("DB_PORT")
	cfg.Database.Name = viper.GetString("DB_NAME")
	cfg.Database.User = viper.GetString("DB_USER")
	cfg.Database.Password = viper.GetString("DB_PASSWORD")
	cfg.Database.SSLMode = viper.GetString("DB_SSLMODE")
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	cfg.JWT.Secret = viper.GetString("JWT_SECRET")
	cfg.JWT.AccessExpiry = getDuration("JWT_ACCESS_EXPIRY", 15*time.Minute)
	cfg.JWT.RefreshExpiry = getDuration("JWT_REFRESH_EXPIRY", 168*time.Hour)

	cfg.Redis.Host = viper.GetString("REDIS_HOST")
	cfg.Redis.Port = viper.GetString("REDIS_PORT")
	cfg.Redis.Password = viper.GetString("REDIS_PASSWORD")
	cfg.Redis.DB = viper.GetInt("REDIS_DB")

	return cfg, nil
}

func LoadFromEnv() *Config {
	cfg := &Config{}

	cfg.App.Env = getEnv("APP_ENV", "development")
	cfg.App.Port = getEnv("APP_PORT", "8080")
	cfg.App.Host = getEnv("APP_HOST", "0.0.0.0")

	origins := getEnv("APP_ALLOWED_ORIGINS", "http://localhost:5173")
	cfg.App.Origins = strings.Split(origins, ",")

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.Name = getEnv("DB_NAME", "tms")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	cfg.JWT.Secret = getEnv("JWT_SECRET", "default-secret-change-me")
	cfg.JWT.AccessExpiry = getEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute)
	cfg.JWT.RefreshExpiry = getEnvDuration("JWT_REFRESH_EXPIRY", 168*time.Hour)

	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvInt("REDIS_DB", 0)

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	return getEnvDuration(key, fallback)
}

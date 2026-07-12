package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
	ServerPort         string   `mapstructure:"PORT"`
	MongoURI           string   `mapstructure:"MONGO_URI"`
	DBName             string   `mapstructure:"DB_NAME"`
	JWTSecretKey       string   `mapstructure:"JWT_SECRET_KEY"`
	JWTExpirationHours int      `mapstructure:"JWT_EXPIRATION_HOURS"`
	EnableCache        bool     `mapstructure:"ENABLE_CACHE"`
	RedisAddr          string   `mapstructure:"REDIS_ADDR"`
	RedisPassword      string   `mapstructure:"REDIS_PASSWORD"`
	LogLevel           string   `mapstructure:"LOG_LEVEL"`
	LogFormat          string   `mapstructure:"LOG_FORMAT"`
	CookieDomains      []string `mapstructure:"COOKIE_DOMAINS"`
	SecureCookie       bool     `mapstructure:"SECURE_COOKIE"`
	AllowedOrigins     []string `mapstructure:"ALLOWED_ORIGINS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENABLE_CACHE", false)
	viper.SetDefault("JWT_EXPIRATION_HOURS", 72)
	viper.SetDefault("COOKIE_DOMAINS", []string{"localhost"})
	viper.SetDefault("SECURE_COOKIE", false)
	viper.SetDefault("ALLOWED_ORIGINS", []string{"http://localhost:5173"})

	// Explicitly bind every remaining field to its environment variable.
	// viper.AutomaticEnv() alone does NOT make Unmarshal() pick up env vars
	// that were never registered via SetDefault/BindEnv/config file - without
	// this, MONGO_URI, DB_NAME, JWT_SECRET_KEY, REDIS_ADDR, REDIS_PASSWORD,
	// LOG_LEVEL, and LOG_FORMAT would silently unmarshal as empty strings
	// even when set correctly in the environment.
	_ = viper.BindEnv("MONGO_URI")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("JWT_SECRET_KEY")
	_ = viper.BindEnv("REDIS_ADDR")
	_ = viper.BindEnv("REDIS_PASSWORD")
	_ = viper.BindEnv("LOG_LEVEL")
	_ = viper.BindEnv("LOG_FORMAT")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	// Manually handle comma-separated strings for slices if viper didn't split them
	if allowedOrigins := viper.GetString("ALLOWED_ORIGINS"); allowedOrigins != "" {
		parts := strings.Split(allowedOrigins, ",")
		var cleaned []string
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			trimmed = strings.Trim(trimmed, "\"'")
			if trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		config.AllowedOrigins = cleaned
	}

	if cookieDomains := viper.GetString("COOKIE_DOMAINS"); cookieDomains != "" {
		parts := strings.Split(cookieDomains, ",")
		var cleaned []string
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			trimmed = strings.Trim(trimmed, "\"'")
			if trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		config.CookieDomains = cleaned
	}

	return
}

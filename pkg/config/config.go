package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Search   SearchConfig
	S3       S3Config
	Auth     AuthConfig
	SMTP     SMTPConfig
	Features FeaturesConfig
}

type ServerConfig struct {
	Domain string `mapstructure:"DOMAIN"`
	Port   int    `mapstructure:"PORT"`
	Env    string `mapstructure:"ENV"`
}

type DatabaseConfig struct {
	URL      string `mapstructure:"DATABASE_URL"`
	MaxConns int32  `mapstructure:"DATABASE_MAX_CONNS"`
	MinConns int32  `mapstructure:"DATABASE_MIN_CONNS"`
}

type RedisConfig struct {
	URL        string `mapstructure:"REDIS_URL"`
	MaxRetries int    `mapstructure:"REDIS_MAX_RETRIES"`
}

type SearchConfig struct {
	URL string `mapstructure:"MEILISEARCH_URL"`
	Key string `mapstructure:"MEILISEARCH_KEY"`
}

type S3Config struct {
	Endpoint  string `mapstructure:"S3_ENDPOINT"`
	AccessKey string `mapstructure:"S3_ACCESS_KEY"`
	SecretKey string `mapstructure:"S3_SECRET_KEY"`
	Bucket    string `mapstructure:"S3_BUCKET"`
	Region    string `mapstructure:"S3_REGION"`
	CDNBase   string `mapstructure:"CDN_BASE_URL"`
}

type AuthConfig struct {
	JWTSecret       string        `mapstructure:"JWT_SECRET"`
	AccessTokenTTL  time.Duration `mapstructure:"JWT_ACCESS_TTL"`
	RefreshTokenTTL time.Duration `mapstructure:"JWT_REFRESH_TTL"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"SMTP_HOST"`
	Port     int    `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
	From     string `mapstructure:"FROM_EMAIL"`
}

type FeaturesConfig struct {
	RegistrationOpen     bool `mapstructure:"REGISTRATION_OPEN"`
	RegistrationApproval bool `mapstructure:"REGISTRATION_APPROVAL"`
	FederationEnabled    bool `mapstructure:"FEDERATION_ENABLED"`
	MaxPostLength        int  `mapstructure:"MAX_POST_LENGTH"`
	MaxImageMB           int  `mapstructure:"MAX_IMAGE_MB"`
	MaxVideoMB           int  `mapstructure:"MAX_VIDEO_MB"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("ENV", "development")
	viper.SetDefault("DATABASE_MAX_CONNS", 20)
	viper.SetDefault("DATABASE_MIN_CONNS", 5)
	viper.SetDefault("REDIS_URL", "redis://localhost:6379")
	viper.SetDefault("REDIS_MAX_RETRIES", 3)
	viper.SetDefault("MEILISEARCH_URL", "http://localhost:7700")
	viper.SetDefault("S3_REGION", "us-east-1")
	viper.SetDefault("JWT_ACCESS_TTL", "15m")
	viper.SetDefault("JWT_REFRESH_TTL", "720h")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("REGISTRATION_OPEN", true)
	viper.SetDefault("REGISTRATION_APPROVAL", false)
	viper.SetDefault("FEDERATION_ENABLED", false)
	viper.SetDefault("MAX_POST_LENGTH", 500)
	viper.SetDefault("MAX_IMAGE_MB", 10)
	viper.SetDefault("MAX_VIDEO_MB", 100)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}

	cfg := &Config{}

	cfg.Server = ServerConfig{
		Domain: viper.GetString("DOMAIN"),
		Port:   viper.GetInt("PORT"),
		Env:    viper.GetString("ENV"),
	}

	cfg.Database = DatabaseConfig{
		URL:      viper.GetString("DATABASE_URL"),
		MaxConns: viper.GetInt32("DATABASE_MAX_CONNS"),
		MinConns: viper.GetInt32("DATABASE_MIN_CONNS"),
	}

	cfg.Redis = RedisConfig{
		URL:        viper.GetString("REDIS_URL"),
		MaxRetries: viper.GetInt("REDIS_MAX_RETRIES"),
	}

	cfg.Search = SearchConfig{
		URL: viper.GetString("MEILISEARCH_URL"),
		Key: viper.GetString("MEILISEARCH_KEY"),
	}

	cfg.S3 = S3Config{
		Endpoint:  viper.GetString("S3_ENDPOINT"),
		AccessKey: viper.GetString("S3_ACCESS_KEY"),
		SecretKey: viper.GetString("S3_SECRET_KEY"),
		Bucket:    viper.GetString("S3_BUCKET"),
		Region:    viper.GetString("S3_REGION"),
		CDNBase:   viper.GetString("CDN_BASE_URL"),
	}

	accessTTL, err := time.ParseDuration(viper.GetString("JWT_ACCESS_TTL"))
	if err != nil {
		accessTTL = 15 * time.Minute
	}
	refreshTTL, err := time.ParseDuration(viper.GetString("JWT_REFRESH_TTL"))
	if err != nil {
		refreshTTL = 720 * time.Hour
	}

	cfg.Auth = AuthConfig{
		JWTSecret:       viper.GetString("JWT_SECRET"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}

	cfg.SMTP = SMTPConfig{
		Host:     viper.GetString("SMTP_HOST"),
		Port:     viper.GetInt("SMTP_PORT"),
		Username: viper.GetString("SMTP_USERNAME"),
		Password: viper.GetString("SMTP_PASSWORD"),
		From:     viper.GetString("FROM_EMAIL"),
	}

	cfg.Features = FeaturesConfig{
		RegistrationOpen:     viper.GetBool("REGISTRATION_OPEN"),
		RegistrationApproval: viper.GetBool("REGISTRATION_APPROVAL"),
		FederationEnabled:    viper.GetBool("FEDERATION_ENABLED"),
		MaxPostLength:        viper.GetInt("MAX_POST_LENGTH"),
		MaxImageMB:           viper.GetInt("MAX_IMAGE_MB"),
		MaxVideoMB:           viper.GetInt("MAX_VIDEO_MB"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Server.Domain == "" {
		return fmt.Errorf("DOMAIN is required")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

func (c *Config) BaseURL() string {
	if c.IsDevelopment() {
		return fmt.Sprintf("http://localhost:%d", c.Server.Port)
	}
	return fmt.Sprintf("https://%s", c.Server.Domain)
}

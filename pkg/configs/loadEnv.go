package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	IsDebug    bool   `mapstructure:"IS_DEBUG"`
	AppName    string `mapstructure:"APP_NAME"`
	DB_DSN     string `mapstructure:"DB_DSN"`
	ServerPort string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
	RedisUri     string `mapstructure:"REDIS_URL"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPUser  string `mapstructure:"SMTP_USER"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

type SiteData struct {
	IsDebug    bool   `mapstructure:"IS_DEBUG"`
	AppEnv     string `mapstructure:"APP_ENV"`
	AppName    string `mapstructure:"APP_NAME"`
	AppKeyword string `mapstructure:"APP_KEYWORD"`
	AppAddress string `mapstructure:"APP_ADDRESS"`

	ServerPort string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
}

func GetSiteData(path string) (config SiteData, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

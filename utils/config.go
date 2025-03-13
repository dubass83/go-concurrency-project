package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config store all configuration of the application
// the values read by viper from file or enviroment variables
type Config struct {
	Enviroment              string        `mapstructure:"ENVIROMENT"`
	DBSource                string        `mapstructure:"DB_SOURCE"`
	DBPoolMaxConns          int32         `mapstructure:"DB_POOL_MAX_CONNS"`
	DBPoolMinConns          int32         `mapstructure:"DB_POOL_MIN_CONNS"`
	DBPoolMaxConnLifetime   time.Duration `mapstructure:"DB_POOL_MAX_CONN_LIFETIME"`
	DBPoolMaxConnIdleTime   time.Duration `mapstructure:"DB_POOL_MAX_CONN_IDLE_TIME"`
	DBPoolHealthCheckPeriod time.Duration `mapstructure:"DB_POOL_HEALTH_CHECK_PERIOD"`
	DBPoolConnectTimeout    time.Duration `mapstructure:"DB_POOL_CONNECT_TIMEOUT"`
	MigrationURL            string        `mapstructure:"MIGRATION_URL"`
	WebPort                 string        `mapstructure:"WEB_PORT"`
	RedisURL                string        `mapstructure:"REDIS_URL"`
	PathToTemplate          string        `mapstructure:"PATH_TO_TEMPLATE"`
	EmailTemplate           string        `mapstructure:"EMAIL_TEMPLATE"`
	EmailService            string        `mapstructure:"EMAIL_SERVICE"`
	EmailLogin              string        `mapstructure:"EMAIL_LOGIN"`
	EmailPassword           string        `mapstructure:"EMAIL_PASSWORD"`
	EmailEncryption         string        `mapstructure:"EMAIL_ENCRYPTION"`
	SenderName              string        `mapstructure:"SENDER_NAME"`
	SenderEmail             string        `mapstructure:"SENDER_EMAIL"`
}

// LoadConfig
func LoadConfig(configPath string) (config Config, err error) {
	v := viper.New()
	v.SetConfigName("conf")
	v.SetConfigType("env")
	v.AddConfigPath(configPath)
	err = v.ReadInConfig()
	if err != nil {
		return
	}
	v.AutomaticEnv()
	err = v.Unmarshal(&config)
	return
}

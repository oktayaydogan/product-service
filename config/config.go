package config

import (
	"github.com/spf13/viper"
)

// Config, uygulama konfigürasyonunu temsil eder
type Config struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

var AppConfig *Config

// LoadConfig, konfigürasyonu yükler
func LoadConfig() error {
	viper.SetConfigName("config.development")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Çevre değişkenlerini kullanarak konfigürasyonu oku
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Konfigürasyonu AppConfig'a yükle
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	return nil
}

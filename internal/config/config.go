package config

import (
    "log"

    "github.com/spf13/viper"
)

type Config struct {
    App struct {
        Name   string
        Port   string
        Mode   string
        APIKey string
    }
    Database struct {
        DSN string
    }
    Providers struct {
        Adyen struct {
            APIKey      string
            MerchantID  string
            Environment string
        }
        Stripe struct {
            APIKey    string
            SecretKey string
        }
    }
}

func Load() (*Config, error) {
    // 設定預設值
    viper.SetDefault("app.name", "Payment Gateway")
    viper.SetDefault("app.port", "8080")
    viper.SetDefault("app.mode", "debug")
    viper.SetDefault("app.api_key", "dev-api-key-12345")
    viper.SetDefault("database.dsn", "host=localhost user=postgres password=postgres dbname=payment_gateway port=5432 sslmode=disable")

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    
    viper.AddConfigPath("./config")
    viper.AddConfigPath(".")

    // 環境變數優先
    viper.AutomaticEnv()

    // 讀取配置檔
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Println("⚠️  Config file not found, using defaults")
        } else {
            return nil, err
        }
    } else {
        log.Printf("✅ Loaded config from: %s", viper.ConfigFileUsed())
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
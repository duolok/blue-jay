package config

import (
    "time"

    "github.com/spf13/viper"
)

type Config struct {
    UserAgent      string
    RequestTimeout time.Duration
    GameListURL    string
    CSVFileName    string
}

func LoadConfig() (*Config, error) {
    viper.SetDefault("USER_AGENT", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
    viper.SetDefault("REQUEST_TIMEOUT", 10*time.Second)
    viper.SetDefault("CSV_FILE_NAME", "games.csv")

    viper.AutomaticEnv()

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    config := &Config{
        UserAgent:      viper.GetString("USER_AGENT"),
        RequestTimeout: viper.GetDuration("REQUEST_TIMEOUT"),
        GameListURL:    viper.GetString("GAME_LIST_URL"),
        CSVFileName:    viper.GetString("CSV_FILE_NAME"),
    }

    return config, nil
}


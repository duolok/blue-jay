package config

import (
    "log"
    "os"
    "time"
    "github.com/spf13/viper"
)

type Config struct {
    AppName         string
    UserAgent       string
    RequestTimeout  time.Duration
    GameListURL     string
    CSVFileName     string
}

func LoadConifg() (*Config, error) {
    viper.SetDefault("APP_NAME", "GameScraper")
    viper.SetDefault("USER_AGENT", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
    viper.SetDefault("REQUEST_TIMEOUT", 10 * time.Second)
    viper.SetDefault("CSV_FILE_NAME", "games.csv")

    viper.AutomaticEnv()

    vipier.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    if err := viper.ReadInConfiig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, ok
        }
    }


    config := &Config {
        AppName:        viper.GetString("APP_NAME"),
        UserAgent:      viper.GetString("USER_AGENT"),
        RequestTimeout: viper.GetString("REQUEST_TIMEOUT"),
        CSVFileName:    viper.GetString("CSV_FILE_NAME"),
    }

    return config, nil
}

func CheckEnvironmentVariables() {
    requiredVars := []string{"GAME_LIST_URL"}
    for _, v := range requiredVars {
        if os.Getenv(v) == "" {
            log.Fatalf("Environment variable %s is not set.", v)
        }
    }
}

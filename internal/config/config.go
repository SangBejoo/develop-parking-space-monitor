package config

import (
    "github.com/spf13/viper"
    "log"
)

type Config struct {
    Tile38Host string `mapstructure:"TILE38_HOST"`
    Tile38Port int    `mapstructure:"TILE38_PORT"`
    RedisHost  string `mapstructure:"REDIS_HOST"`
    RedisPort  int    `mapstructure:"REDIS_PORT"`
    DBHost     string `mapstructure:"DB_HOST"`
    DBPort     int    `mapstructure:"DB_PORT"`
    DBUser     string `mapstructure:"DB_USER"`
    DBPassword string `mapstructure:"DB_PASSWORD"`
    DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig() (*Config, error) {
    viper.AutomaticEnv()

    // Set defaults
    viper.SetDefault("TILE38_HOST", "localhost")
    viper.SetDefault("TILE38_PORT", 9851)
    viper.SetDefault("REDIS_HOST", "localhost")
    viper.SetDefault("REDIS_PORT", 6379)
    viper.SetDefault("DB_HOST", "localhost")
    viper.SetDefault("DB_PORT", 5432)
    viper.SetDefault("DB_USER", "postgres")
    viper.SetDefault("DB_PASSWORD", "postgres")  // Add default password
    viper.SetDefault("DB_NAME", "parking_monitor")

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    // Add debug logging
    log.Printf("Database Config - Host: %s, Port: %d, User: %s, Password: %s, DBName: %s",
        config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

    return &config, nil
}

// ...existing code...
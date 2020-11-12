// This package reads configuration file.
package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/FTi130/keep-the-moment-app/back/lib/mail"
	"github.com/FTi130/keep-the-moment-app/back/lib/redis"
	"github.com/FTi130/keep-the-moment-app/back/postgres"
	"github.com/FTi130/keep-the-moment-app/back/server"
)

type (
	// Flags structure contains variables, which should be received as commandline flags.
	Flags struct {
		Path *string
	}
	// Config structure contains configurable options of the whole program.
	Config struct {
		Server   server.Config
		Postgres postgres.Config
		Redis    redis.Config
		Email    mail.Config
	}
)

// setDefaults sets default values for the configuration file.
func setDefaults() {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 5000)
	viper.SetDefault("server.secret", "") // Must always be reassigned!
	viper.SetDefault("server.domains", []string{"https://keepthemoment.ru", "https://www.keepthemoment.ru"})

	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.username", "postgres")
	viper.SetDefault("postgres.password", "")
	viper.SetDefault("postgres.database", "postgres")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")

	viper.SetDefault("email.host", "localhost")
	viper.SetDefault("email.port", 465)
	viper.SetDefault("email.email", "")
	viper.SetDefault("email.password", "")
	viper.SetDefault("email.service", "")
}

// Read reads configuration file from disk.
func Read(f Flags) (c Config) {
	setDefaults()

	viper.SetConfigName("config")
	viper.AddConfigPath(*f.Path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	c.Server.Host = viper.GetString("server.host")
	c.Server.Port = viper.GetInt("server.port")
	c.Server.Secret = viper.GetString("server.secret")
	c.Server.Domains = viper.GetStringSlice("server.domains")

	c.Postgres.Host = viper.GetString("postgres.host")
	c.Postgres.Port = viper.GetInt("postgres.port")
	c.Postgres.Username = viper.GetString("postgres.username")
	c.Postgres.Password = viper.GetString("postgres.password")
	c.Postgres.Database = viper.GetString("postgres.database")

	c.Redis.Host = viper.GetString("redis.host")
	c.Redis.Port = viper.GetInt("redis.port")
	c.Redis.Password = viper.GetString("redis.password")

	c.Email.Host = viper.GetString("email.host")
	c.Email.Port = viper.GetInt("email.port")
	c.Email.From = viper.GetString("email.email")
	c.Email.Password = viper.GetString("email.password")

	fmt.Printf("⇨ configuration loaded from %s\n", *f.Path)
	return c
}

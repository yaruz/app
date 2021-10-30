package config

import (
	"flag"

	yaruz_config "github.com/yaruz/app/pkg/yarus_platform/config"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/log"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

const (
	LangEng = "eng"
	LangRus = "rus"
)

// Configuration is the struct for app configuration
type Configuration struct {
	Server struct {
		HTTPListen string
	}
	Log           log.Config
	DB            DB
	YaruzMetadata yaruz_config.Metadata
	// JWT signing key. required.
	JWTSigningKey string
	// JWT expiration in hours. Defaults to 72 hours (3 days)
	JWTExpiration   uint
	SessionLifeTime uint
	CacheLifeTime   uint
}

type DB struct {
	Identity  minipkg_gorm.Config
	Reference minipkg_gorm.Config
	Data      minipkg_gorm.Config
	Search    minipkg_gorm.Config
	Redis     redis.Config
}

func (c *Configuration) YaruzConfig() yaruz_config.Configuration {
	return yaruz_config.Configuration{
		Infrastructure: &yaruz_config.Infrastructure{
			Log:           c.Log,
			ReferenceDB:   c.DB.Reference,
			DataDB:        c.DB.Data,
			SearchDB:      c.DB.Search,
			Redis:         c.DB.Redis,
			CacheLifeTime: c.CacheLifeTime,
		},
		Metadata: &c.YaruzMetadata,
	}
}

// defaultPathToConfig is the default path to the app config
const defaultPathToConfig = "config/config.yaml"
const defaultPathToMetadata = "config/metadata.yaml"

// Get func return the app config
func Get() (*Configuration, error) {
	// config is the app config
	var config Configuration = Configuration{}
	// pathToConfig is a path to the app config
	var pathToConfig string
	var pathToMetadata string

	flag.StringVar(&pathToConfig, "config", defaultPathToConfig, "path to YAML/JSON config file")
	flag.StringVar(&pathToMetadata, "metadata", defaultPathToMetadata, "path to YAML/JSON metadata file")
	flag.Parse()

	viper.AutomaticEnv() // read in environment variables that match

	if err := config.readConfig(pathToConfig); err != nil {
		return &config, err
	}

	if err := config.readMetadata(pathToMetadata); err != nil {
		return &config, err
	}

	return &config, nil
}

func (c *Configuration) readConfig(pathToConfig string) error {
	viper.SetConfigFile(pathToConfig)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.Errorf("Config file not found in %q", pathToConfig)
		} else {
			return errors.Errorf("Config file was found in %q, but was produced error: %v", pathToConfig, err)
		}
	}

	err := viper.Unmarshal(c)
	if err != nil {
		return errors.Errorf("Config unmarshal error: %v", err)
	}
	return nil
}

func (c *Configuration) readMetadata(pathToMetadata string) error {
	viper.SetConfigFile(pathToMetadata)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.Errorf("Metadata file not found in %q", pathToMetadata)
		} else {
			return errors.Errorf("Metadata file was found in %q, but was produced error: %v", pathToMetadata, err)
		}
	}

	err := viper.Unmarshal(&c.YaruzMetadata)
	if err != nil {
		return errors.Errorf("Metadata unmarshal error: %v", err)
	}
	return nil
}

func addition4Test(cfg *Configuration, logAppPostfix string) {
	cfg.Log.OutputPaths = []string{
		"stdout",
	}
	cfg.Log.InitialFields = map[string]interface{}{"app": "carizza-test: " + logAppPostfix}
	cfg.Log.Level = "debug"
	return
}

func Get4Test(logAppPostfix string) (*Configuration, error) {
	cfg, err := Get()
	if err != nil {
		return nil, err
	}
	addition4Test(cfg, logAppPostfix)

	return cfg, nil
}

func Get4UnitTest(logAppPostfix string) *Configuration {
	cfg := &Configuration{
		Log: log.Config{
			Encoding: "json",
		},
		DB: DB{
			Identity: minipkg_gorm.Config{
				Dialect:       "postgres",
				DSN:           "host=localhost port=5401 dbname=postgres user=postgres password=postgres sslmode=disable",
				IsAutoMigrate: true,
				Log: minipkg_gorm.LogConfig{
					LogLevel: 4,
				},
			},
			Redis: redis.Config{},
		},
		JWTSigningKey:   "test",
		JWTExpiration:   1,
		SessionLifeTime: 1,
		CacheLifeTime:   1,
	}
	addition4Test(cfg, logAppPostfix)

	return cfg
}

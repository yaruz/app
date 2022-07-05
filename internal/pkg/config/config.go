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

	// defaultPathToConfig is the default path to the app config
	defaultPathToConfig   = "config/config.yaml"
	defaultPathToMetadata = "metadata/metadata.yaml"
)

// Configuration is the struct for app configuration
type Configuration struct {
	Server struct {
		HTTPListen string
	}
	Log           log.Config
	DB            DB
	Auth          Auth
	YaruzMetadata yaruz_config.Metadata
	CacheLifeTime uint
}

type DB struct {
	Identity     minipkg_gorm.Config
	Reference    minipkg_gorm.Config
	DataSharding yaruz_config.Sharding
	Search       minipkg_gorm.Config
	Redis        redis.Config
}

type Auth struct {
	Endpoint               string
	ClientId               string
	ClientSecret           string
	Organization           string
	Application            string
	SignInRedirectURL      string
	JWTSigningKey          string
	JWTExpiration          uint
	SessionlifeTime        uint
	DefaultAccountSettings DefaultAccountSettings
}

type DefaultAccountSettings struct {
	Lang string
}

func (c *Configuration) YaruzConfig() *yaruz_config.Configuration {
	return &yaruz_config.Configuration{
		Infrastructure: &yaruz_config.Infrastructure{
			Log:           c.Log,
			ReferenceDB:   c.DB.Reference,
			DataSharding:  c.DB.DataSharding,
			SearchDB:      c.DB.Search,
			Redis:         c.DB.Redis,
			CacheLifeTime: c.CacheLifeTime,
		},
		Metadata: &c.YaruzMetadata,
	}
}

// Get func return the app config
func Get() (*Configuration, error) {
	// config is the app config
	var config Configuration = Configuration{}
	// pathToConfig is a path to the app config
	var pathToConfig string
	var pathToMetadata string

	viper.AutomaticEnv() // read in environment variables that match
	//viper.BindEnv("pathToConfig")
	defPathToConfig := defaultPathToConfig
	if viper.Get("pathToConfig") != nil {
		defPathToConfig = viper.Get("pathToConfig").(string)
	}

	flag.StringVar(&pathToConfig, "config", defPathToConfig, "path to YAML/JSON config file")
	flag.StringVar(&pathToMetadata, "metadata", defaultPathToMetadata, "path to YAML/JSON metadata file")
	flag.Parse()

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
		CacheLifeTime: 1,
	}
	addition4Test(cfg, logAppPostfix)

	return cfg
}

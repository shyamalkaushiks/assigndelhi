package config

import (
	"errors"
	"runtime"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ConfigStruct :
type ConfigStruct struct {
	SERVICE_PORT string
	DATABASE     string
	DB_SERVER    string
	DB_PORT      string
	DB_USER      string
	DB_PASSWORD  string
	DB_DATABASE  string
	API_SECRET   string
}

var showVersion = false

// const defaultConfigFile string = "/etc/nova_us_jobs/trainers/config/trainers.toml"
// const windowFile string = ""

const defaultConfigFile string = "F:\\novaultraaura\\novajobs-ultra-aura-shyamal2\\trainers\\config\\trainer.toml"
const windowFile string = "F:\\novaultraaura\\novajobs-ultra-aura-shyamal2\\trainers\\config\\trainer.toml"

// Config :
var Config ConfigStruct
var configFile string

// LoadConfig :
func LoadConfig() error {
	if runtime.GOOS == "windows" {
		flag.StringVarP(&configFile, "config", "c", windowFile, "assets configuration file.")
	} else {
		flag.StringVarP(&configFile, "config", "c", defaultConfigFile, "assets configuration file.")
	}
	flag.BoolVarP(&showVersion, "version", "v", false, "Display version information and exit")
	flag.Parse()

	viper.SetConfigFile(configFile)
	//loadDefaultConfig()

	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			// Log.Error().Err(err).Msg("error in loading config file, viper.ConfigParseError")
			return errors.New("error in loading config file, viper.ConfigParseError")
		} else {
			// Log.Error().Err(err).Msg("error in loading config file")
			return errors.New("error in loading config file")
		}
	}

	viper.Unmarshal(&Config)
	return nil
}

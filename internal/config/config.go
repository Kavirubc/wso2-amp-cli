package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	KeyAPIURL       = "api_url"
	KeyAPIKeyHeader = "api_key_header"
	KeyAPIKeyValue  = "api_key"
	KeyDefaultOrg   = "default_org"
	KeyDefaultProj  = "default_project"
)

//ConfigDir returns the path to .amp
func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".amp")
}

//ConfigFile returns the path to the config file
func ConfigFile() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func Init() error {
	if err:= os.MkdirAll(ConfigDir(), 0755); err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigDir())  // Look in ~/.amp/
	// Set default values
	viper.SetDefault(KeyAPIURL, "http://localhost:8080")
	viper.SetDefault(KeyAPIKeyHeader, "X-API-Key")
	viper.SetDefault(KeyAPIKeyValue, "")
	viper.SetDefault(KeyDefaultOrg, "")
	viper.SetDefault(KeyDefaultProj, "")
	// Try to read existing config (ignore error if file doesn't exist yet)
	_ = viper.ReadInConfig()
	return nil
}

func Get(key string) string {
	return viper.GetString(key)
}

func Set(key, value string) error {
	viper.Set(key, value)

	if err := os.MkdirAll(ConfigDir(), 0755); err != nil {
		return err
	}

	return viper.WriteConfigAs(ConfigFile())
}

func GetAPIURL() string       { return viper.GetString(KeyAPIURL) }
func GetAPIKeyHeader() string { return viper.GetString(KeyAPIKeyHeader) }
func GetAPIKeyValue() string  { return viper.GetString(KeyAPIKeyValue) }
func GetDefaultOrg() string   { return viper.GetString(KeyDefaultOrg) }
func GetDefaultProject() string { return viper.GetString(KeyDefaultProj) }
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type (
	baseConfiguration struct {
		// The Alphabill home directory
		HomeDir string
		// Configuration file URL. If it's relative, then it's relative from the HomeDir.
		CfgFile string
		// Logger configuration file URL.
		LogCfgFile string
		Metrics    bool
	}
)

const (
	// The prefix for configuration keys inside environment.
	envPrefix = "AB"
	// The default name for config file.
	defaultConfigFile = "config.props"
	// the default alphabill directory.
	defaultAlphabillDir = ".alphabill"
	// The default logger configuration file name.
	defaultLoggerConfigFile = "logger-config.yaml"
	// The default rootchain directory
	defaultRootChainDir = "rootchain"
	// The configuration key for home directory.
	keyHome = "home"
	// The configuration key for config file name.
	keyConfig = "config"
	// Enables or disables metrics collection
	keyMetrics = "metrics"
)

func (r *baseConfiguration) addConfigurationFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&r.HomeDir, keyHome, "", fmt.Sprintf("set the AB_HOME for this invocation (default is %s)", alphabillHomeDir()))
	cmd.PersistentFlags().StringVar(&r.CfgFile, keyConfig, "", fmt.Sprintf("config file URL (default is $AB_HOME/%s)", defaultConfigFile))
	cmd.PersistentFlags().StringVar(&r.LogCfgFile, "logger-config", defaultLoggerConfigFile, "logger config file URL. Considered absolute if starts with '/'. Otherwise relative from $AB_HOME.")
	cmd.PersistentFlags().BoolVar(&r.Metrics, keyMetrics, false, "Enables metrics collection.")
}

func (r *baseConfiguration) initConfigFileLocation() {
	// Home directory and config file are special configuration values as these are used for loading in rest of the configuration.
	// Handle these manually, before other configuration loaded with Viper.

	// Home dir is loaded from command line argument. If it's not set, then from env. If that's not set, then default is used.
	if r.HomeDir == "" {
		r.HomeDir = os.Getenv(envKey(keyHome))
		if r.HomeDir == "" {
			r.HomeDir = alphabillHomeDir()
		}
	}

	// Config file name is loaded from command line argument. If it's not set, then from env. If that's not set, then default is used.
	if r.CfgFile == "" {
		r.CfgFile = os.Getenv(envKey(keyConfig))
		if r.CfgFile == "" {
			r.CfgFile = defaultConfigFile
		}
	}
	if !filepath.IsAbs(r.CfgFile) {
		// Config file name is using relative path
		r.CfgFile = filepath.Join(r.HomeDir, r.CfgFile)
	}
}

func (r *baseConfiguration) configFileExists() bool {
	_, err := os.Stat(r.CfgFile)
	return err == nil
}

func (r *baseConfiguration) defaultRootGenesisDir() string {
	return filepath.Join(r.HomeDir, defaultRootChainDir)
}

func envKey(key string) string {
	return strings.ToUpper(envPrefix + "_" + key)
}

func alphabillHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic("default user home dir not defined")
	}
	return filepath.Join(dir, defaultAlphabillDir)
}

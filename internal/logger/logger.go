package logger

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"

	yaml "gopkg.in/yaml.v3"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/errors"
)

func init() {
	initializeGlobalFactory()
}

func initializeGlobalFactory() {
	globalFactoryImpl = &globalFactory{
		loggers:                 make(map[string]*ContextLogger),
		context:                 make(Context),
		consoleTimeFormat:       "15:04:05.000",
		callerSkipFrames:        4, // This depends on the logger code, not meant to be changed by callers.
		packageNameResolver:     &PackageNameResolver{BasePackage: "alphabill/alphabill"},
		nonAlphaNumericRegex:    regexp.MustCompile(`[^a-zA-Z0-9]`),
		globalLoggerInitialized: false,
	}
}

func loadGlobalConfigFromFile(fileName string) (GlobalConfig, error) {
	type (
		LoggerConfiguration struct {
			DefaultLevel    string            `yaml:"defaultLevel"` // tags enable to parse yaml file for the configuration in the future
			PackageLevels   map[string]string `yaml:"packageLevels"`
			OutputPath      string            `yaml:"outputPath"`
			ConsoleFormat   bool              `yaml:"consoleFormat"`
			ShowCaller      bool              `yaml:"showCaller"`
			TimeLocation    string            `yaml:"timeLocation"`
			ShowGoroutineID bool              `yaml:"showGoroutineID"`
		}
	)

	yamlFile, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		pe, ok := err.(*os.PathError)
		if ok {
			return GlobalConfig{}, errors.Wrap(errors.ErrFileNotFound, pe.Error())
		}
		return GlobalConfig{}, errors.Wrap(err, "failed to read logger config file")
	}
	config := &LoggerConfiguration{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return GlobalConfig{}, errors.Wrap(err, "failed to unmarshal logger config")
	}

	// --- Setup globals
	globalConfig := GlobalConfig{
		DefaultLevel:    LevelFromString(config.DefaultLevel),
		PackageLevels:   make(map[string]LogLevel),
		Writer:          nil,
		ConsoleFormat:   config.ConsoleFormat,
		ShowCaller:      config.ShowCaller,
		TimeLocation:    config.TimeLocation,
		ShowGoroutineID: config.ShowGoroutineID,
	}
	// Output writer
	if config.OutputPath != "" {
		file, err := os.Create(config.OutputPath)
		if err != nil {
			return GlobalConfig{}, errors.Wrap(err, "failed to create output writer")
		}
		globalConfig.Writer = bufio.NewWriter(file)
	} else {
		globalConfig.Writer = os.Stdout
	}
	// Log levels for individual packages
	for k, v := range config.PackageLevels {
		globalConfig.PackageLevels[k] = LevelFromString(v)
	}

	return globalConfig, nil
}

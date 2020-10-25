package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kushsharma/go-kafutil/cmd"
	"github.com/kushsharma/go-kafutil/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// Version is app version
	Version = "0"
	// AppName of this executable
	AppName = "kwrite"

	Config config.App
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	initConfig()

	rootCmd := cmd.InitCommands(AppName, Version, Config)
	rootCmd.Execute()
}

func initConfig() {
	viper.SetDefault("host", "localhost:9092")
	viper.SetDefault("topic", "kafutil-default")
	viper.SetDefault("descriptor_set_path", "./protos/desc.set")
	viper.SetDefault("schema", "kafutil.internal.KafutilSample")

	viper.SetEnvPrefix("KAFUTIL")
	viper.SetConfigName(".kafutil")
	viper.SetConfigType("yaml")
	if currentHomeDir, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(currentHomeDir, ".config"))
	}
	viper.AddConfigPath(".")      // directory of binary
	viper.AddConfigPath("../../") // when running in debug mode
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
		} else {
			panic(fmt.Errorf("unable to read optimus config file %v", err))
		}
	}
	Config.Host = viper.GetString("host")
	Config.Topic = viper.GetString("topic")
	Config.Schema = viper.GetString("schema")
	Config.DescriptorSetPath = viper.GetString("descriptor_set_path")
}

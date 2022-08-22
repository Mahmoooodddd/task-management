package config

import (
	"flag"

	"github.com/spf13/viper"
)

const (
	FileName           = "config"
	TestConfigFileName = "config_test"
	FilePath           = "./config"
	TestConfigFilePath = "../config"
)

func SetConfigs() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yml")
	v.SetConfigName(FileName)

	if !isTestEnv() {
		v.AddConfigPath(FilePath)
	} else {
		v.AddConfigPath(TestConfigFilePath)
	}

	err := v.ReadInConfig()
	if err != nil {
		panic("can not read config file" + err.Error())
	}

	//rewrite configs for test env
	if isTestEnv() {
		v.SetConfigName(TestConfigFileName)
		err = v.MergeInConfig()
		if err != nil {
			panic("can not merge config file" + err.Error())
		}
	}

	_ = v.AllSettings()
	return v
}

func isTestEnv() bool {
	return flag.Lookup("test.v") != nil
}

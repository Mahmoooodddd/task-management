package platform

import (
	"flag"
	"github.com/spf13/viper"
)


const (
	domainConfigKey = "task-management.domain"
	EnvConfigKey    = "task-management.environment"
	EnvProd         = "prod"
	EnvTest         = "test"
	EnvDev          = "dev"
)


type Configs interface {
	GetEnv() string
	GetString(name string) string
	GetInt(name string) int
	GetDomain() string
	Set(key string, value interface{})
}

type configs struct {
	viper *viper.Viper
}

func (c *configs) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

func (c *configs) GetString(name string) string {
	return c.viper.GetString(name)
}

func (c *configs) GetInt(name string) int {
	return c.viper.GetInt(name)
}

func (c *configs) GetEnv() string {
	return c.GetString(EnvConfigKey)
}

func (c *configs) GetDomain() string {
	return c.GetString(domainConfigKey)

}

func NewConfigs(viper *viper.Viper) Configs {
	c := &configs{viper}
	if flag.Lookup("test.v") != nil {
		c.Set(EnvConfigKey, EnvTest)
	}
	return c
}

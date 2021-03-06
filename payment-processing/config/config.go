package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"spenmo/payment-processing/payment-processing/internal/pkg/constant"
	"spenmo/payment-processing/payment-processing/internal/pkg/timeutil"
)

// AppConfig is config object to use across application
var AppConfig config

type EnvName string

const (
	ProductionEnv = "production"
	TestEnv       = "test"
	UatEnv        = "uat"
)

func NewEnvName(name string) EnvName {
	switch name {
	case "test":
		return TestEnv
	case "uat":
		return UatEnv
	case "production":
		return ProductionEnv
	}
	return ""
}

type Country struct {
	Name string
	Dial string
	Code string
}

type config struct {
	AppName string
	Env     EnvName

	RedisAddress                string
	RedisPort                   string
	RedisConnPool               int
	RedisReadWriteTimeoutMillis int64
	RedisMaxIdle                int
	RedisMaxIdleTimeoutMinute   int64

	MySQLUser    string
	MySQLPass    string
	MySQLHost    string
	MySQLPort    string
	MySQLDB      string
	MySQLMaxIdle int
	MySQLMaxOpen int

	RateLimitEnabled bool
	RateLimitValue   int64
	RateLimitUnit    string
}

func parseConfigFilePath() string {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(workPath, "config")
}

func InitializeAppConfig() {
	configPath := parseConfigFilePath()
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if AppConfig.AppName = viper.GetString("appname"); AppConfig.AppName == "" {
		panic("appName is missing in config")
	}

	AppConfig.Env = TestEnv
	if env := NewEnvName(viper.GetString("env")); env != "" {
		AppConfig.Env = env
	}

	initRedis(&AppConfig)
	initMySQL(&AppConfig)
	initRateLimit(&AppConfig)
}

func initRedis(c *config) {
	if c.RedisAddress = viper.GetString("redisaddress"); c.RedisAddress == "" {
		panic("redis address is missing in config")
	}
	if c.RedisPort = viper.GetString("redisport"); c.RedisPort == "" {
		panic("redis port is missing in config")
	}

	if c.RedisConnPool = viper.GetInt("redisconnpool"); c.RedisConnPool == 0 {
		c.RedisConnPool = 1000
	}

	if c.RedisReadWriteTimeoutMillis = viper.GetInt64("redisreadwritetimeoutmillis"); c.RedisReadWriteTimeoutMillis == 0 {
		c.RedisReadWriteTimeoutMillis = timeutil.ConvertMinuteToMillis(1)
	}

	if c.RedisMaxIdle = viper.GetInt("redismaxidle"); c.RedisMaxIdle == 0 {
		c.RedisMaxIdle = 5
	}

	if c.RedisMaxIdleTimeoutMinute = viper.GetInt64("redismaxidletimeoutminute"); c.RedisMaxIdleTimeoutMinute == 0 {
		c.RedisMaxIdleTimeoutMinute = 2
	}
}

func initRateLimit(c *config) {
	c.RateLimitEnabled = viper.GetBool("ratelimitenabled")

	if c.RateLimitEnabled {
		if c.RateLimitValue = viper.GetInt64("ratelimitvalue"); c.RateLimitValue <= 0 {
			panic("rate limit value missing in config")
		}
		if c.RateLimitUnit = viper.GetString("ratelimitunit"); c.RateLimitUnit != constant.Minute && c.RateLimitUnit != constant.Hour {
			panic("rate limit unit invalid or missing in config")
		}
	}
}

func initMySQL(c *config) {
	if c.MySQLUser = viper.GetString("mysqluser"); c.MySQLUser == "" {
		panic("mysql username is missing in config")
	}
	if c.MySQLPass = viper.GetString("mysqlpass"); c.MySQLPass == "" {
		panic("mysql password is missing in config")
	}
	if c.MySQLHost = viper.GetString("mysqlhost"); c.MySQLHost == "" {
		panic("mysql host is missing in config")
	}
	if c.MySQLPort = viper.GetString("mysqlport"); c.MySQLPort == "" {
		panic("mysql port is missing in config")
	}
	if c.MySQLDB = viper.GetString("mysqldb"); c.MySQLDB == "" {
		panic("mysql database is missing in config")
	}
	if c.MySQLMaxIdle = viper.GetInt("mysqlmaxidle"); c.MySQLMaxIdle == 0 {
		c.MySQLMaxIdle = 2
	}
	c.MySQLMaxOpen = viper.GetInt("mysqlmaxopen")
}

package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type (
	PgConfig struct {
		DefaultDSN   string `yaml:"DefaultDSN"`   // postgres default database connect dsn
		IfindPgDSN   string `yaml:"IfindPgDSN"`   // postgres ifind database connect dsn
		QueryTimeout int    `yaml:"QueryTimeout"` // time out of pg query
		MaxIdleConns int    `yaml:"MaxIdleConns"` // max number of idles existed
		MaxOpenConns int    `yaml:"MaxOpenConns"` // max number of idles opened
		LogLevel     string `yaml:"LogLevel"`     // log level of pg connection
	}

	ServiceConfig struct {
		HttpPort int `yaml:"HttpPort"` // http port
	}

	LogConfig struct {
		LogPath     string `yaml:"LogPath"`     // log file path
		StatLogPath string `yaml:"StatLogPath"` // status log file path
		GinLogPath  string `yaml:"GinLogPath"`  // gin log file path
		LogLevel    string `yaml:"LogLevel"`    // log level
	}

	DataApiConfig struct {
		AppName string    `yaml:"AppName"` // app name
		RpsCfg  RpsConfig `yaml:"Rps"`     // rps config
		FxjCfg  FxjConfig `yaml:"Fxj"`     // fxj config
	}

	Config struct {
		PgConfig      PgConfig      `yaml:"Pgsql"`
		ServiceConfig ServiceConfig `yaml:"Service"`
		LogConfig     LogConfig     `yaml:"Log"`
		DataApiConfig DataApiConfig `yaml:"Api"`
	}
)

type (
	RpsConfig struct {
		Url  string `yaml:"Url"`
		Cron string `yaml:"Cron"`
	}
	FxjConfig struct {
		Url     string `yaml:"Url"`
		Cron    string `yaml:"Cron"`
		Markets struct {
			Bond  []uint8 `yaml:"Bond"`
			Fund  []uint8 `yaml:"Fund"`
			Index []uint8 `yaml:"Index"`
		} `yaml:"Market"`
	}
)

var (
	cfgFile *string
	cfg     Config
)

/*ConfigureInit
 * @Description: 获取配置
 */
func ConfigureInit() {
	// 注释中为线上环境配置
	//cfgFile = flag.String("f", "/usr/local/conf/conf.yaml", "config file path")
	cfgFile = flag.String("f", "/Users/heqimin/Code/Go/finance/importer/conf/conf.yaml",
		"config file path")
	flag.Parse()
	content, err := ioutil.ReadFile(*cfgFile)
	if yaml.Unmarshal(content, &cfg) != nil {
		log.Fatalf("解析config.yaml出错: %v", err)
	}
	log.Println(cfg)
}

func GetPgConfig() PgConfig {
	return cfg.PgConfig
}

func GetServiceConfig() ServiceConfig {
	return cfg.ServiceConfig
}

func GetLogConfig() LogConfig {
	return cfg.LogConfig
}

func GetApiConfig() DataApiConfig {
	return cfg.DataApiConfig
}

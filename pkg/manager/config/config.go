/**
 * @Time: 2023/10/20 16:29
 * @Author: jzechen
 * @File: dto.go
 * @Software: GoLand collector
 */

package config

import (
	"errors"
	"github.com/jzechen/toresa/cmd/manager/cmd/options"
	"github.com/jzechen/toresa/pkg/manager/config/scrape"
	. "github.com/jzechen/toresa/pkg/manager/contants"
	"github.com/jzechen/toresa/pkg/manager/utils"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"os"
	"path"
	"time"
)

type CollectorManager struct {
	OriginConf *viper.Viper
	Server     ServerConfig `yaml:"server"`
	Mongo      MongoConfig  `yaml:"mongo"`
	Scraper    ScrapeConfig `yaml:"scraper"`
	Drive      DriveConfig  `yaml:"drive"`
}

type ServerConfig struct {
	Addr           string        `yaml:"addr"`
	Port           int           `yaml:"port"`
	Limit          int           `yaml:"limit"`
	Burst          int           `yaml:"burst"`
	RequestTimeout time.Duration `yaml:"requestTimeout"`
}

type MongoConfig struct {
	Addr        string        `yaml:"addr"`
	Database    string        `yaml:"database"`
	DialTimeout time.Duration `yaml:"dialTimeout"`
}

type ScrapeConfig struct {
	RuntimePath string                 `yaml:"runtimePath"`
	Sina        scrape.SinaWeiboConfig `yaml:"sina"`
}

type DriveConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
	Port int    `yaml:"port"`
}

func BuildConfig(opts *options.CollectorManagerOptions) (*CollectorManager, *viper.Viper, error) {
	conf := viper.New()

	if opts.CfgFile != "" {
		// Use config file from the flag.
		conf.SetConfigFile(opts.CfgFile)
	} else {
		// Find work directory
		var wd string
		wd, err := os.Getwd()
		if err != nil {
			return nil, nil, err
		}
		// Search config in work directory with name "config" (without extension).
		conf.AddConfigPath(wd)
		conf.SetConfigType("yaml")
		conf.SetConfigName("config")
	}
	conf.AutomaticEnv()
	err := conf.ReadInConfig()
	if err != nil {
		return nil, nil, err
	}
	klog.Infoln("Using config file:", conf.ConfigFileUsed())

	var _config CollectorManager
	err = conf.Unmarshal(&_config)
	if err != nil {
		return nil, nil, err
	}
	if klog.V(4).Enabled() {
		klog.Infoln("Using config: ", _config)
	}

	return &_config, conf, nil
}

func (conf *CollectorManager) SetDefaultConfig() {
	// server
	if conf.Server.Addr == "" {
		conf.Server.Addr = DefaultServerAddr
	}
	if conf.Server.Port == 0 {
		conf.Server.Port = DefaultServerPort
	}
	if conf.Server.Limit < 0 {
		conf.Server.Limit = 10
	}
	if conf.Server.Burst < conf.Server.Limit {
		conf.Server.Burst = conf.Server.Limit
	}
	if conf.Server.RequestTimeout == 0 {
		conf.Server.RequestTimeout = time.Duration(60) * time.Second
	}

	// drive
	if conf.Drive.Type == "" {
		conf.Drive.Type = DefaultDriveType
	}
	if conf.Drive.Port == 0 {
		conf.Drive.Port = DefaultDrivePort
	}
	if conf.Drive.Path == "" {
		conf.Drive.Path = DefaultDrivePath
	}
	_, err := os.Stat(conf.Drive.Path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	if errors.Is(err, os.ErrNotExist) {
		conf.Drive.Path = path.Join(utils.ExecPath, conf.Drive.Path)
	}
}

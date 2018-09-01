package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

var (
	Conf              conf
	defaultConfigFile = "conf/config.toml"
)

type conf struct {
	Server server
	Sql    sql
	Token  token
}

type server struct {
	Addr string `toml:"addr"`
}

type sql struct {
	User       string `toml:"user"`
	Password   string `toml:"password"`
	Addr       string `toml:"addr"`
	DB         string `toml:"db"`
	TimeLayout string `toml:"time_layout"`
}

type token struct {
	TokenExpire int    `toml:"token_expire"`
	Salt        string `toml:"salt"`
}

func InitConfig(runmode string) error {
	var configFile string
	if runmode == "" || runmode == "prod" {
		configFile = defaultConfigFile
	} else if runmode == "dev" {
		configFile = "conf/config_dev.toml"
	} else {
		return fmt.Errorf("runmode error")
	}

	if _, err := os.Stat(configFile); err != nil {
		return err
	}
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(configBytes), &Conf)
	if err != nil {
		return err
	}
	return nil
}

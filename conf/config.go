package conf

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"time"
)

var (
	Conf              conf
	defaultConfigFile = "conf/config.toml"
)

type conf struct {
	Server server
}

type server struct {
	Addr string `toml:"addr"`
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

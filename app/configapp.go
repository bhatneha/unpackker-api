package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
)

const (
	AppPort   = "APP_PORT"
	UIPort    = "UI_PORT"
	LogPath   = "LOG_PATH"
	UI        = "UI"
	AppConfig = "APP_CONFIG"
)

func (c *ConfigData) configLog() (io.Writer, error) {

	if !statFile(os.Getenv(LogPath)) {
		return nil, fmt.Errorf("Unable to find the configuration file in the specified path, writing to STDOUT")
	}
	file, err := os.OpenFile(c.LogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("Unable to open log file, writing to STDOUT")
	}
	return file, nil

}

func getConfig() (*ConfigData, error) {

	//checking if the environment variable is set
	if len(os.Getenv(AppConfig)) == 0 {
		fmt.Fprintf(os.Stdout, "Configuration file was not specified, switching to default configuration!\n")
		con, err := decodeConfig([]byte(defaultconfig))
		if err != nil {
			return nil, err
		}
		return con, nil
	}

	//checking if there is a file in the specified path set using env variable
	if !statFile(os.Getenv(AppConfig)) {
		fmt.Fprintf(os.Stdout, "Unable to locate configuration file, switching to default configuration\n")
		con, err := decodeConfig([]byte(defaultconfig))
		if err != nil {
			return nil, err
		}
		return con, nil
	}

	//if config file is present, take configurations from there
	jsonCont, err := ioutil.ReadFile(os.Getenv(AppConfig))
	if err != nil {
		return nil, err
	}
	con, err := decodeConfig(jsonCont)
	if err != nil {
		return nil, err
	}
	return con, nil
}

//unmarshal the configurations
func decodeConfig(dconf []byte) (*ConfigData, error) {

	var conf ConfigData
	if err := json.Unmarshal(dconf, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func statFile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func (c *ConfigData) mergeConfig() error {
	envcon := gatherConfigFromEnv()
	// fmt.Println(envcon)
	conf := new(ConfigData)
	if err := mapstructure.Decode(envcon, &conf); err != nil {
		fmt.Println(err)
		return err
	}
	// if err := mergo.Merge(c, conf, mergo.WithOverride); err != nil {
	// 	return err
	// }
	return nil
}

func gatherConfigFromEnv() map[string]interface{} {
	envs := make(map[string]interface{})
	envs["AppConfig"] = os.Getenv(AppConfig)
	envs["AppPort"] = os.Getenv(AppPort)
	envs["UIPort"] = os.Getenv(UIPort)
	envs["UI"] = os.Getenv(UI)
	envs["LogPath"] = os.Getenv(LogPath)

	return envs
}

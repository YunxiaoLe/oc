package conf

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func ReadSettingsFromFile(settingFilePath string) (config *Config) {
	var Config Config
	jsonFile, err := os.Open(settingFilePath)
	if err != nil {
		panic(any("No such file named " + settingFilePath))
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &Config)
	if err != nil {
		logrus.Error(err)
	}
	config = &Config
	return config
}

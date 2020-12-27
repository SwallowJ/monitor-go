package conf

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/SwallowJ/loggo"
	"gopkg.in/yaml.v2"
)

var logger *loggo.Logger

func init() {

	logger = loggo.New("config")
	logger.SetServiceName("config")

	source := "./config"
	if _, err := os.Stat(source); err != nil {
		logger.Fatal(err)
	}

	loadSource(source)

	if Config.Debug {
		loggo.SetLevel(loggo.LevelDebug)
	}
}

func loadSource(source string) {
	sourceInfo, err := ioutil.ReadDir(source)
	if err != nil {
		logger.Error(err)
	}

	for _, f := range sourceInfo {
		file := path.Join(source, f.Name())
		if f.IsDir() {
			loadSource(file)
		} else if path.Ext(file) == ".yml" || path.Ext(file) == ".yaml" {
			if err := load(file); err != nil {
			}
		}
	}
}

func load(filePath string) error {
	logger.Debug("加载配置文件:", filePath)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		return err
	}
	return nil
}

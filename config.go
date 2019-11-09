package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

/* Получить параметры из конфигурационного YAML файла */
func getConfigParameters(configFileName string, config *Config) {

	// Открыть файл
	file, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Fail to open config file '%s': %v", configFileName, err)
	}
	defer func() {
		err = file.Close()
	}()

	// Прочитать весь файл
	yamlData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Fail to read config file '%s': %v", configFileName, err)
	}

	// Десериализовать
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		log.Fatal("Fail to unmarshal config file '%s': %v", configFileName, err)
		os.Exit(1)
	}

	// Перевыставить уровень логирования на основе конфига
	debugLevel := "DEBUG" // "DEBUG", "INFO"
	if config.LogLevel == "INFO" {
		debugLevel = "INFO"
		SetLog(log.InfoLevel)
	}
	log.Debugf("Log level: '%s'", debugLevel)

	log.Debugf("=--> config:\n%v\n\n", *config)

	err = file.Close()
	if err != nil {
		log.Errorf("Fail to close config file '%s': %v", configFileName, err)
	}
}

/*********************************************************/
/* Программа для автоматического создания резевных копий */
/*********************************************************/

package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// Параметры каждого из Бекапов
type Backup struct {
	From        string `yaml:"From"`
	To          string `yaml:"To"`
	OutFileName string `yaml:"OutFileName"`
	Crypt       bool   `yaml:"Crypt"`
	DateTime    bool   `yaml:"DateTime"`
	Ssh         string `yaml:"Ssh"`
}

// Конфиг полностью
type Config struct {
	LogLevel       string   `yaml:"LogLevel"`
	ThreadsCount   int      `yaml:"ThreadsCount"`
	DateTimeFormat string   `yaml:"DateTimeFormat"`
	Backups        []Backup `yaml:"Backups"`
}

func main() {

	const configFileName = "gobackuper.yaml"

	var (
		config Config
	)

	// Выставить параметры логирования
	SetLog(log.DebugLevel)

	// Получить параметры из конфигурационного файла
	getConfigParameters(configFileName, &config)

}

/* Получить параметры из конфигурационного YAML файла */
func getConfigParameters(configFileName string, config *Config) {

	// Открыть файл
	file, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Fail to open config file '%s': %v", configFileName, err)
	}

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

/* Выставить параметры логирования */
func SetLog(debugLevel log.Level) {
	log.SetOutput(os.Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(debugLevel)
}

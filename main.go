/*********************************************************/
/* Программа для автоматического создания резевных копий */
/*********************************************************/

package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Config1 struct {
	// Общие параметры приложения
	LogLevel       string `yaml:"LogLevel"`
	ThreadsCount   string `yaml:"ThreadsCount"`
	DateTimeFormat string `yaml:"DateTimeFormat"`

	// Параметры каждого из Бекапов
	From        string `yaml:"From"`
	To          string `yaml:"To"`
	OutFileName string `yaml:"OutFileName"`
	Crypt       string `yaml:"Crypt"`
	DateTime    string `yaml:"DateTime"`
	Ssh         string `yaml:"Ssh"`
}

type Backup struct {
	// Параметры каждого из Бекапов
	From        string `yaml:"From"`
	To          string `yaml:"To"`
	OutFileName string `yaml:"OutFileName"`
	Crypt       bool   `yaml:"Crypt"`
	DateTime    bool   `yaml:"DateTime"`
	Ssh         string `yaml:"Ssh"`
}

type Config2 struct {
	LogLevel       string `yaml:"LogLevel"`
	ThreadsCount   int    `yaml:"ThreadsCount"`
	DateTimeFormat string `yaml:"DateTimeFormat"`
	Backups        []Backup
}

func main() {

	const configFileName = "gobackuper.yaml"

	var (
		//debugLevel = "DEBUG" // "DEBUG", "INFO"
		//configs    []Config1
		config3 Config2
	)

	// Выставить параметры логирования
	SetLog(log.DebugLevel)

	// Получить параметры из конфигурационного файла
	//getConfigParameters(configFileName, &configs)

	// Перевыставить уровень логирования на основе конфига
	//if configs[0].LogLevel == "INFO" {
	//	debugLevel = "INFO"
	//	SetLog(log.InfoLevel)
	//}
	//log.Infof("Log level: '%s'", debugLevel)

	// Сериализовать в YAML

	config3.LogLevel = "INFO"

	var config2 Backup
	config2.From = "AAA"
	config2.To = "BBB"
	config2.Crypt = true
	config2.DateTime = true
	config2.Ssh = "test.com.ru"

	var config5 Backup
	config5.From = "CCC"
	config5.To = "DDD"
	config5.Crypt = false
	config5.DateTime = true
	config5.Ssh = "qa.com.ru"

	config3.Backups = append(config3.Backups, config2)
	config3.Backups = append(config3.Backups, config5)

	d, err := yaml.Marshal(&config3)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- config3 dump:\n%s\n\n", string(d))

}

/* Получить параметры из конфигурационного YAML файла */
func getConfigParameters(configFileName string, configs *[]Config1) {

	// Открыть файл
	file, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Fail to open configs file '%s': %v", configFileName, err)
	}

	// Прочитать весь файл
	yamlData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Fail to read configs file '%s': %v", configFileName, err)
	}

	// Десериализовать
	err = yaml.Unmarshal(yamlData, &configs)
	if err != nil {
		log.Fatal("Fail to unmarshal configs file '%s': %v", configFileName, err)
		os.Exit(1)
	}
	log.Infof("=--> configs:\n%v\n\n", *configs)

	err = file.Close()
	if err != nil {
		log.Errorf("Fail to close configs file '%s': %v", configFileName, err)
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

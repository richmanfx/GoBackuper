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

type Config struct {
	// Общие параметры приложения
	LogLevel       string `LogLevel`
	ThreadsCount   string `ThreadsCount`
	DateTimeFormat string `DateTimeFormat`

	// Параметры каждого из Бекапов
	From        string `From`
	To          string `To`
	OutFileName string `OutFileName`
	Crypt       string `Crypt`
	DateTime    string `DateTime`
	Ssh         string `Ssh`
}

func main() {

	const configFileName = "gobackuper.yaml"

	var (
		debugLevel = "DEBUG" // "DEBUG", "INFO"
		configs    []Config
	)

	// Выставить параметры логирования
	SetLog(log.DebugLevel)

	// Получить параметры из конфигурационного файла
	getConfigParameters(configFileName, &configs)

	// Перевыставить уровень логирования на основе конфига
	if configs[0].LogLevel == "INFO" {
		debugLevel = "INFO"
		SetLog(log.InfoLevel)
	}
	log.Infof("Log level: '%s'", debugLevel)

}

/* Получить параметры из конфигурационного YAML файла */
func getConfigParameters(configFileName string, configs *[]Config) {

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

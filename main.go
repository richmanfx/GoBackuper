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
	Debuglevel string
	//B struct {
	//	RenamedC int   `yaml:"c"`
	//	D        []int `yaml:",flow"`
	//}
}

//var data = `
//debuglevel: DEBUG
//`

func main() {

	const configFileName = "gobackuper.yaml"

	var (
		debugLevel = "DEBUG"
	)

	// Выставить параметры логирования
	SetLog(log.DebugLevel)

	// Получить параметры из конфигурационного файла
	getConfigParameters(configFileName, debugLevel)

}

/* Получить параметры из конфигурационного YAML файла */
func getConfigParameters(configFileName string, debugLevel string) {

	config := Config{}
	var err error

	// Читать файл
	file, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Fail to open config file '%s': %v", configFileName, err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Fail to read config file '%s': %v", configFileName, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Fail to unmarshal config file '%s': %v", configFileName, err)
		os.Exit(1)
	}

	log.Infof("=--> config:\n%v\n\n", config)

	//debugLevel = config.Section("").Key("DEBUG_LEVEL").String()
	//if debugLevel == "INFO" {
	//	SetLog(log.InfoLevel)
	//}
	//log.Debugf("Debug level: %s", debugLevel)

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

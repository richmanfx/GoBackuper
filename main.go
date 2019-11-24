/*********************************************************/
/* Программа для автоматического создания резевных копий */
/*********************************************************/

package main

import (
	log "github.com/Sirupsen/logrus"
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
	LogLevel             string   `yaml:"LogLevel"`
	CompressionLevel     int      `yaml:"CompressionLevel"`
	SelectiveCompression bool     `yaml:"SelectiveCompression"`
	ThreadsCount         int      `yaml:"ThreadsCount"`
	DateTimeFormat       string   `yaml:"DateTimeFormat"`
	Backups              []Backup `yaml:"Backups"`
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

	// В ZIP
	toZip(&config)

	// Перенести LZMA архив в директорию хранения
	moveArchive(&config)

}

func moveArchive(config *Config) {

	for _, backup := range config.Backups {

		var oldLocation string

		if backup.DateTime == false {
			oldLocation = backup.OutFileName + ".zip"
		} else if backup.DateTime == true {
			dataTimeSuffix := getSuffix(config.DateTimeFormat)
			oldLocation = backup.OutFileName + "-" + dataTimeSuffix + ".zip"
		}

		newLocation := backup.To + "\\" + oldLocation

		err := os.Rename(oldLocation, newLocation)
		errLogAndExit(err)

	}

}

/* Обработать фатальную ошибку и выйти */
func errLogAndExit(err error) {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

/* Логировать ошибку */
func errLog(err error) {
	if err != nil {
		log.Errorln(err)
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

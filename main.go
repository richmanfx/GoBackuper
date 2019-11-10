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

	// В TAR архив
	toTar(&config)

	// Сжать LZMA2
	toLZMA2(&config)

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

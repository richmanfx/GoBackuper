/*********************************************************/
/* Программа для автоматического создания резевных копий */
/*********************************************************/

package main

import (
	"archive/tar"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

}

/* Поместить директории в TAR архивы */
func toTar(config *Config) {

	// Бежать по конфигам бекапов
	for _, backup := range config.Backups {

		parentDir := ""
		parentDirs := strings.Split(backup.From, string(filepath.Separator))
		currentDir := parentDirs[len(parentDirs)-1:][0] // Последний оставляем
		parentDirs = parentDirs[1:len(parentDirs)]      // Первый удаляем
		parentDirs = parentDirs[:len(parentDirs)-1]     // Последний удаляем
		log.Infof("parentDirs: %v", parentDirs)
		for _, dir := range parentDirs {
			parentDir = parentDir + string(filepath.Separator) + dir
		}

		addToTarArchive(backup.OutFileName+".tar", currentDir, parentDir)
	}
}

func addToTarArchive(tarFileName string, currentDir string, parentDir string) {

	entry := parentDir + string(filepath.Separator) + currentDir
	//entry := currentDir

	// Открыть директорию
	dir, err := os.Open(entry)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	// Все файлы в директории
	files, err := dir.Readdir(0)
	if err != nil {
		log.Errorln(err)
	}

	// Создать TAR файл
	tarFile, err := os.Create(tarFileName)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	var fileWriter io.WriteCloser = tarFile
	tarFileWriter := tar.NewWriter(fileWriter)

	for _, fileInfo := range files {

		// Если файл - директория
		if fileInfo.IsDir() {
			log.Infof("Есть директория: %s", fileInfo.Name())

			// Файлы-дети в директории
			//File[] children = file.listFiles();
			subDir, err := os.Open(parentDir + string(filepath.Separator) + currentDir + string(filepath.Separator) + fileInfo.Name())
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}
			subFiles, err := subDir.Readdir(0)
			if len(subFiles) != 0 {
				for _, subFile := range subFiles {
					currentSubDirs := strings.Split(subFile.Name(), string(filepath.Separator))
					currentSubDir := currentSubDirs[len(currentSubDirs)-1:][0] // Последний оставляем
					addToTarArchive(tarFileName, currentSubDir, entry)
				}
			}
			//continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		if err != nil {
			log.Errorln(err)
		}

		// Подготовка TAR заголовков
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()
		header.Format = tar.FormatPAX

		err = tarFileWriter.WriteHeader(header)
		if err != nil {
			log.Errorln(err)
		}

		_, err = io.Copy(tarFileWriter, file)
		if err != nil {
			log.Errorln(err)
		}

		err = file.Close()
		if err != nil {
			log.Errorln(err)
		}
	}

	err = tarFileWriter.Close()
	if err != nil {
		log.Errorln(err)
	}

	err = tarFile.Close()
	if err != nil {
		log.Errorln(err)
	}

	err = dir.Close()
	if err != nil {
		log.Errorln(err)
	}
}

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

/* Выставить параметры логирования */
func SetLog(debugLevel log.Level) {
	log.SetOutput(os.Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(debugLevel)
}

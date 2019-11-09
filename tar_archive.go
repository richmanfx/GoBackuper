package main

import (
	"archive/tar"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

/* Поместить директории в TAR архивы */
func toTar(config *Config) {

	// Бежать по конфигам бекапов
	for _, backup := range config.Backups {
		wg.Add(1) // Новая горутина
		go backupToTar(backup)
	}

	wg.Wait() // Ждать окончания работы всех горутин
}

func backupToTar(backup Backup) {

	defer wg.Done() // После окончания работы функции счётчик именьшить на 1

	parentDir := ""
	parentDirs := strings.Split(backup.From, string(filepath.Separator))
	parentDirs = parentDirs[1 : len(parentDirs)-1] // Первый и последний удаляем
	log.Debugf("parentDirs: %v", parentDirs)

	// Собрать в кучу
	for _, dir := range parentDirs {
		parentDir = parentDir + string(filepath.Separator) + dir
	}

	file, err := os.Open(backup.From)
	errLogAndExit(err)

	// Создать TAR файл
	tarFile, err := os.Create(backup.OutFileName + ".tar")
	errLogAndExit(err)

	// Создать Writer
	var fileWriter io.WriteCloser = tarFile
	tarFileWriter := tar.NewWriter(fileWriter)

	addToTarArchive(tarFileWriter, file, parentDir)

	err = tarFileWriter.Close()
	errLog(err)

	err = tarFile.Close()
	errLog(err)

	err = file.Close()
	errLog(err)
}

func addToTarArchive(tarFile *tar.Writer, fileOrDirToArchive *os.File, parentDir string) {

	entry := fileOrDirToArchive.Name()

	fileInfo, err := os.Stat(entry)
	errLog(err)

	if fileInfo.Mode().IsRegular() {
		// Если файл

		file, err := os.Open(entry)
		errLog(err)

		// Подготовить TAR заголовок
		header := new(tar.Header)
		header.Name, err = filepath.Rel(parentDir, file.Name()) // Исключаем полный путь
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()
		header.Format = tar.FormatPAX

		// Добавить TAR зоголовок
		err = tarFile.WriteHeader(header)
		errLog(err)

		// Добавить файл
		_, err = io.Copy(tarFile, file)
		errLog(err)

		// Закрыть файл
		err = file.Close()
		errLog(err)

	} else if fileInfo.Mode().IsDir() {
		// Если директория

		// Открыть директорию
		dir, err := os.Open(entry)
		errLogAndExit(err)

		// Все файлы-дети в директории
		fileInfos, err := dir.Readdir(0)
		errLog(err)

		if len(fileInfos) != 0 {
			// Непустая директория
			for _, fileInfo := range fileInfos {

				file, err := os.Open(entry + string(filepath.Separator) + fileInfo.Name())
				errLogAndExit(err)

				addToTarArchive(tarFile, file, parentDir)
			}
		}

		// Закрыть
		err = dir.Close()
		errLog(err)

	} else {
		// Ни файл, ни директория ٩(｡•́‿•̀｡)۶
		log.Errorf("Not directory and not file - '{}'", fileOrDirToArchive)
	}
}

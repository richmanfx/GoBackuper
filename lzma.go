package main

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"os"
)

func toLZMA2(config *Config) {

	// Бежать по конфигам бекапов
	for _, backup := range config.Backups {
		wg.Add(1) // Новая горутина
		go toLZMACompressOneFile(backup.OutFileName)
	}
	wg.Wait() // Ждать окончания работы всех горутин
}

func toLZMACompressOneFile(outFileName string) {

	log.Infof("Start compress: '%s'", outFileName)

	// После окончания работы функции счётчик именьшить на 1
	defer wg.Done()

	// Входной TAR файл
	inFile, err := os.Open(outFileName + ".tar")
	errLogAndExit(err)
	reader := bufio.NewReader(inFile)

	// Выходной LZMA файл
	outFile, err := os.Create(outFileName + ".tar.lzma")
	fileWriter := bufio.NewWriter(outFile)
	lzmaFileWriter, err := lzma.NewWriter(fileWriter)
	errLogAndExit(err)

	// Компресиия
	_, err = io.Copy(lzmaFileWriter, reader)
	errLogAndExit(err)

	// Всё очистить и закрыть
	_ = fileWriter.Flush()
	_ = lzmaFileWriter.Close()
	_ = outFile.Close()
	_ = inFile.Close()

	// Удалить TAR файл
	err = os.Remove(outFileName + ".tar")
	errLog(err)

	log.Infof("Ending compress: '%s'", outFileName)
}

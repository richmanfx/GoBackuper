package main

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"os"
)

func toLZMA2(config *Config) {

	//pipeReader, pipeWriter := io.Pipe()
	//defer func() {
	//	err := pipeReader.Close()
	//	errLog(err)
	//}()

	// Source file
	inFile, err := os.Open(config.Backups[0].OutFileName + ".tar")
	errLog(err)
	//myDefer(inFile)
	reader := bufio.NewReader(inFile)

	// Destination file
	outFile, err := os.Create(config.Backups[0].OutFileName + ".tar.lzma")
	//myDefer(outFile)

	var fileWriter io.WriteCloser = outFile
	lzmaFileWriter, err := lzma.NewWriter(fileWriter)

	//bufioWriter := bufio.NewWriter(outFile)
	//writer, err := lzma.NewWriter2(bufioWriter)

	errLogAndExit(err)

	count, err := io.Copy(lzmaFileWriter, reader)
	errLog(err)
	log.Infof("Скописровано: %d", count)

	//_ = bufioWriter.Flush()
	_ = lzmaFileWriter.Close()
	_ = fileWriter.Close()
	_ = outFile.Close()
	_ = inFile.Close()

}

func myDefer(file *os.File) {
	defer func() {
		err := file.Close()
		errLog(err)
	}()
}

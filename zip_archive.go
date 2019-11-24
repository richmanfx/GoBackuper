package main

import (
	"github.com/mholt/archiver"
	"sync"
)

var wg sync.WaitGroup

func toZip(config *Config) {

	// Бежать по конфигам бекапов
	for _, backup := range config.Backups {
		wg.Add(1) // Новая горутина
		go dirToZip(backup, config.CompressionLevel, config.SelectiveCompression, config.DateTimeFormat)
	}
	wg.Wait() // Ждать окончания работы всех горутин

}

func dirToZip(backup Backup, compressionLevel int, selectiveCompression bool, format string) {

	var source []string
	var destination string
	defer wg.Done() // После окончания работы функции счётчик именьшить на 1

	zipArchive := archiver.Zip{
		CompressionLevel:       compressionLevel,
		OverwriteExisting:      true,
		MkdirAll:               false,
		SelectiveCompression:   selectiveCompression,
		ImplicitTopLevelFolder: false,
		ContinueOnError:        true,
	}

	source = append(source, backup.From)
	if backup.DateTime == false {
		destination = backup.OutFileName + ".zip"
	} else if backup.DateTime == true {
		dataTimeSuffix := getSuffix(format)
		destination = backup.OutFileName + "-" + dataTimeSuffix + ".zip"
	}

	err := zipArchive.Archive(source, destination)
	errLog(err)

	err = zipArchive.Close()
	errLog(err)
}

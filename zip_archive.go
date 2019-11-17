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
		go dirToZip(backup, config.CompressionLevel, config.SelectiveCompression)
	}
	wg.Wait() // Ждать окончания работы всех горутин

}

func dirToZip(backup Backup, compressionLevel int, selectiveCompression bool) {

	var source []string
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

	err := zipArchive.Archive(source, backup.To+"/"+backup.OutFileName+".zip")
	errLog(err)

	err = zipArchive.Close()
	errLog(err)
}

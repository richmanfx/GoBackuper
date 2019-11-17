package main

//
//import (
//	"github.com/mholt/archiver"
//)
//
//func toTar2()  {
//
//	//t := archiver.Tar{
//	//	OverwriteExisting:      true,
//	//}
//
//	//b := archiver.TarBz2{
//	//	Tar:              &t,
//	//	CompressionLevel: 9,
//	//}
//
//	z := archiver.Zip{
//		CompressionLevel:       9,
//		OverwriteExisting:      true,
//		MkdirAll:               false,
//		SelectiveCompression:   false,
//		ImplicitTopLevelFolder: false,
//		ContinueOnError:        true,
//	}
//
//	//b := archiver.TarLz4{
//	//	Tar:              &t,
//	//	CompressionLevel: 12,
//	//}
//
//	var source []string
//	source = append(source, "C:\\Users\\zoer\\Documents")
//
//	err := z.Archive(source, "Документы.zip")
//	errLog(err)
//
//
//	//err := t.Archive(source, "Документы.tar")
//	//errLog(err)
//
//	//err = b.Close()
//	//errLog(err)
//
//	err = z.Close()
//	errLog(err)
//}
//
//

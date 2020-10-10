package cmd

import (
	"os"
	"path/filepath"
	"time"
)

type FileScanner struct {
	Logging *Logging
}

type File struct {
	path string
	info os.FileInfo
}

func(fileScanner FileScanner) list(path string) {

	fileScanner.Logging.Stdout.Printf("List Files")

	startTime := time.Now()

	var files []File
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, File{path, info})
		return nil
	})
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fileScanner.Logging.Stdout.Printf("Scanned %d Files in %d ms", len(files), duration * time.Millisecond)

	if err != nil {
		fileScanner.Logging.ErrorLog.Printf("%v", err)
	}

	for _, file := range files {
		fileScanner.Logging.InfoLog.Printf("%s %d", file.path, file.info.Size())
	}

}

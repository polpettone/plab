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

type FileListResult struct {
	Files []File
	StartTime time.Time
	EndTime time.Time
}

func listFiles(path string) (*FileListResult, error) {
	startTime := time.Now()
	var files []File
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, File{path, info})
		return nil
	})
	endTime := time.Now()

	if err != nil {
		return nil, err
	}

	return &FileListResult{
		files,
		startTime,
		endTime,
	}, nil
}

func(fileScanner FileScanner) list(path string) {
	fileScanner.Logging.Stdout.Printf("List Files")

	fileListResult, err := listFiles(path)

	if err != nil {
		fileScanner.Logging.ErrorLog.Printf("%v", err)
		return
	}

	duration := fileListResult.EndTime.Sub(fileListResult.StartTime).Milliseconds()
	fileScanner.Logging.Stdout.Printf("Scanned %d Files in %d ms", len(fileListResult.Files), duration)

	for _, file := range fileListResult.Files {
		fileScanner.Logging.InfoLog.Printf("%s %d", file.path, file.info.Size())
	}
}

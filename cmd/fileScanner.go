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

type Stat struct {
	FileExtensionMap map[string][]File
}


func NewStat(files []File) Stat {
	extensions := getAllFileExtensions(files)
	fileExtionsionMap := make(map[string][]File)
	for ext := range extensions.m {
		fileExtionsionMap[ext] = filterByExtension(files, ext)
	}
	return Stat{
		FileExtensionMap: fileExtionsionMap,
	}
}

func listFiles(path string) (*FileListResult, error) {
	startTime := time.Now()
	var files []File
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, File{path, info})
		}
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


func getAllFileExtensions(files []File) *Set {
	extensions := NewSet()
	for _, file := range files {
		extensions.Add(filepath.Ext(file.path))
	}
	return extensions
}


func filterByExtension(files []File, extension string) []File {
	var filtered []File
	for _, file := range files {
		if filepath.Ext(file.path) == extension {
			filtered = append(filtered, file)
		}
	}
	return filtered
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

	extensions := getAllFileExtensions(fileListResult.Files)

	stat := NewStat(fileListResult.Files)
	for ext := range extensions.m {
		count := len(stat.FileExtensionMap[ext])
		fileScanner.Logging.InfoLog.Printf("%s %d", ext, count)
	}
}

var exists = struct{}{}

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]struct{})
	return s
}

func (s *Set) Add(value string) {
	s.m[value] = exists
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

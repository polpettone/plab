package cmd

import (
	"os"
	"path/filepath"
	"sort"
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
	FilesByExtension []FilesByExtension
}

type FilesByExtension struct {
	Extension string
	Files []File
}

func NewStat(files []File, logging *Logging) Stat {
	extensions := getAllFileExtensions(files)
	var filesByExtensionSlice []FilesByExtension

	startTime := time.Now()
	logging.DebugLog.Printf("filter files by extension")

	filesByExtensionChan := make(chan *FilesByExtension)

	for ext := range extensions.m {
		go func(ext string) {
			filesByExtension :=  &FilesByExtension{ext, filterByExtension(files, ext)}
			filesByExtensionChan <- filesByExtension
		}(ext)
	}

	for {
		filesByExtension := <-filesByExtensionChan
		filesByExtensionSlice = append(filesByExtensionSlice, *filesByExtension)
		if len(filesByExtensionSlice) == len(extensions.m) {
			break
		}
	}

	endTime := time.Now()
	logging.DebugLog.Printf("filter files by extension done in %d ms", endTime.Sub(startTime).Milliseconds())

	startTime = time.Now()
	logging.DebugLog.Printf("sort files by extensions count")
	sort.Slice(filesByExtensionSlice, func(i, j int) bool {
		return len(filesByExtensionSlice[i].Files) > len(filesByExtensionSlice[j].Files)
	})
	endTime = time.Now()
	logging.DebugLog.Printf("sort files by extensions count done in %d ms", endTime.Sub(startTime).Milliseconds())

	return Stat{
		FilesByExtension: filesByExtensionSlice,
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

	stat := NewStat(fileListResult.Files, fileScanner.Logging)

	fileScanner.Logging.Stdout.Printf("Created Stats")

	for _, filesByExtension := range stat.FilesByExtension {
		fileScanner.Logging.InfoLog.Printf("%s %d", filesByExtension.Extension, len(filesByExtension.Files))
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

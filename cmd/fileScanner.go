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
	FilesByExtensionMap map[string]*FilesByExtension
	Files []File
	StartTime time.Time
	EndTime time.Time
}

type Stat struct {
	FilesByExtension []* FilesByExtension
}

type FilesByExtension struct {
	Extension string
	Files []File
}

func NewStat(files map[string] *FilesByExtension, logging *Logging) Stat {
	var l []*FilesByExtension

	for _ , x := range files {
		l = append(l, x)
	}

	sort.Slice(l, func(i, j int) bool {
		return len(l[i].Files) > len(l[j].Files)
	})

	return Stat{
		FilesByExtension: l,
	}
}

func listFiles(path string) (*FileListResult, error) {
	startTime := time.Now()
	var files []File

	filesByExtensionMap := make(map[string]*FilesByExtension)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if x, found := filesByExtensionMap[ext]; found {
				x.Files = append(x.Files, File{path, info})
			} else {
				var files []File
				files = append(files, File{path, info})
				filesByExtensionMap[ext] = &FilesByExtension{Extension: ext, Files: files}
			}
		}
		return nil
	})
	endTime := time.Now()

	if err != nil {
		return nil, err
	}

	return &FileListResult{
		filesByExtensionMap,
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
	fileScanner.Logging.Stdout.Printf("Scanned files with %d different extensions in %d ms", len(fileListResult.FilesByExtensionMap), duration)

	stat := NewStat(fileListResult.FilesByExtensionMap, fileScanner.Logging)

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

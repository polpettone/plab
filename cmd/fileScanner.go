package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

//TODO: check permissions before scanning

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

func NewStat(files map[string] *FilesByExtension) Stat {
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

func(fileScanner FileScanner) list(path string) {
	fileScanner.Logging.Stdout.Printf("List Files")

	fileListResult, err := listFiles(path)

	if err != nil {
		fileScanner.Logging.ErrorLog.Printf("%v", err)
		return
	}

	duration := fileListResult.EndTime.Sub(fileListResult.StartTime).Milliseconds()
	fileScanner.Logging.Stdout.Printf("Scanned files with %d different extensions in %d ms", len(fileListResult.FilesByExtensionMap), duration)

	stat := NewStat(fileListResult.FilesByExtensionMap)

	fileScanner.Logging.Stdout.Printf("Created Stats")

	for _, filesByExtension := range stat.FilesByExtension {
		fileScanner.Logging.InfoLog.Printf("%s %d", filesByExtension.Extension, len(filesByExtension.Files))
	}
}

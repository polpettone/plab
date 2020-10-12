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
	StartTime           time.Time
	EndTime             time.Time
}

type Stat struct {
	FilesByExtensionSlice []*FilesByExtension
	FileCount             int
}

type FilesByExtension struct {
	Extension string
	Files     []File
}

func NewStat(files map[string]*FilesByExtension) Stat {
	var l []*FilesByExtension

	for _, x := range files {
		l = append(l, x)
	}

	sort.Slice(l, func(i, j int) bool {
		return len(l[i].Files) > len(l[j].Files)
	})

	var count int
	for _, f := range files {
		count += len(f.Files)
	}

	return Stat{
		FilesByExtensionSlice: l,
		FileCount:             count,
	}
}

func listFiles(path string, logging *Logging) (*FileListResult, error) {
	startTime := time.Now()

	filesByExtensionMap := make(map[string]*FilesByExtension)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info != nil {
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
		} else {
			logging.DebugLog.Printf("info is nil for path %s error: %v", path, err)
		}
		return nil
	})
	endTime := time.Now()

	if err != nil {
		return nil, err
	}

	return &FileListResult{
		filesByExtensionMap,
		startTime,
		endTime,
	}, nil
}

func (fileScanner FileScanner) list(path string) {
	fileScanner.Logging.Stdout.Printf("List Files")

	fileListResult, err := listFiles(path, fileScanner.Logging)

	if err != nil {
		fileScanner.Logging.ErrorLog.Printf("%v", err)
		return
	}

	duration := fileListResult.EndTime.Sub(fileListResult.StartTime).Milliseconds()
	fileScanner.Logging.Stdout.Printf("Scanned files with %d different extensions in %d ms", len(fileListResult.FilesByExtensionMap), duration)

	stat := NewStat(fileListResult.FilesByExtensionMap)

	fileScanner.Logging.Stdout.Printf("total file count: %d", stat.FileCount)
	for _, filesByExtension := range stat.FilesByExtensionSlice {
		fileScanner.Logging.Stdout.Printf("%s %d", filesByExtension.Extension, len(filesByExtension.Files))
	}

	largestFiles := findNthLargesFiles(fileListResult)

	for _, f := range largestFiles {
		fileScanner.Logging.Stdout.Printf("%v %d %s", f.info.ModTime(), f.info.Size() / (1000 * 1000), f.path)
	}
}

func findNthLargesFiles(fileListResult *FileListResult) []File {
	var files []File
	for _, f := range fileListResult.FilesByExtensionMap {
		files = append(files, f.Files...)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].info.Size() > files[j].info.Size()
	})

	return files[:10]
}

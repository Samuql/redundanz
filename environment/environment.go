package environment

import (
	"log"
	"os"
	"path/filepath"
)

// return all folders in the current directory
// currently []string as i dont know any other better DS
// also get the name of the current folder
func GetFolderSelection() []string {
	folders := []string{GetWd()}
	files, _ := os.ReadDir(".")
	for _, f := range files {
		if f.IsDir() {
			folders = append(folders, f.Name())
		}
	}
	return folders
}

func GetWd() string {
	wd, err := os.Getwd()
	if err != nil {
		// TODO implement error handling with different logging levels
		panic(err)
	}
	return wd
}

// should later be implemented as a reentrant function => multiple calls of getContent to collect faster?

// TODO recursively get all files in a folder and its subfolders
// TODO no string return
// TODO handle multiple params

// // Gets all File names excluding the name of the selected folder
// func GetDirContent(dir string, logger *log.Logger) []string {
// 	files, err := os.ReadDir(dir)
// 	if err != nil {
// 		// TODO implement error handling with different logging levels
// 		panic(err)
// 	}
// 	fileNames := make([]string, 0, len(files))
// 	for _, f := range files {
// 		log.Printf("Found %s", f.Name())
// 		fileNames = append(fileNames, f.Name())
// 	}
// 	log.Printf("Found %d files in %s", len(fileNames), dir)
// 	return fileNames
// }

// OLDER
// // Gets all File names excluding the name of the selected folder
// func GetFilesInDir(dir string, logger *log.Logger) []os.DirEntry {
// 	files, err := os.ReadDir(dir)
// 	if err != nil {
// 		// TODO implement error handling with different logging levels
// 		panic(err)
// 	}
// 	files_sl := make([]os.DirEntry, 0, len(files))
// 	// fmt.Printf("\n\n== Folder %s/ ==\n", dir)
// 	for _, f := range files {
// 		// log.Printf("Found %s", f.Name())
// 		if f.IsDir() {
// 			continue // TODO get content
// 		}

// 		files_sl = append(files_sl, f)
// 	}
// 	// log.Printf("Found %d files in %s", len(fileNames), dir)
// 	return files_sl
// }

func GetFilesInDir(dir string, logger *log.Logger, filesChan chan<- []os.DirEntry) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		// TODO implement error handling with different logging levels
		panic(err)
	}
	filesSlice := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			// Call GetFilesInDir recursively for directories
			subDir := filepath.Join(dir, entry.Name())
			GetFilesInDir(subDir, logger, filesChan)
		} else {
			filesSlice = append(filesSlice, entry)
		}
	}
	filesChan <- filesSlice
}

func BuildEnvironment() {
	// TODO
}

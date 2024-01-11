package main

import (
	"log"
	"os"
)

// return all folders in the current directory
// currently []string as i dont know any other better DS
// also get the name of the current folder
func getSelection() []string {
	folders := []string{getWd()}
	files, _ := os.ReadDir(".")
	for _, f := range files {
		if f.IsDir() {
			folders = append(folders, f.Name())
		}
	}
	return folders
}

func getWd() string {
	wd, err := os.Getwd()
	if err != nil {
		// TODO implement error handling with different logging levels
		panic(err)
	}
	return wd
}

// should later be implemented as a reentrant function => multiple calls of getContent to collect faster?
// recursively get all files in a folder and its subfolders

// Gets all File names excluding the name of the selected folder
// TODO no string return
// TODO handle multiple params
func getFileNames(dir string, verbose bool) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		// TODO implement error handling with different logging levels
		panic(err)
	}
	fileNames := make([]string, 0, len(files))
	for _, f := range files {
		log.Printf("Found %s", f.Name())
		fileNames = append(fileNames, f.Name())
	}
	log.Printf("Found %d files in %s", len(fileNames), dir)
	return fileNames
}

package fileutils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func IsDirectoryWritable(path string) bool {
	file, err := ioutil.TempFile(path, "")
	if err != nil {
		return false
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	return true
}

func copyFile(sourceFile string, destinationFile string) {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Fatalf("An error occurred while reading the file %s:\n%v", sourceFile, err)
	}

	if err := ioutil.WriteFile(destinationFile, data, os.ModeDevice); err != nil {
		log.Fatalf("An error occurred while writing the file %s:\n%v", destinationFile, err)
	}
}

func copyTreeInternal(sourceDirectory string, destinationDirectory string, templateFiles *[]string) {
	fileInfos, err := ioutil.ReadDir(sourceDirectory)
	if err != nil {
		log.Fatalf("An error occurred while copying %s into %s:\n%v", sourceDirectory, destinationDirectory, err)
	}

	for _, fileInfo := range fileInfos {
		sourcePath := filepath.Join(sourceDirectory, fileInfo.Name())
		destinationPath := filepath.Join(destinationDirectory, fileInfo.Name())
		if fileInfo.IsDir() {
			if err := os.Mkdir(destinationPath, os.ModeDir); err != nil {
				log.Fatalf("An error occurred while creating the directory %s:\n%v", destinationPath, err)
			}
			copyTreeInternal(sourcePath, destinationPath, templateFiles)
		} else {
			copyFile(sourcePath, destinationPath)
			if filepath.Ext(fileInfo.Name())  == ".tmpl" {
				*templateFiles = append(*templateFiles, destinationPath)
			}
		}
	}
}

func CopyTree(sourceDirectory string, destinationDirectory string) []string {
	var templateFiles []string

	copyTreeInternal(sourceDirectory, destinationDirectory, &templateFiles)

	return templateFiles
}

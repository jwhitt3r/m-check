// Package directory manages the interactions of the Filesystem during
// the applications runtime.
package directory

import (
	"fmt"
	"io"
	"log"
	"os"
)

// CreateDirectory creates a new directory to store the Github Repository
// documentation within, currently it is /docs/<owner>/<RepoName>/
func CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 1)
		return nil
	}
	return err
}

// GetFilePathTemplate formats a filepath to be used for the creation of a new file.
func GetFilePathTemplate(base string, owner string, repoName string) string {
	filePathTemplate := fmt.Sprintf("./%s/%s/%s/", base, owner, repoName)
	return filePathTemplate
}

// OutputToFile creates a file to store the link and the responding status code or error.
func OutputToFile(path string, val string) error {
	f, err := os.OpenFile(path+"output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Failed to create output file: %v\n", err)
	}

	defer f.Close()

	_, err = io.WriteString(f, val+"\n")
	if err != nil {
		log.Printf("Failed to print io line: %v", err)
	}
	return nil
}

package directory

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Creates a new directory to store the Github Repository
// documentation within, currently it is /docs/<owner>/<RepoName>/
func CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 1)
		return nil
	}
	return err
}

func GetFilePathTemplate(owner string, repoName string) string {
	filePathTemplate := fmt.Sprintf("./docs/%s/%s/", owner, repoName)
	return filePathTemplate
}

func OutputToFile(path string, link string, code string) error {
	f, err := os.OpenFile(path+"output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Failed to create output file: %v\n", err)
	}

	defer f.Close()

	_, err = io.WriteString(f, link+" - "+code+"\n")

	return nil
}

package directory

import (
	"fmt"
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

func GetFileBaseTemplate() string {
	fileBaseTemplate := "./docs/"
	return fileBaseTemplate
}

func GetFilePathTemplate(owner string, repoName string) string {
	filePathTemplate := fmt.Sprintf("./docs/%s/%s/", owner, repoName)
	return filePathTemplate
}

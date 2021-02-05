package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jwhitt3r/gondola/internal/platform/repo"
	"github.com/jwhitt3r/gondola/internal/platform/urlchecker"
)

func main() {
	fmt.Println("[+] Welcome to Gondola [+]")
	owner, reponame, token := "jwhitt3r", "test_repo", ""
	client := http.Client{Timeout: 5 * time.Second}
	checker := urlchecker.NewURLChecker(client)

	myRepo := repo.NewGithubConnection(owner, reponame, token)

	myRepo.GetGithubContents(context.Background(), directory.GetFileBaseTemplate())

	fmt.Println("[+] Saving All Documentation Found")
	for _, fileURL := range myRepo.FilesURL {
		myRepo.Fetch(fileURL)
	}

	myRepo.GetFileNames()

	err := myRepo.Parse()
	if err != nil {
		log.Fatalf("Failed to parse documentation")
	}

	fmt.Println("[+] Checking Connectivty of Markdown Links")
	for _, link := range myRepo.Links {
		checker.URLCheck(fmt.Sprintf(filePathTemplate, myRepo.Owner, myRepo.RepoName), link)
	}

}

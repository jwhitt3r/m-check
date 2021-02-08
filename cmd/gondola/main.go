package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jwhitt3r/gondola/internal/platform/directory"
	"github.com/jwhitt3r/gondola/internal/repo"
	"github.com/jwhitt3r/gondola/internal/urlcheck"
)

const fileBaseTemplate = "./docs/"

var (
	o = flag.String("o", "", "Used to specify the owner of the repository.")
	r = flag.String("r", "", "Used to specify the Repository that you would like to search in.")
	t = flag.String("t", "", "Used to specify Your GitHub Personal Token.")
	l = flag.Bool("l", false, "Used to specify a local scan, this indicates that you have already downloaded the documentation.")
)

var usage = `Usage: Gondola [options...]

Options:
	-o Owner of the repository you would like to search.
	-r Repository that you would like to search in.

Optional:
	-t Your GitHub Personal Token if you would like to have a higher level of searchers.
	-l Indicates that there is a local copy of the documentation already downloaded.
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(usage))
	}

	flag.Parse()
	if flag.NFlag() < 2 {
		usageAndExit("")
	}

	owner := *o
	reponame := *r
	token := *t
	local := *l

	if owner == "" {
		usageAndExit(fmt.Sprintf("The repository owner has not been set"))
	}

	if reponame == "" {
		usageAndExit(fmt.Sprintf("The repository name has not been set"))
	}

	myRepo := repo.NewRepository(owner, reponame, token)
	client := http.Client{Timeout: 5 * time.Second}
	checker := urlcheck.NewURLChecker(client)
	if local == false {
		myRepo.NewGithubConnection()
		myRepo.GetGithubContents(context.Background(), fileBaseTemplate)

		fmt.Println("[+] Saving All Documentation Found")
		for _, fileURL := range myRepo.FilesURL {
			myRepo.Fetch(fileURL)
		}

	}
	myRepo.GetFileNames()

	err := myRepo.Parse()
	if err != nil {
		log.Fatalf("Failed to parse documentation")
	}

	fmt.Println("[+] Checking Connectivty of Markdown Links")
	response, err := checker.URLCheck(myRepo.Links)
	for key, val := range response {
		directory.OutputToFile(directory.GetFilePathTemplate(myRepo.Owner, myRepo.RepoName), key, val)
	}

}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	os.Exit(1)
}

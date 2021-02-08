// Gondola is a Markdown URL checker that is to verify
// links that are typically found within the scope of
// Github documentation directories.
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

// Used to generate the base location to store the downloaded documentation
// and output file.
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

Example For Downloading Content: ./gondola -o jwhitt3r -r gondola -t 12345678975336985

Example For Working On A Local Copy: ./gondola -o jwhitt3r -r gondola -l
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
	checker := urlcheck.NewURLCheck(client)
	if local == false {
		myRepo.NewGithubConnection()
		myRepo.GetGithubContents(context.Background(), fileBaseTemplate)

		fmt.Println("[+] Saving All Documentation Found")
		err := directory.CreateDirectory(directory.GetFilePathTemplate(myRepo.Owner, myRepo.RepoName))
		if err != nil {
			log.Fatalf("An error occured while making a new directory: %v\n", err)
		}
		myRepo.Fetch(myRepo.FilesURL)

	}
	myRepo.GetFileNames()

	links := myRepo.ParseBatch()

	fmt.Println("[+] Checking Connectivty of Markdown Links")
	webConnectionResponse := checker.URLCheckBatch(links)

	for _, val := range webConnectionResponse {
		directory.OutputToFile(directory.GetFilePathTemplate(myRepo.Owner, myRepo.RepoName), val)
	}

}

// A simple function to present the usage of flags when running the command.
// This is typically called when there are not enough flags have been passed at runtime.
func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	os.Exit(1)
}

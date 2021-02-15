// m-check is a Markdown URL checker that is to verify
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

	"github.com/jwhitt3r/m-check/internal/platform/directory"
	"github.com/jwhitt3r/m-check/internal/repo"
	"github.com/jwhitt3r/m-check/internal/urlcheck"
)

var (
	o = flag.String("o", "", "Used to specify the owner of the repository.")
	r = flag.String("r", "", "Used to specify the Repository that you would like to search in.")
	t = flag.String("t", "", "Used to specify Your GitHub Personal Token.")
	b = flag.String("b", "./docs", "Used to specify the Base Path to save your documents")
	p = flag.String("p", "docs", "Used to specify the remote documentation location")

	l = flag.Bool("l", false, "Used to specify a local scan, this indicates that you have already downloaded the documentation.")
)

var usage = `Usage: m-check [mandatory...] [options...]

Mandatory:
	-o Owner of the repository you would like to search.
	-r Repository that you would like to search in.

Optional:
	-t Your GitHub Personal Token if you would like to have a higher level of searchers.
	-l Indicates that there is a local copy of the documentation already downloaded.
	-b Used to specify the Base Path to save your documents, by default this will be ./docs.
	-p Used to specify the remote documentation location, by default this will be "docs".

Output:
	The output of the check will be stored within the specified basepath, under the name output.txt

Examples:
	Example For Downloading Content: ./m-check -o jwhitt3r -r m-check -t 12345678975336985

	Example For Working On A Local Copy: ./m-check -o jwhitt3r -r m-check -l

	Example For Saving To Non-Default Destination: ./m-check -o jwhitt3r -r m-check -b ./tmp

	Example For Non-Default Remote Directory ./m-check -o jwhitt3r -r test_repo -p "documentation"
`

func main() {
	var FilesDownloadURL []string
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
	basepath := *b
	remotepath := *p

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

		fmt.Println("[+] Finding Repository")
		myRepo.NewGithubConnection()

		myRepo.GetGithubContents(context.Background(), remotepath, &FilesDownloadURL)

		fmt.Println("[+] Saving All Documentation Found")
		err := directory.CreateDirectory(directory.GetFilePathTemplate(basepath, myRepo.Owner, myRepo.RepoName))
		if err != nil {
			log.Fatalf("An error occured while making a new directory: %v\n", err)
		}
		myRepo.FetchAndCreate(basepath, FilesDownloadURL)

	}

	fmt.Println("[+] Gathering Filenames")
	files := myRepo.GetFileNames(basepath)

	links := myRepo.ParseBatch(basepath, files)

	fmt.Println("[+] Checking Connectivty Of Markdown Links")
	webConnectionResponse := checker.URLCheckBatch(links)

	fmt.Printf("[+] Findings Are Saved To %s/output.txt", basepath)
	for _, val := range webConnectionResponse {
		directory.OutputToFile(directory.GetFilePathTemplate(basepath, myRepo.Owner, myRepo.RepoName), val)
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

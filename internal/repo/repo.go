package repo

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jwhitt3r/gondola/internal/platform/directory"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// Repository holds all the key information for managing a Github repository.
type Repository struct {
	// Owner is the owner of the repository we are evaluating e.g., jwhitt3r.
	Owner string
	// RepoName is the repository that is going to be downloaded, e.g., Gondola.
	RepoName string
	// Files holds the names of all downloaded documentation.
	files []string
	// FilesURL are the URLS for the fetch function to download against.
	// This data is gathered by the GetGitHub Contents.
	FilesURL []string
	// Token is the personal token used to authenticate with the Github server.
	// By supplying a token, a user is allowed more requests to the Github server.
	token string
	// Links holds all the URLS that have been gathered from the the documents that
	// have been downloaded from the Documentation folder of a repository
	Links  []string
	client *github.Client
}

// GithubContents recursively looks through any directory within the Documentation folder
// of a repository and appends the FilesURL of a Markdown file to a slice of strings to be
// downloaded later.
func (r *Repository) GetGithubContents(ctx context.Context, path string) error {
	_, dirContents, _, err := r.client.Repositories.GetContents(ctx, r.Owner, r.RepoName, path, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	for _, element := range dirContents {
		switch element.GetType() {
		case "file":
			if filepath.Ext(element.GetName()) == ".md" {
				r.FilesURL = append(r.FilesURL, element.GetDownloadURL())
			}
		case "dir":
			r.GetGithubContents(ctx, element.GetPath())
		}
	}
	return nil
}

// NewRepoistory wraps the creation of a Repository type
func NewRepository(owner string, reponame string, token string) *Repository {
	r := Repository{
		Owner:    owner,
		RepoName: reponame,
		token:    token,
	}
	return &r
}

// GitHubConnection creates a connection to GitHub through a personal access token
// this increases the number of times you can connect to a repository
func (r *Repository) NewGithubConnection() {
	fmt.Println("[+] Finding Repository")
	ctx := context.Background()
	if r.token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: r.token},
		)
		tc := oauth2.NewClient(ctx, ts)

		r.client = github.NewClient(tc)
		return
	}
	r.client = github.NewClient(nil)

}

// FetchAndCreate will download all the files that have been
// collected by the GithubContents function, and save them into
// the local repository.
func (r *Repository) FetchAndCreate(fileURLS []string) error {

	for _, fileURL := range fileURLS {
		resp, err := http.Get(fileURL)
		if err != nil {
			log.Printf("Failed to fetch URL: %v\n", err)
		}

		u, err := url.Parse(fileURL)
		if err != nil {
			log.Printf("Failed to parse URL: %v\n", err)
		}
		path := strings.ReplaceAll(u.Path, "/", ".")
		pathFirstIndex := strings.Index(path, ".docs")

		f, err := os.Create(directory.GetFilePathTemplate(r.Owner, r.RepoName) + path[pathFirstIndex+6:])
		if err != nil {
			log.Fatalf("Failed to create file: %v\n", err)
		}
		defer f.Close()

		_, err = io.Copy(f, resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatalf("Failed to copy contents to file: %v\n", err)
		}
	}

	return nil

}

// GetFileNames gathers all the downloaded files found within the docs
// directory and stores them into the Files Slice.
func (r *Repository) GetFileNames() error {

	fmt.Println("[+] Gathering Filenames")
	files, err := ioutil.ReadDir(directory.GetFilePathTemplate(r.Owner, r.RepoName))

	if err != nil {
		log.Fatalf("Could not read files from directory: %v\n", err)
	}

	for _, fileName := range files {
		r.files = append(r.files, fileName.Name())
	}

	return nil
}

// Parse traverses a markdown file that has been downloaded within the
// documentation folder within the repository, and compares a regular
// expression to find any possible links within the documentation.
func (r *Repository) Parse() error {
	fmt.Println("[+] Parsing All Markdown Documentation")
	// The Regex will aim to locate any address that has the following structure:
	// https://github.com/jwhitt3r. An example of this would be within a markdown
	// file as: [Jwhitt3rs GitHub](https://github.com/jwhitt3r) or
	// file as: [Jwhitt3rs GitHub]("https://github.com/jwhitt3r")
	markdownURL := regexp.MustCompile(`https?://[^()]+?[^)"]+`)
	for _, fileName := range r.files {
		if filepath.Ext(fileName) == ".md" {
			f, err := os.Open(directory.GetFilePathTemplate(r.Owner, r.RepoName) + fileName)
			if err != nil {
				log.Fatalf("Failed to open file: %v\n", err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)

			for scanner.Scan() {

				submatchall := markdownURL.FindAllString(scanner.Text(), -1)

				for _, element := range submatchall {
					r.Links = append(r.Links, strings.TrimSpace(element))
				}

			}
		}
	}
	return nil
}

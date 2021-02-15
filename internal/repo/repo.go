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
	"sync"

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
	// Token is the personal token used to authenticate with the Github server.
	// By supplying a token, a user is allowed more requests to the Github server.
	token string
	// Links holds all the URLS that have been gathered from the the documents that
	// have been downloaded from the Documentation folder of a repository
	client *github.Client
}

// GithubContents recursively looks through any directory within the Documentation folder
// of a repository and appends the FilesURL of a Markdown file to a slice of strings to be
// downloaded later.
func (r *Repository) GetGithubContents(ctx context.Context, path string, filesDownloadURL *[]string) {

	_, dirContents, _, err := r.client.Repositories.GetContents(ctx, r.Owner, r.RepoName, path, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	for _, element := range dirContents {
		switch element.GetType() {
		case "file":
			if filepath.Ext(element.GetName()) == ".md" {
				*filesDownloadURL = append(*filesDownloadURL, element.GetDownloadURL())
			}
		case "dir":
			r.GetGithubContents(ctx, element.GetPath(), filesDownloadURL)
		}
	}

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

// GitHubConnection creates a connection to GitHub with or without a
// personal access token. However, with a personal token this increases
// the number of times you can connect to a repository.
func (r *Repository) NewGithubConnection() {
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
func (r *Repository) FetchAndCreate(basepath string, fileURLS []string) error {

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

		f, err := os.Create(directory.GetFilePathTemplate(basepath, r.Owner, r.RepoName) + path[pathFirstIndex+6:])
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
func (r *Repository) GetFileNames(basepath string) []string {
	var f []string
	files, err := ioutil.ReadDir(directory.GetFilePathTemplate(basepath, r.Owner, r.RepoName))

	if err != nil {
		log.Fatalf("Could not read files from directory: %v\n", err)
	}

	for _, fileName := range files {
		f = append(f, fileName.Name())
	}

	return f
}

// Parse traverses a markdown file that has been downloaded within the
// documentation folder within the repository, and compares a regular
// expression to find any possible links within the documentation.
func (r *Repository) Parse(f io.Reader) []string {

	var links []string
	// The Regex will aim to locate any address that has the following structure:
	// https://github.com/jwhitt3r. An example of this would be within a markdown
	// file as: [Jwhitt3rs GitHub](https://github.com/jwhitt3r) or
	// file as: [Jwhitt3rs GitHub]("https://github.com/jwhitt3r")
	markdownURL := regexp.MustCompile(`https?://[^()]+?[^)"]+`)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		submatchall := markdownURL.FindAllString(scanner.Text(), -1)

		for _, element := range submatchall {
			links = append(links, strings.TrimSpace(element))
		}

	}

	return links
}

// ParseFileHandler, will generate a file handler, which is then passed to the parse
// method to be analysed. This allows for the seperation of duties between the parser
// and the handling of files. This function will return the links that have been gathered
// from the parsed file.
func (r *Repository) ParseFileHandler(basepath string, fileName string) []string {
	var links []string
	if filepath.Ext(fileName) == ".md" {
		f, err := os.Open(directory.GetFilePathTemplate(basepath, r.Owner, r.RepoName) + fileName)
		if err != nil {
			log.Fatalf("Failed to open file: %v\n", err)
		}
		defer f.Close()
		links = r.Parse(f)

	}
	return links
}

// ParseBatch wraps a concurrent method for parsing a file
// which the outcome is then appended to a slice of strings,
// to be passed to the URLCheckBatch function.
func (r *Repository) ParseBatch(basepath string, files []string) []string {
	ch := make(chan []string, len(files))
	var links []string
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, fileName := range files {
		go func(fileName string) {
			ch <- r.ParseFileHandler(basepath, fileName)
			wg.Done()
		}(fileName)

	}
	wg.Wait()
	close(ch)
	for value := range ch {
		links = append(links, value...)
	}
	return links
}

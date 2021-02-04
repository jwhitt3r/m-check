package main

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
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// Repository holds all the key information to review a GitHub repo
type Repository struct {
	Owner    string
	RepoName string
	FilesURL []string
	Files    []string
	Links    []string
	Token    string
}

// GithubContents recursively looks through any directory within the Documentation folder
// of a repository and appends the FilesURL to a slice of strings to be downloaded later.
func (r *Repository) githubContents(client *github.Client, ctx context.Context, path string) {
	_, dirContents, _, err := client.Repositories.GetContents(ctx, r.Owner, r.RepoName, path, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	for _, element := range dirContents {
		switch element.GetType() {
		case "file":
			r.FilesURL = append(r.FilesURL, element.GetDownloadURL())
		case "dir":
			r.githubContents(client, ctx, element.GetPath())
		}
	}
}

// GitHubConnection creates a connection to GitHub through a personal access token
// this increases the number of times you can connect to a repository
func (r *Repository) GitHubConnection() error {

	path := "docs"

	fmt.Println("[+] Finding Repository")
	ctx := context.Background()
	if r.Token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: r.Token},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)
		r.githubContents(client, ctx, path)
	} else {
		client := github.NewClient(nil)
		r.githubContents(client, ctx, path)
	}

	return nil
}

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

// Fetch will download all the files that have been collected by the GithubContents
// function, and save them into the local repository.
func (r *Repository) Fetch(fileURL string) error {
	filePath := "./docs/" + r.Owner + "/" + r.RepoName
	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch URL: %v\n", err)
	}

	err = CreateDirectory(filePath)
	if err != nil {
		log.Fatalf("An error occured while making a new directory: %v\n", err)
	}

	f, err := os.Create(filePath + "/" + url.QueryEscape(path.Base(fileURL)))
	if err != nil {
		log.Fatalf("Failed to create file: %v\n", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalf("Failed to copy contents to file: %v\n", err)
	}

	return nil

}

// URLCheck makes a connection to the list of URLS found within the
// Markdown documentation, and provides the HTTP status_code to be
// acted upon
func (r *Repository) URLCheck(link string) error {
	f, err := os.OpenFile("./docs/"+r.Owner+"/"+r.RepoName+"/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Failed to create output file: %v\n", err)
	}

	defer f.Close()

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(link)
	if err != nil {
		io.WriteString(f, link+" - "+"Broken Link"+"\n")

		log.Printf("Failed to connect to page: %v\n", err)
		return err
	}
	_, err = io.WriteString(f, link+" - "+strconv.Itoa(resp.StatusCode)+"\n")
	resp.Body.Close()

	return nil

}

// GetFileNames gathers all the downloaded files found within the docs
// directory and stores them into the Files Slice.
func (r *Repository) GetFileNames() error {

	fmt.Println("[+] Gathering Filenames")
	files, err := ioutil.ReadDir("./docs/" + r.Owner + "/" + r.RepoName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read files from directory: %v\n", err)
	}

	for _, fileName := range files {
		r.Files = append(r.Files, fileName.Name())
	}

	return nil
}

// Parse traverses a markdown file that has been downloaded within the
// documentation folder within the repository, and compares a regular
// expression to find any possible links within the documentation.
func (r *Repository) Parse() error {
	fmt.Println("[+] Parsing All Markdown Documentation")
	markdownURL := regexp.MustCompile(`https?://[^()]+?[^"]+`)
	for _, fileName := range r.Files {
		f, err := os.Open("./docs/" + r.Owner + "/" + r.RepoName + "/" + fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
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
	return nil
}

func main() {
	fmt.Println("[+] Welcome to Gondola [+]")

	myRepo := Repository{
		Owner:    "jwhitt3r",
		RepoName: "test_repo",
		Token:    "Token Here",
	}

	myRepo.GitHubConnection()

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
		err := myRepo.URLCheck(link)
		if err != nil {
			log.Fatalf("Failed to do URL Check: %v\n", err)
		}
	}

}

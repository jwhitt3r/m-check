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
)

type Repository struct {
	Owner    string
	RepoName string
	FilesURL []string
	Files    []string
	Links    []string
}

// CollectDocs goes through the chosen Repository and pulls the
// Contents of the docs folder
func (r *Repository) CollectDocs() error {
	fmt.Println("[+] Finding Repository")
	client := github.NewClient(nil)

	_, dirContent, _, err := client.Repositories.GetContents(context.Background(), r.Owner, r.RepoName, "docs", nil)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
	for i := range dirContent {
		if dirContent[i].GetType() == "file" {
			r.FilesURL = append(r.FilesURL, dirContent[i].GetDownloadURL())
		} else if dirContent[i].GetType() == "dir" {
			_, dirContent, _, err := client.Repositories.GetContents(context.Background(), r.Owner, r.RepoName, "docs"+"/"+dirContent[i].GetName(), nil)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			r.FilesURL = append(r.FilesURL, dirContent[i].GetDownloadURL())
		}

	}

	return nil
}

// Fetch will download all the files that have been collected by the GetContents
// function, and save them into the local repository
func (r *Repository) Fetch(fileURL string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch URL: %v\n", err)
	}

	// Create a directory
	if _, err := os.Stat("./docs/" + r.Owner + "/" + r.RepoName); os.IsNotExist(err) {
		fmt.Println("[+] Creating a Docs directory")
		os.MkdirAll("./docs/"+r.Owner+"/"+r.RepoName, 1)
	}

	// Download the files
	f, err := os.Create("./docs/" + r.Owner + "/" + r.RepoName + "/" + url.QueryEscape(path.Base(fileURL)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to copy contents to file: %v\n", err)
	}

	return nil

}

func (r *Repository) URLCheck(link string) error {

	f, err := os.OpenFile("./docs/"+r.Owner+"/"+r.RepoName+"/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Failed to create output file: %v\n", err)
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
	fmt.Println("Welcome to Gondola")
	myRepo := Repository{
		Owner:    "jwhitt3r",
		RepoName: "test_repo",
	}

	myRepo.CollectDocs()

	fmt.Println("[+] Saving All Documentation Found")
	for _, fileURL := range myRepo.FilesURL {
		myRepo.Fetch(fileURL)
	}

	myRepo.GetFileNames()
	err := myRepo.Parse()
	if err != nil {
		fmt.Println("Failed to parse")
	}

	fmt.Println("[+] Checking Connectivty of Markdown Links")
	for _, link := range myRepo.Links {
		err := myRepo.URLCheck(link)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to do URL Check: %v\n", err)
		}
	}

}

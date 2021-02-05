package urlchecker

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type URLChecker struct {
	client http.Client
}

func NewURLChecker(client http.Client) *URLChecker {
	return &URLChecker{
		client: client,
	}
}

// URLCheck makes a connection to the list of URLS found within the
// Markdown documentation, and provides the HTTP status_code to be
// acted upon
func (u *URLChecker) URLCheck(path string, link string) error {
	f, err := os.OpenFile(path+"output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Failed to create output file: %v\n", err)
	}

	defer f.Close()

	resp, err := u.client.Get(link)
	if err != nil {
		io.WriteString(f, link+" - "+"Broken Link"+"\n")
		log.Printf("Failed to connect to page: %v\n", err)
		return err
	}
	_, err = io.WriteString(f, link+" - "+strconv.Itoa(resp.StatusCode)+"\n")
	defer resp.Body.Close()

	return nil

}

package urlcheck

import (
	"fmt"
	"net/http"
	"strconv"
)

// URLChecker represents a URL that is being used to verify the URI's status code
type URLChecker struct {
	// Client specifies a http.Client type to be used to make GET requests
	// within the URLCheck function.
	client http.Client
}

// NewURLChecker is a wrapper for the creation of a URLChecker type
// which returns the address of the newly created URLChecker type.
func NewURLCheck(client http.Client) *URLChecker {
	return &URLChecker{
		client: client,
	}
}

// URLCheck makes a connection to the list of URLS found within the
// Markdown documentation, and appends the link and status code to
// a slice of strings. This will be passed to the OutputFile function
// found in the /internal/platform/directory.go file.
func (u *URLChecker) URLCheck(links []string) ([]string, error) {
	var webConnectionResponse []string

	for _, link := range links {
		resp, err := u.client.Get(link)
		if err != nil {
			webConnectionResponse = append(webConnectionResponse, fmt.Sprintf("%s - Broken Link", link))
			continue
		}
		webConnectionResponse = append(webConnectionResponse, fmt.Sprintf("%s - %s", link, strconv.Itoa(resp.StatusCode)))
		defer resp.Body.Close()
	}
	return webConnectionResponse, nil
}

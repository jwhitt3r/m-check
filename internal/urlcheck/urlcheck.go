package urlcheck

import (
	"fmt"
	"net/http"
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

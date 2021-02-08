package urlcheck

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
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

// URLCheck makes a connection to a url found within the
// Markdown documentation and returns the formatted string
// to be appended to a list of links and status codes to
// be examined later on.
func (u *URLChecker) URLCheck(link string) string {
	resp, err := u.client.Get(link)
	if err != nil {
		return fmt.Sprintf("%s - Broken Link", link)
	}
	defer resp.Body.Close()
	return fmt.Sprintf("%s - %s", link, strconv.Itoa(resp.StatusCode))
}

// URLCheckBatch takes a list of urls and wraps a concurrent
// check of each url found within the documentation. The method
// then returns a slice of the outcomes to be saved to file.
func (u *URLChecker) URLCheckBatch(links []string) []string {
	var webConnectionResponse []string
	ch := make(chan string, len(links))
	var wg sync.WaitGroup
	wg.Add(len(links))
	for _, link := range links {
		go func(link string) {
			ch <- u.URLCheck(link)
			wg.Done()
		}(link)
	}
	wg.Wait()
	close(ch)
	for value := range ch {
		webConnectionResponse = append(webConnectionResponse, value)
	}
	return webConnectionResponse
}

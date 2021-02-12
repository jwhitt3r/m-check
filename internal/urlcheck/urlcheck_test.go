package urlcheck

import (
	"net"
	"net/http"
	"testing"
	"time"
)

const success = "\u2713"
const failure = "\u2717"

func TestURLCheckStatusCode(t *testing.T) {
	tt := []struct {
		url        string
		statusCode int
	}{
		{"http://github.com/jwhitt3r/", http.StatusOK},
		{"http://github.com/jwhitt3rasdasdasd3/", http.StatusNotFound},
	}
	t.Log("Given the need to fetch a URL")
	{
		for testID, test := range tt {
			t.Logf("Test %d:\tWhen checking %q for a status code %d", testID, test.url, test.statusCode)
			{
				resp, err := http.Get(test.url)
				if err != nil {
					t.Fatalf("\t%s\tTest %d: Should be able to make a Get call %v", failure, testID, err)
				}
				t.Logf("\t%s\tTest %d: Should be able to make a Get call.", success, testID)
				defer resp.Body.Close()
				if resp.StatusCode == test.statusCode {
					t.Logf("\t%s\tTest %d: Should receive a %d status code.", success, testID, test.statusCode)
				} else {
					t.Errorf("\t%s\tTest %d: Should receive a %d status code : %v", failure, testID, test.statusCode, resp.StatusCode)
				}
			}
		}
	}
}

func TestURLCheckBrokenLink(t *testing.T) {
	tt := []struct {
		url string
	}{
		{"https://askldjalskdjlkasdjlskadjlkas.com"},
	}
	client := http.Client{
		Timeout: 1 * time.Millisecond,
	}

	t.Log("Given a timeout connection to a URL")
	for testID, test := range tt {
		t.Logf("Test %d:\tWhen checking %q for a timeout a network error should be given", testID, test.url)
		{
			_, err := client.Get(test.url)
			_, ok := err.(net.Error)
			if ok == true {
				t.Logf("\t%s\tTest %d: Should receive a network error.", success, testID)
			} else {
				t.Logf("\t%s\tTest %d: Should receive a network error.", failure, testID)
			}

		}
	}
}

package repo

import (
	"regexp"
	"testing"
)

const success = "\u2713"
const failure = "\u2717"

func TestParse(t *testing.T) {
	tt := []struct {
		text string
	}{
		{"[Jwhitt3rs Github](https://github.com/jwhitt3r)"},
		{"[Jwhitt3rs Github](http://github.com/jwhitt3r)"},
	}

	markdownURL := regexp.MustCompile(`https?://[^()]+?[^)"]+`)
	t.Log("Given the need to check if a string is matched by a regualr expresssion")
	for testID, test := range tt {
		t.Logf("Test %d:\tWhen checking %s for a match", testID, test.text)
		submatch := markdownURL.MatchString(test.text)
		if submatch == true {
			t.Logf("\t%s\tTest %d:\t%s Should be matched by the regular expression", success, testID, test.text)
		} else {
			t.Fatalf("\t%s\tTest %d:\t%s Should be matched by the regular expression", failure, testID, test.text)
		}

	}
}

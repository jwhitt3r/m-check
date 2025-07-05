[![Go Report Card](https://goreportcard.com/badge/github.com/jwhitt3r/m-check)](https://goreportcard.com/report/github.com/jwhitt3r/m-check)
# m-check
m-check is a markdown parser aimed at reviewing links found within the documentation of Github repositories.

m-check, can work with both remote repositories and local repositories.

If a link is detected within a markdown file, a GET request will be made to establish if the connection is valid. Depending on the outcome, a HTTP Status will be provided, e.g., `200 , 400, 404, etc`. 

If however, the link could not be reached, for example, due to network error or timeout, a `Broken Link` will be added to the `output.txt` file.

It should be advised that any link that does not appear to work, should be manually investigated.

# Installation
To download m-check, simply run:
```
go get github.com/jwhitt3r/m-check
```

Finally, to build run the makefile:
```
make build
```

or alternatively:
```
go build cmd/m-check/m-check.go
```

# Usage
Below is a detailed breakdown of each flag, with working examples.

```
$ go run cmd/m-check/m-check.go
Usage: m-check [mandatory...] [options...]

Mandatory:
        -o Owner of the repository you would like to search.
        -r Repository that you would like to search in.

Optional:
        -t Your GitHub Personal Token if you would like to have a higher level of searchers.
        -l Indicates that there is a local copy of the documentation already downloaded.
        -b Used to specify the Base Path to save your documents, by default this will be ./docs.
        -p Used to specify the remote documentation location, by default this will be "docs".

Output:
        The output of the check will be stored within the specified basepath, under the name output.txt

Examples:
        Example For Downloading Content: ./m-check -o jwhitt3r -r m-check -t 12345678975336985

        Example For Working On A Local Copy: ./m-check -o jwhitt3r -r m-check -l

        Example For Saving To Non-Default Destination: ./m-check -o jwhitt3r -r m-check -b ./tmp

        Example For Non-Default Remote Directory ./m-check -o jwhitt3r -r test_repo -p "documentation"
```

# Thank You's and Inspirations
Thank you to [@mneverov](https://github.com/mneverov) for his mentorship through the development of this project!

Thank you to [@ardanlabs](https://github.com/ardanlabs) for the awesome training Ultimate Go!

Inspiration on how to handle flags was found at: [@rakyll](https://github.com/rakyll/hey)!

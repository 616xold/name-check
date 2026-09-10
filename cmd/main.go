package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/616xold/namecheck/bluesky"
	"github.com/616xold/namecheck/github"
)

type Checker interface {
	fmt.Stringer

	IsValid(string) bool
	IsAvailable(string) (bool, error)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "usage: namecheck <username>\n")
		os.Exit(1)
	}
	username := os.Args[1]
	checkers := []Checker{
		&github.GitHub{Client: http.DefaultClient},
		&bluesky.Bluesky{},
	}
	for _, checker := range checkers {
		if !checker.IsValid(username) {
			continue
		}
		avail, err := checker.IsAvailable(username)
		if err != nil {
			log.Fatal(err)
		}
		if !avail {
			continue
		}
		fmt.Printf("%q is valid and available on %s\n", username, checker)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/616xold/namecheck/bluesky"
	"github.com/616xold/namecheck/github"
)

func main() {
	username := "jub0bs"
	if !github.IsValid(username) {
		return
	}
	avail, err := github.IsAvailable(username)
	if err != nil {
		log.Fatal(err)
	}
	if !avail {
		return
	}
	fmt.Printf("%q is valid and available on GitHub\n", username)
	if !bluesky.IsValid(username) {
		return
	}
	avail, err = bluesky.IsAvailable(username)
	if err != nil {
		log.Fatal(err)
	}
	if !avail {
		return
	}
	fmt.Printf("%q is valid and available on Bluesky\n", username)
}

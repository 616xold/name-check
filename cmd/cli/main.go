package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/616xold/namecheck/github"
)

type Checker interface {
	fmt.Stringer

	IsValid(string) bool
	IsAvailable(string) (bool, error)
}

type Result struct {
	Platform  string
	Valid     bool
	Available bool
	Err       error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "usage: namecheck <username>\n")
		os.Exit(1)
	}
	username := os.Args[1]

	gh := github.GitHub{
		Client: http.DefaultClient,
	}
	checkers := make([]Checker, 16)

	for i := range checkers {
		checkers[i] = &gh
	}

	var wg sync.WaitGroup
	resultCh := make(chan Result)

	for _, checker := range checkers {
		wg.Add(1)
		go check(checker, username, &wg, resultCh)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []Result
	for result := range resultCh {
		results = append(results, result)
	}
	fmt.Println(results)
}

func check(checker Checker, username string, wg *sync.WaitGroup, resultCh chan Result) {
	defer wg.Done()

	result := Result{
		Platform: checker.String(),
	}

	result.Valid = checker.IsValid(username)

	if !result.Valid {
		resultCh <- result
		return
	}

	avail, err := checker.IsAvailable(username)
	if err != nil {
		result.Err = err
		resultCh <- result
		return
	}
	result.Available = avail
	resultCh <- result
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/616xold/namecheck/github"

	"github.com/jub0bs/cors"
)

type Checker interface {
	fmt.Stringer

	IsValid(string) bool
	IsAvailable(string) (bool, error)
}

type Result struct {
	Platform  string `json:platform`
	Valid     bool   `json:valid`
	Available bool   `json:available`
	Err       error
}

func handleCheck(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	response := struct {
		Username string   `json:"username"`
		Results  []Result `json:"results,omitempty`
	}{
		Username: username,
		Results:  results,
	}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(
			w,
			"failed to encode response ",
			http.StatusInternalServerError,
		)
		return
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /check", handleCheck)

	// instantiate a CORS middleware whose configuration suits your needs
	corsMw, err := cors.NewMiddleware(cors.Config{
		Origins: []string{"https://namecheck.jub0bs.dev"}, // TODO: adapt this!
	})
	if err != nil {
		log.Fatal(err)
	}

	// apply your CORS middleware to your HTTP request multiplexer
	handler := corsMw.Wrap(mux)

	// pass the result to ListenAndServe
	if err := http.ListenAndServe(":8080", handler); err != http.ErrServerClosed {
		log.Fatal(err)
	}
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

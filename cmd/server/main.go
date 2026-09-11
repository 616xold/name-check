package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"sync"

	"github.com/616xold/namecheck/github"
	"golang.org/x/net/context"

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

var checks = make(map[string]int)
var mu sync.Mutex

func handleCheck(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")

	mu.Lock()
	checks[username]++
	mu.Unlock()

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
		go check(context.TODO(), checker, username, &wg, resultCh)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []Result
	for result := range resultCh {
		if result.Err != nil {
			http.Error(
				w,
				"availability check failed",
				http.StatusInternalServerError,
			)
			return
		}
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

func handleStats(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	snapshot := maps.Clone(checks)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(
			w,
			"failed to encode stats",
			http.StatusInternalServerError,
		)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /check", handleCheck)
	mux.HandleFunc("GET /stats", handleStats)

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

func check(ctx context.Context, checker Checker, username string, wg *sync.WaitGroup, resultCh chan Result) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	result := Result{
		Platform: checker.String(),
	}

	result.Valid = checker.IsValid(username)

	if !result.Valid {
		select {
		case resultCh <- result:
		case <-ctx.Done():
			return
		}
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
	}

	avail, err := checker.IsAvailable(username)
	if err != nil {
		result.Err = err

		select {
		case resultCh <- result:
		case <-ctx.Done():
			return
		}
		return
	}

	result.Available = avail

	select {
	case resultCh <- result:
	case <-ctx.Done():
		return
	}
}

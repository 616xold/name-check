package github_test

import (
	"strings"
	"testing"

	"github.com/616xold/namecheck/github"
)

func TestIsValid(t *testing.T) {
	type TestCase struct {
		username string
		want     bool
	}
	testCases := map[string]TestCase{
		"contains two consecutive hyphens": {username: "jub0bs--on-GitHub", want: false},
		"starts with a hyphen":             {username: "-616xold", want: false},
		"ends with a hyphen":               {username: "jub0bs-", want: false},
		"too short":                        {username: "ab", want: false},
		"too long":                         {username: strings.Repeat("a", 40), want: false},
		"contains illegal chars":           {username: "jub&bs", want: false},
		"all good":                         {username: "jub0bs", want: true},
	}

	for desc, tc := range testCases {
		f := func(t *testing.T) {
			got := github.IsValid(tc.username)
			if got != tc.want {
				const tmpl = "github.IsValid(%q): got %t; want %t"
				t.Errorf(tmpl, tc.username, got, tc.want)
			}
		}
		t.Run(desc, f)
	}
}

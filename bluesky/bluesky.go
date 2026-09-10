package bluesky

import "github.com/616xold/namecheck"

type Bluesky struct {
	Client namecheck.Doer
}

func (b *Bluesky) IsValid(username string) bool { return false }

func (b *Bluesky) IsAvailable(username string) (bool, error) { return false, nil }

func (gh *Bluesky) String() string {
	return "Bluesky"
}

package bluesky

import "net/http"

type Bluesky struct {
	Client *http.Client
}

func (b *Bluesky) IsValid(username string) bool { return false }

func (b *Bluesky) IsAvailable(username string) (bool, error) { return false, nil }

func (gh *Bluesky) String() string {
	return "Bluesky"
}

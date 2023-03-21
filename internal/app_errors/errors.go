package app_errors

import "fmt"

// ErrorNoMoreLinks used when we ask for a new unvisited link but  already visited every link
type ErrorNoMoreLinks struct {
}

func (e ErrorNoMoreLinks) Error() string {
	return "No more links to visit"
}

// ErrorLinkAlreadyExists used when we try to add a link that already exists
type ErrorLinkAlreadyExists struct {
	Link string
}

func (e ErrorLinkAlreadyExists) Error() string {
	return fmt.Sprintf("the link %s already exists", e.Link)
}

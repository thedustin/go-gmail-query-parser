package filter

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Action interface {
	// Perform executes the action on the server.
	// Return true to call ""
	Perform(*imap.Message, *client.Client) error
}

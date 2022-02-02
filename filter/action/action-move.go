package action

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type ActionMove struct {
	dest string
}

func NewMove(dest string) ActionMove {
	return ActionMove{dest}
}

func (a ActionMove) Perform(msg *imap.Message, c *client.Client) error {
	seqset := new(imap.SeqSet)
	seqset.AddNum(msg.Uid)

	err := c.UidMove(seqset, a.dest)

	if err != nil {
		return fmt.Errorf("move message failed: %w", err)
	}

	return nil
}

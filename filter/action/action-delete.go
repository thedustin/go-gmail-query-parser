package action

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var (
	deleteFlagItem = imap.FormatFlagsOp(imap.AddFlags, true)
	deleteFlag     = []interface{}{imap.DeletedFlag}
)

type ActionDelete struct{}

func NewDelete() ActionDelete {
	return ActionDelete{}
}

func (a ActionDelete) Perform(msg *imap.Message, c *client.Client) error {
	seqset := new(imap.SeqSet)
	seqset.AddNum(msg.Uid)

	err := c.UidStore(seqset, deleteFlagItem, deleteFlag, nil)

	if err != nil {
		return fmt.Errorf("mark as deleted failed: %w", err)
	}

	// @todo: Could be called once
	err = c.Expunge(nil)

	if err != nil {
		return fmt.Errorf("expunge failed: %w", err)
	}

	return nil
}

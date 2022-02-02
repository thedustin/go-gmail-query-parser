package action

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var (
	seenFlagItem = imap.FormatFlagsOp(imap.AddFlags, true)
	seenFlag     = []interface{}{imap.SeenFlag}
)

type ActionMarkAsRead struct{}

func NewMarkAsRead() ActionMarkAsRead {
	return ActionMarkAsRead{}
}

func (a ActionMarkAsRead) Perform(msg *imap.Message, c *client.Client) error {
	seqset := new(imap.SeqSet)
	seqset.AddNum(msg.Uid)

	err := c.UidStore(seqset, seenFlagItem, seenFlag, nil)

	if err != nil {
		return fmt.Errorf("mark as seen/read failed: %w", err)
	}

	return nil
}

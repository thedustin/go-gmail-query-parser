package criteria

import (
	"fmt"
	"strings"

	"github.com/emersion/go-imap"
)

func addressesContains(addrs []*imap.Address, substr string) bool {
	for _, addr := range addrs {
		if addressContains(addr, substr) {
			return true
		}
	}

	return false
}

func addressContains(addr *imap.Address, substr string) bool {
	s := fmt.Sprintf("%s <%s>", addr.PersonalName, addr.Address())

	return strings.Contains(s, substr)
}

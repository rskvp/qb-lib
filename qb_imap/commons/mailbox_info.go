package commons

import qbc "github.com/rskvp/qb-core"

type MailboxInfo struct {
	// The mailbox attributes.
	Attributes []string `json:"attributes"`
	// The server's path separator.
	Delimiter string `json:"delimiter"`
	// The mailbox name.
	Name string `json:"name"`
}

func (instance *MailboxInfo) String() string {
	return qbc.JSON.Stringify(instance)
}

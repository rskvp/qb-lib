package lib_imap

import (
	"crypto/tls"
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_imap/commons"
)

type DriverImap struct {
	settings *commons.ImapClientSettings

	open bool

	_client *client.Client
}

func NewDriverImap(settings *commons.ImapClientSettings) (*DriverImap, error) {
	instance := new(DriverImap)
	instance.settings = settings

	return instance, nil
}

func (instance *DriverImap) Open() error {
	if nil != instance {
		c, err := instance.client()
		if nil != err {
			return err
		}

		if !instance.open {
			err = c.Login(instance.settings.Auth.User, instance.settings.Auth.Pass)
			if nil == err {
				instance.open = true
			}
			return err
		}
	}
	return nil
}

func (instance *DriverImap) Close() error {
	if nil != instance {
		instance.open = false
		c, err := instance.client()
		if nil != err {
			return err
		}
		return c.Logout()
	}
	return nil
}

func (instance *DriverImap) IsOpen() bool {
	if nil != instance && nil != instance._client {
		return instance.open
	}
	return false
}

func (instance *DriverImap) ListMailboxes() ([]*commons.MailboxInfo, error) {
	response := make([]*commons.MailboxInfo, 0)
	c, err := instance.client()
	if nil != err {
		return nil, err
	}

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	for m := range mailboxes {
		response = append(response, &commons.MailboxInfo{
			Attributes: m.Attributes,
			Delimiter:  m.Delimiter,
			Name:       m.Name,
		})
	}

	// wait for error
	if err = <-done; err != nil {
		return nil, err
	}
	return response, nil
}

func (instance *DriverImap) GetMailboxFlags(mailboxName string) ([]string, error) {
	_, box, err := instance.mailbox(mailboxName, true)
	if nil != err {
		return nil, err
	}
	return box.Flags, nil
}

func (instance *DriverImap) ReadMailbox(mailboxName string, onlyNew bool) ([]*commons.ImapMessage, error) {
	response := make([]*commons.ImapMessage, 0)
	messages, err := instance.readEnvelope(mailboxName, onlyNew)
	for _, message := range messages {
		m := new(commons.ImapMessage)
		parseErr := m.Parse(message)
		if nil == parseErr {
			response = append(response, m)
		}
	}
	return response, err
}

func (instance *DriverImap) ReadMessage(uid interface{}) (*commons.ImapMessage, error) {
	seqNum := uint32(qbc.Convert.ToInt(uid))
	im, err := instance.readFullMessage(seqNum)
	if nil != err {
		return nil, err
	}
	m := new(commons.ImapMessage)
	parseErr := m.Parse(im)
	return m, parseErr
}

func (instance *DriverImap) MarkMessageAsSeen(uid interface{}) error {
	seqNum := uint32(qbc.Convert.ToInt(uid))
	flags := []string{imap.SeenFlag}
	return instance.markMessage(seqNum, flags)
}

func (instance *DriverImap) MarkMessageAsDeleted(uid interface{}) error {
	seqNum := uint32(qbc.Convert.ToInt(uid))
	flags := []string{imap.DeletedFlag}
	return instance.markMessage(seqNum, flags)
}

func (instance *DriverImap) MarkMessageAsAnswered(uid interface{}) error {
	seqNum := uint32(qbc.Convert.ToInt(uid))
	flags := []string{imap.AnsweredFlag}
	return instance.markMessage(seqNum, flags)
}

func (instance *DriverImap) MarkMessageAsFlagged(uid interface{}) error {
	seqNum := uint32(qbc.Convert.ToInt(uid))
	flags := []string{imap.FlaggedFlag}
	return instance.markMessage(seqNum, flags)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DriverImap) mailbox(name string, readOnly bool) (*client.Client, *imap.MailboxStatus, error) {
	c, err := instance.client()
	if nil != err {
		return nil, nil, err
	}
	if len(name) == 0 {
		name = "INBOX"
	}
	box, err := c.Select(name, readOnly)
	if nil != err {
		return nil, nil, err
	}
	return c, box, nil
}

func (instance *DriverImap) client() (*client.Client, error) {
	if nil == instance._client {
		settings := instance.settings
		address := fmt.Sprintf("%s:%v", settings.Host, settings.Port)
		// TLS config
		var tlsConfig *tls.Config
		if settings.Tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         settings.Host,
			}
		} else {
			tlsConfig = nil
		}
		c, err := client.DialTLS(address, tlsConfig)
		if nil != err {
			return nil, err
		}
		instance._client = c
	}
	return instance._client, nil
}

func (instance *DriverImap) readEnvelope(mailboxName string, onlyNew bool) ([]*imap.Message, error) {
	fetchItems := []imap.FetchItem{imap.FetchEnvelope, imap.FetchInternalDate}
	c, box, err := instance.mailbox(mailboxName, true)
	if nil != err {
		return nil, err
	}
	if onlyNew {
		seqNums, err := c.Search(&imap.SearchCriteria{
			WithoutFlags: []string{imap.SeenFlag},
		})
		if nil != err {
			return nil, err
		}
		return fetchSeq(c, seqNums, fetchItems)
	} else {
		return fetchFrom(c, 1, box.Messages, fetchItems)
	}
	// return []*imap.Message{}, nil
}

func (instance *DriverImap) readFullMessage(seqNum uint32) (*imap.Message, error) {
	c, err := instance.client()
	if nil != err {
		return nil, err
	}

	data, err := fetchFrom(c, seqNum, seqNum, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody, imap.FetchRFC822, imap.FetchInternalDate})
	if nil != err {
		return nil, err
	}
	if len(data) == 1 {
		return data[0], err
	}
	return nil, err
}

func (instance *DriverImap) markMessage(seqNum uint32, flags []string) error {
	c, err := instance.client()
	if nil != err {
		return err
	}

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(seqNum, seqNum)

	return c.Store(seqSet, item, flags, nil)
}

func (instance *DriverImap) unmarkMessage(seqNum uint32, flags []string) error {
	c, err := instance.client()
	if nil != err {
		return err
	}

	item := imap.FormatFlagsOp(imap.RemoveFlags, true)
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(seqNum, seqNum)

	return c.Store(seqSet, item, flags, nil)
}

// ---------------------------------------------------------------------------------------------------------------------
//	S T A T I C
// ---------------------------------------------------------------------------------------------------------------------

func toFetchFlags(flags []interface{}) []imap.FetchItem {
	fetchFlags := make([]imap.FetchItem, 0)
	for _, flag := range flags {
		if v, b := flag.(imap.FetchItem); b {
			fetchFlags = append(fetchFlags, v)
		} else if s, b := flag.(string); b {
			fetchFlags = append(fetchFlags, imap.FetchItem(s))
		}
	}
	return fetchFlags
}

func fetchSeq(c *client.Client, seqNums []uint32, flags []imap.FetchItem) ([]*imap.Message, error) {
	if len(seqNums) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(seqNums...)
		return fetch(c, seqset, flags)
	}
	return []*imap.Message{}, nil
}

func fetchFrom(c *client.Client, from, to uint32, flags []imap.FetchItem) ([]*imap.Message, error) {
	if from > 0 && to > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddRange(from, to)
		return fetch(c, seqset, flags)
	}
	return []*imap.Message{}, nil
}

func fetch(c *client.Client, seqset *imap.SeqSet, flags []imap.FetchItem) ([]*imap.Message, error) {
	response := make([]*imap.Message, 0)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, flags, messages)
	}()
	// read messages
	for msg := range messages {
		response = append(response, msg)
	}
	// check error
	if err := <-done; err != nil {
		return nil, err
	}

	return response, nil
}

func isSeen(message *imap.Message) bool {
	for _, flag := range message.Flags {
		if flag == imap.SeenFlag {
			return true
		}
	}
	return false
}

func isRecent(message *imap.Message) bool {
	for _, flag := range message.Flags {
		if flag == imap.RecentFlag {
			return true
		}
	}
	return false
}

func isDeleted(message *imap.Message) bool {
	for _, flag := range message.Flags {
		if flag == imap.DeletedFlag {
			return true
		}
	}
	return false
}

func isDraft(message *imap.Message) bool {
	for _, flag := range message.Flags {
		if flag == imap.DraftFlag {
			return true
		}
	}
	return false
}

func isAnswered(message *imap.Message) bool {
	for _, flag := range message.Flags {
		if flag == imap.AnsweredFlag {
			return true
		}
	}
	return false
}

func isFlagged(message *imap.Message) bool {
	for _, flag := range message.Flags {
		if flag == imap.FlaggedFlag {
			return true
		}
	}
	return false
}

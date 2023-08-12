package lib_imap

import (
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_imap/commons"
)

type ImapClient struct {
	settings *commons.ImapClientSettings
	driver   *DriverImap
}

// ---------------------------------------------------------------------------------------------------------------------
//	ImapClient
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ImapClient) String() string {
	m := map[string]interface{}{
		"settings": instance.settings,
	}
	return qbc.Convert.ToString(m)
}

func (instance *ImapClient) Open() error {
	if nil == instance.driver {
		driver, err := NewDriverImap(instance.settings)
		if nil != err {
			return err
		}
		instance.driver = driver
	}
	return instance.driver.Open()
}

func (instance *ImapClient) Close() error {
	if nil != instance.driver {
		err := instance.driver.Close()
		instance.driver = nil
		return err
	}
	return nil
}

func (instance *ImapClient) CloseSilent() {
	if nil != instance {
		_ = instance.Close()
	}
}

func (instance *ImapClient) IsOpen() bool {
	if nil != instance && nil != instance.driver {
		return instance.driver.IsOpen()
	}
	return false
}

func (instance *ImapClient) ListMailboxes() ([]*commons.MailboxInfo, error) {
	if nil != instance.driver {
		return instance.driver.ListMailboxes()
	}
	return nil, nil
}

func (instance *ImapClient) GetMailboxFlags(mailboxName string) ([]string, error) {
	if nil != instance.driver {
		return instance.driver.GetMailboxFlags(mailboxName)
	}
	return nil, nil
}

func (instance *ImapClient) ReadMailbox(mailboxName string, onlyNew bool) ([]*commons.ImapMessage, error) {
	if nil != instance.driver {
		return instance.driver.ReadMailbox(mailboxName, onlyNew)
	}
	return nil, nil
}

func (instance *ImapClient) ReadMessage(seqNum interface{}) (*commons.ImapMessage, error) {
	if nil != instance.driver {
		return instance.driver.ReadMessage(seqNum)
	}
	return nil, nil
}
func (instance *ImapClient) MarkMessageAsSeen(seqNum interface{}) error {
	if nil != instance.driver {
		return instance.driver.MarkMessageAsSeen(seqNum)
	}
	return nil
}
func (instance *ImapClient) MarkMessageAsAnswered(seqNum interface{}) error {
	if nil != instance.driver {
		return instance.driver.MarkMessageAsAnswered(seqNum)
	}
	return nil
}
func (instance *ImapClient) MarkMessageAsDeleted(seqNum interface{}) error {
	if nil != instance.driver {
		return instance.driver.MarkMessageAsDeleted(seqNum)
	}
	return nil
}
func (instance *ImapClient) MarkMessageAsFlagged(seqNum interface{}) error {
	if nil != instance.driver {
		return instance.driver.MarkMessageAsFlagged(seqNum)
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
//	S T A T I C
// ---------------------------------------------------------------------------------------------------------------------

func NewImapClient(params ...interface{}) (*ImapClient, error) {
	instance := new(ImapClient)
	instance.settings = new(commons.ImapClientSettings)
	switch len(params) {
	case 1:
		// settings
		if s, b := params[0].(string); b {
			if qbc.Regex.IsValidJsonObject(s) {
				err := qbc.JSON.Read(s, &instance.settings)
				if nil != err {
					return nil, err
				}
				return instance, nil
			} else {
				// try loading data from file
				text, err := qbc.IO.ReadTextFromFile(s)
				if nil != err {
					return nil, err
				}
				return NewImapClient(text)
			}
		} else if settings, b := params[0].(*commons.ImapClientSettings); b {
			err := instance.settings.LoadFromText(qbc.JSON.Stringify(settings))
			if nil != err {
				return nil, err
			}
			return instance, nil
		} else if settings, b := params[0].(commons.ImapClientSettings); b {
			err := instance.settings.LoadFromText(qbc.JSON.Stringify(settings))
			if nil != err {
				return nil, err
			}
			return instance, nil
		} else if settings, b := params[0].(map[string]interface{}); b {
			err := instance.settings.LoadFromText(qbc.JSON.Stringify(settings))
			if nil != err {
				return nil, err
			}
			return instance, nil
		}
	default:
		return nil, commons.ErrorMismatchConfiguration
	}
	return instance, nil
}

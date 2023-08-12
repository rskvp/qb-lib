package commons

import (
	_ "embed"

	qbc "github.com/rskvp/qb-core"
)

// ---------------------------------------------------------------------------------------------------------------------
//	ClientSettingsAuth
// ---------------------------------------------------------------------------------------------------------------------

type MailboxerClientSettingsAuth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

// ---------------------------------------------------------------------------------------------------------------------
//	ImapClientSettings
// ---------------------------------------------------------------------------------------------------------------------

type ImapClientSettings struct {
	Type string                       `json:"type"`
	Host string                       `json:"host"`
	Port int                          `json:"port"`
	Tls  bool                         `json:"tls"`
	Auth *MailboxerClientSettingsAuth `json:"auth"`
}

func (instance *ImapClientSettings) String() string {
	return qbc.JSON.Stringify(instance)
}

func (instance *ImapClientSettings) LoadFromFile(filename string) error {
	text, err := qbc.IO.ReadTextFromFile(filename)
	if nil != err {
		return err
	}
	return instance.LoadFromText(text)
}

func (instance *ImapClientSettings) LoadFromText(text string) error {
	return qbc.JSON.Read(text, &instance)
}

func (instance *ImapClientSettings) LoadFromMap(m map[string]interface{}) error {
	return instance.LoadFromText(qbc.JSON.Stringify(m))
}

// ---------------------------------------------------------------------------------------------------------------------
//	S T A T I C
// ---------------------------------------------------------------------------------------------------------------------

func NewMailboxerClientSettings(stype, shost, iport, suser, spass interface{}) *ImapClientSettings {
	return &ImapClientSettings{
		Type: qbc.Convert.ToString(stype),
		Host: qbc.Convert.ToString(shost),
		Port: qbc.Convert.ToInt(iport),
		Auth: &MailboxerClientSettingsAuth{
			User: qbc.Convert.ToString(suser),
			Pass: qbc.Convert.ToString(spass),
		},
	}
}

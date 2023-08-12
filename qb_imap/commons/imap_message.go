package commons

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_html"
	imap_mime "github.com/rskvp/qb-lib/qb_imap/mime"
)

const (
	HeaderContentType     = "Content-Type"
	HeaderMessageId       = "Message-Id"
	HeaderParentMessageId = "Parent-Message-Id"
	HeaderDate            = "Date"
	HeaderSubject         = "Subject"
	HeaderSender          = "Sender"
	HeaderReplyTo         = "Reply-To"
	HeaderFrom            = "From"
	HeaderTo              = "To"
	HeaderCc              = "Cc"
	HeaderBcc             = "Bcc"
)

// ---------------------------------------------------------------------------------------------------------------------
//	ImapAttachment
// ---------------------------------------------------------------------------------------------------------------------

type ImapAttachment struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}

// ---------------------------------------------------------------------------------------------------------------------
//	MailboxerAddress
// ---------------------------------------------------------------------------------------------------------------------

type ImapAddress struct {
	PersonalName string `json:"personal-name"`  // The personal name.
	AtDomainList string `json:"at-domain-list"` // The SMTP at-domain-list (source route).
	MailboxName  string `json:"mailbox-name"`   // The mailbox name.
	HostName     string `json:"host-name"`      // The host name.
}

func (instance *ImapAddress) Parse(text string) {
	text = strings.ReplaceAll(text, "<", ",")
	text = strings.ReplaceAll(text, ">", "")
	fields := strings.Split(text, ",")
	var name, email string
	email = text
	if len(fields) == 2 {
		name = strings.TrimSpace(fields[0])
		email = strings.TrimSpace(fields[1])
	}
	instance.PersonalName = name
	if len(email) > 0 {
		tokens := strings.Split(email, "@")
		instance.MailboxName = qbc.Arrays.GetAt(tokens, 0, "").(string)
		instance.HostName = qbc.Arrays.GetAt(tokens, 1, "").(string)
	}
}

func (instance *ImapAddress) String() string {
	if len(instance.PersonalName) > 0 {
		return fmt.Sprintf("%s <%s@%s>", instance.PersonalName, instance.MailboxName, instance.HostName)
	}
	return fmt.Sprintf("%s@%s", instance.MailboxName, instance.HostName)
}

func (instance *ImapAddress) ToMailAddress() *mail.Address {
	return &mail.Address{
		Name:    instance.PersonalName,
		Address: fmt.Sprintf("%s@%s", instance.MailboxName, instance.HostName),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	MailboxerMessageHeader
// ---------------------------------------------------------------------------------------------------------------------

type ImapMessageHeader struct {
	MessageId       string         `json:"message-id"`
	ParentMessageId string         `json:"parent-message-id"`
	ReplyTo         []*ImapAddress `json:"reply-to"`
	Sender          []*ImapAddress `json:"sender"`
	From            []*ImapAddress `json:"from"`
	To              []*ImapAddress `json:"to"`
	Cc              []*ImapAddress `json:"cc"`
	Bcc             []*ImapAddress `json:"bcc"`
	Subject         string         `json:"subject"`
	Date            time.Time      `json:"date"`
}

func NewImapMessageHeader() *ImapMessageHeader {
	instance := new(ImapMessageHeader)
	instance.ReplyTo = make([]*ImapAddress, 0)
	instance.Sender = make([]*ImapAddress, 0)
	instance.From = make([]*ImapAddress, 0)
	instance.To = make([]*ImapAddress, 0)
	instance.Cc = make([]*ImapAddress, 0)
	instance.Bcc = make([]*ImapAddress, 0)
	return instance
}

func (instance *ImapMessageHeader) String() string {
	return qbc.JSON.Stringify(instance)
}

func (instance *ImapMessageHeader) SetFrom(value string) {
	instance.SetAddresses("from", value)
}

func (instance *ImapMessageHeader) SetTo(value string) {
	instance.SetAddresses("to", value)
}

func (instance *ImapMessageHeader) SetCc(value string) {
	instance.SetAddresses("cc", value)
}

func (instance *ImapMessageHeader) SetBcc(value string) {
	instance.SetAddresses("bcc", value)
}

func (instance *ImapMessageHeader) SetSender(value string) {
	instance.SetAddresses("sender", value)
}

func (instance *ImapMessageHeader) SetReplyTo(value string) {
	instance.SetAddresses("reply-to", value)
}

func (instance *ImapMessageHeader) SetAddresses(field, value string) {
	if len(value) > 0 {
		tokens := qbc.Strings.Split(value, ";,")
		for _, token := range tokens {
			address := new(ImapAddress)
			address.Parse(token)
			switch field {
			case "to":
				instance.To = append(instance.To, address)
			case "from":
				instance.From = append(instance.From, address)
			case "cc":
				instance.Cc = append(instance.Cc, address)
			case "bcc":
				instance.Bcc = append(instance.Bcc, address)
			case "sender":
				instance.Sender = append(instance.Sender, address)
			case "reply-to":
				instance.ReplyTo = append(instance.ReplyTo, address)
			}
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	MailboxerMessageBody
// ---------------------------------------------------------------------------------------------------------------------

type ImapMessageBody struct {
	MIMEType          string            `json:"mime-type"`          // The MIME type (e.g. "text", "image")
	MIMESubType       string            `json:"mime-sub-type"`      // The MIME subtype (e.g. "plain", "png")
	Params            map[string]string `json:"mime-params"`        // The MIME parameters.
	Id                string            `json:"id"`                 // The Content-Id header.
	Description       string            `json:"description"`        // The Content-Description header.
	Encoding          string            `json:"encoding"`           // The Content-Encoding header.
	Size              uint32            `json:"size"`               // The Content-Length header.
	Extended          bool              `json:"extended"`           // True if the body structure contains extension data.
	Disposition       string            `json:"disposition"`        // The Content-Disposition header field value.
	DispositionParams map[string]string `json:"disposition-params"` // The Content-Disposition header field parameters.
	Language          []string          `json:"language"`           // The Content-Language header field, if multipart.
	Location          []string          `json:"location"`           // The content URI, if multipart.
	MD5               string            `json:"md5"`                // The MD5 checksum.
	Multipart         bool              `json:"multipart"`
	PartsCount        int               `json:"parts-count"`
	Text              string            `json:"text"`
	HTML              string            `json:"html"`
	Attachments       []*ImapAttachment `json:"attachments"`
	Headers           map[string]string `json:"headers"`
}

func NewImapMessageBody() *ImapMessageBody {
	instance := new(ImapMessageBody)
	instance.Params = make(map[string]string)
	instance.DispositionParams = make(map[string]string)
	instance.Headers = make(map[string]string)
	instance.Attachments = make([]*ImapAttachment, 0)
	instance.Location = make([]string, 0)
	instance.Language = make([]string, 0)

	return instance
}

func (instance *ImapMessageBody) String() string {
	return qbc.JSON.Stringify(instance)
}

func (instance *ImapMessageBody) SetMIMEType(contentType string) {
	if len(contentType) > 0 {
		tokens := qbc.Strings.Split(contentType, ";")
		if len(tokens) > 0 {
			types := qbc.Strings.Split(tokens[0], "/")
			instance.MIMEType = qbc.Arrays.GetAt(types, 0, "").(string)
			instance.MIMESubType = qbc.Arrays.GetAt(types, 1, "").(string)
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	MailboxerMessage
// ---------------------------------------------------------------------------------------------------------------------

type ImapMessage struct {
	SeqNum       uint32             `json:"seq-num"`
	InternalDate time.Time          `json:"internal-date"` // The date the message was received by the server.
	Header       *ImapMessageHeader `json:"header"`
	Body         *ImapMessageBody   `json:"body"`
	BodyData     []byte
}

func NewImapMessage() *ImapMessage {
	target := &ImapMessage{}
	target.Header = NewImapMessageHeader()
	return target
}

func (instance *ImapMessage) Json() string {
	return qbc.JSON.Stringify(instance)
}

func (instance *ImapMessage) String() string {
	return string(instance.BodyData)
}

func (instance *ImapMessage) IsEmpty() bool {
	return nil == instance || nil == instance.Header || nil == instance.Body
}

func (instance *ImapMessage) IsRootMessage() bool {
	if nil != instance.Header {
		id := instance.Header.MessageId
		parentId := instance.Header.ParentMessageId
		if len(id) > 0 {
			return id == parentId || len(parentId) == 0
		}
	}
	return false
}

func (instance *ImapMessage) Subject() string {
	if nil != instance.Header && len(instance.Header.Subject) > 0 {
		return instance.Header.Subject
	}
	return "undefined"
}

func (instance *ImapMessage) PlainText() string {
	if nil != instance {
		if nil != instance.Body {
			if len(instance.Body.Text) > 0 {
				// first of all the text
				return instance.Body.Text
			} else if len(instance.Body.HTML) > 0 {
				// if no plain text try to parse HTML
				parser, err := qb_html.NewHtmlParser(instance.Body.HTML)
				if nil == err {
					return parser.TextAll()
				}
			}
		}
	}
	return ""
}

func (instance *ImapMessage) To() []*mail.Address {
	response := make([]*mail.Address, 0)
	if nil != instance && nil != instance.Header {
		for _, a := range instance.Header.To {
			response = append(response, a.ToMailAddress())
		}
	}
	return response
}

func (instance *ImapMessage) From() []*mail.Address {
	response := make([]*mail.Address, 0)
	tmp := make([]string, 0)
	if nil != instance && nil != instance.Header {
		for _, a := range instance.Header.From {
			ma := a.ToMailAddress()
			if qbc.Arrays.IndexOf(ma.String(), tmp) == -1 {
				tmp = append(tmp, ma.String())
				response = append(response, ma)
			}
		}
		for _, a := range instance.Header.Sender {
			ma := a.ToMailAddress()
			if qbc.Arrays.IndexOf(ma.String(), tmp) == -1 {
				tmp = append(tmp, ma.String())
				response = append(response, ma)
			}
		}
	}
	return response
}

func (instance *ImapMessage) ReplyTo() []*mail.Address {
	response := make([]*mail.Address, 0)
	if nil != instance && nil != instance.Header {
		for _, a := range instance.Header.ReplyTo {
			response = append(response, a.ToMailAddress())
		}
	}
	return response
}

func (instance *ImapMessage) Attachments() []*ImapAttachment {
	response := make([]*ImapAttachment, 0)
	if nil != instance && nil != instance.Body {
		return instance.Body.Attachments
	}
	return response
}

func (instance *ImapMessage) Date() time.Time {
	if nil != instance && nil != instance.Header {
		return instance.Header.Date
	}
	return time.Now()
}

func (instance *ImapMessage) SaveToFile(args ...interface{}) (err error) {
	filename := qbc.Strings.Slugify(instance.Subject()) + ".eml"
	if len(args) > 0 {
		filename = qbc.Convert.ToString(args[0])
	}
	ext := strings.ToLower(qbc.Paths.ExtensionName(filename))
	if ext == "json" {
		_, err = qbc.IO.WriteTextToFile(instance.Json(), filename)
	} else if len(instance.BodyData) > 0 {
		if len(ext) == 0 {
			filename = filename + ".eml"
		}
		_, err = qbc.IO.WriteBytesToFile(instance.BodyData, filename)
	}
	return
}

func (instance *ImapMessage) LoadFromFile(filename string) (err error) {
	data, err := qbc.IO.ReadBytesFromFile(filename)
	if nil != err {
		return err
	}
	envelope, err := parseBodyData(data, instance)
	if nil != err {
		return err
	}
	if nil != envelope {
		parseMime(envelope, instance)
	}
	return nil
}

func (instance *ImapMessage) ParseImap(m *imap.Message) *ImapMessage {
	if nil != m {
		parseImap(m, instance)
	}
	return instance
}

func (instance *ImapMessage) Parse(arg interface{}) (err error) {
	var data []byte
	switch v := arg.(type) {
	case io.Reader:
		data, err = ioutil.ReadAll(v)
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case *imap.Message:
		parseImap(v, instance)
	}
	if nil != err {
		return err
	}

	if len(data) > 0 {
		envelope, err := parseBodyData(data, instance)
		if nil != err {
			return err
		}
		if nil != envelope {
			parseMime(envelope, instance)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e    (MailboxerMessage )
// ---------------------------------------------------------------------------------------------------------------------

func parseMime(envelope *imap_mime.Envelope, target *ImapMessage) {

}

func parseImap(source *imap.Message, target *ImapMessage) {
	if nil != source && nil != target {
		target.SeqNum = source.SeqNum
		target.InternalDate = source.InternalDate
		if nil == target.Header {
			target.Header = NewImapMessageHeader()
		}
		if nil == target.Body {
			target.Body = NewImapMessageBody()
		}

		// ENVELOPE
		if len(target.Header.MessageId) == 0 {
			envelope := source.Envelope
			if nil != envelope {
				target.Header.MessageId = envelope.MessageId
				target.Header.Subject = envelope.Subject
				target.Header.Date = envelope.Date
				target.Header.ParentMessageId = envelope.InReplyTo // parent MessageId

				// REPLY-TO
				target.Header.ReplyTo = make([]*ImapAddress, 0)
				for _, a := range envelope.ReplyTo {
					address := &ImapAddress{
						PersonalName: a.PersonalName,
						AtDomainList: a.AtDomainList,
						MailboxName:  a.MailboxName,
						HostName:     a.HostName,
					}
					target.Header.ReplyTo = append(target.Header.ReplyTo, address)
				}
				// SENDER
				target.Header.Sender = make([]*ImapAddress, 0)
				for _, a := range envelope.Sender {
					address := &ImapAddress{
						PersonalName: a.PersonalName,
						AtDomainList: a.AtDomainList,
						MailboxName:  a.MailboxName,
						HostName:     a.HostName,
					}
					target.Header.Sender = append(target.Header.Sender, address)
				}
				// FROM
				target.Header.From = make([]*ImapAddress, 0)
				for _, a := range envelope.From {
					address := &ImapAddress{
						PersonalName: a.PersonalName,
						AtDomainList: a.AtDomainList,
						MailboxName:  a.MailboxName,
						HostName:     a.HostName,
					}
					target.Header.From = append(target.Header.From, address)
				}
				// TO
				target.Header.To = make([]*ImapAddress, 0)
				for _, a := range envelope.To {
					address := &ImapAddress{
						PersonalName: a.PersonalName,
						AtDomainList: a.AtDomainList,
						MailboxName:  a.MailboxName,
						HostName:     a.HostName,
					}
					target.Header.To = append(target.Header.To, address)
				}
				// CC
				target.Header.Cc = make([]*ImapAddress, 0)
				for _, a := range envelope.Cc {
					address := &ImapAddress{
						PersonalName: a.PersonalName,
						AtDomainList: a.AtDomainList,
						MailboxName:  a.MailboxName,
						HostName:     a.HostName,
					}
					target.Header.Cc = append(target.Header.Cc, address)
				}
				// BCC
				target.Header.Bcc = make([]*ImapAddress, 0)
				for _, a := range envelope.Bcc {
					address := &ImapAddress{
						PersonalName: a.PersonalName,
						AtDomainList: a.AtDomainList,
						MailboxName:  a.MailboxName,
						HostName:     a.HostName,
					}
					target.Header.Bcc = append(target.Header.Bcc, address)
				}
			}
		}

		// BODY STRUCTURE
		if nil != source.BodyStructure {
			// The children parts, if multipart.
			parseImapBodyStructure(source.BodyStructure, source.Body, target)
		}
	}
}

func parseImapBodyStructure(source *imap.BodyStructure, body map[*imap.BodySectionName]imap.Literal, targetMessage *ImapMessage) {
	if nil == targetMessage.Body {
		targetMessage.Body = new(ImapMessageBody)
	}
	targetBody := targetMessage.Body
	// copy body structure into target output message
	imapMergeBodyStructure(source, targetBody)

	if len(body) > 0 {
		for _, r := range body {
			_ = parseBody(r, targetMessage)
		}
	}

}

func imapMergeBodyStructure(source *imap.BodyStructure, target *ImapMessageBody) {
	target.Id = source.Id
	target.Description = source.Description
	target.MIMEType = source.MIMEType
	target.MIMESubType = source.MIMESubType
	target.Params = source.Params
	target.Encoding = source.Encoding
	target.Size = source.Size
	target.Extended = source.Extended
	target.MD5 = source.MD5
	target.Disposition = source.Disposition
	target.DispositionParams = source.DispositionParams
	target.Language = source.Language
	target.Location = source.Location
	target.PartsCount = len(source.Parts)
	target.Multipart = len(source.Parts) > 0
}

func parseBody(r io.Reader, targetMessage *ImapMessage) error {
	data, err := ioutil.ReadAll(r)
	if nil != err {
		return err
	}
	_, err = parseBodyData(data, targetMessage)
	return err
}

func parseBodyData(data []byte, targetMessage *ImapMessage) (*imap_mime.Envelope, error) {
	targetMessage.BodyData = data
	if nil == targetMessage.Body {
		targetMessage.Body = new(ImapMessageBody)
	}

	// mime parser
	env, envErr := imap_mime.ReadEnvelope(bytes.NewReader(data))
	if nil == envErr && nil != env {
		// content
		targetMessage.Body.Text = env.Text
		targetMessage.Body.HTML = env.HTML

		// headers
		targetMessage.Body.Headers = make(map[string]string)
		keys := env.GetHeaderKeys()
		for _, key := range keys {
			value := env.GetHeader(key)
			targetMessage.Body.Headers[key] = value
		}

		// attachments
		if len(env.Attachments) > 0 {
			targetMessage.Body.Attachments = make([]*ImapAttachment, 0)
			for _, attachment := range env.Attachments {
				targetMessage.Body.Attachments = append(targetMessage.Body.Attachments, &ImapAttachment{
					Filename: attachment.FileName,
					Data:     attachment.Content,
				})
			}
		}

		// is header ready?
		if nil == targetMessage.Header {
			parseHeader(env, targetMessage)
		}

	}
	return env, nil
}

func parseHeader(env *imap_mime.Envelope, target *ImapMessage) {
	if nil == target.Header {
		target.Header = NewImapMessageHeader()
	}
	if nil == target.Body {
		target.Body = NewImapMessageBody()
	}

	// HEADER
	if len(target.Header.MessageId) == 0 {
		target.Header.MessageId = env.GetHeader(HeaderMessageId)
	}
	if len(target.Header.ParentMessageId) == 0 {
		target.Header.ParentMessageId = env.GetHeader(HeaderParentMessageId)
	}
	if qbc.Dates.IsZero(target.Header.Date) {
		dt, err := qbc.Dates.ParseAny(env.GetHeader(HeaderDate))
		if nil == err {
			target.Header.Date = dt
		}
	}
	if len(target.Header.Subject) == 0 {
		target.Header.Subject = env.GetHeader(HeaderSubject)
	}
	if len(target.Header.Sender) == 0 {
		target.Header.SetSender(env.GetHeader(HeaderSender))
	}
	if len(target.Header.ReplyTo) == 0 {
		target.Header.SetReplyTo(env.GetHeader(HeaderReplyTo))
	}
	if len(target.Header.From) == 0 {
		target.Header.SetFrom(env.GetHeader(HeaderFrom))
	}
	if len(target.Header.To) == 0 {
		target.Header.SetTo(env.GetHeader(HeaderTo))
	}
	if len(target.Header.Cc) == 0 {
		target.Header.SetCc(env.GetHeader(HeaderCc))
	}
	if len(target.Header.Bcc) == 0 {
		target.Header.SetBcc(env.GetHeader(HeaderBcc))
	}

	// MESSAGE & BODY
	if len(target.Body.MIMEType) == 0 {
		target.Body.SetMIMEType(env.GetHeader(HeaderContentType))
	}
	if qbc.Dates.IsZero(target.InternalDate) {
		dt, err := qbc.Dates.ParseAny(env.GetHeader(HeaderDate))
		if nil == err {
			target.InternalDate = dt
		}
	}

}

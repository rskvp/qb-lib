package nodemailer

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_email"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "nodemailer"

type Path struct {
	runtime *goja.Runtime
	util    *goja.Object
	root    string
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// nodemailer.createTransport(options[, defaults])
func (instance *Path) createTransport(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) > 0 {
		options := call.Argument(0).Export()
		if o, b := options.(map[string]interface{}); b {
			// creates transport object
			transport := &map[string]interface{}{
				"sendMail": instance.sendMail(o),
			}
			return instance.runtime.ToValue(transport)
		} else {
			// invalid options
			panic(instance.runtime.NewTypeError("Invalid Options."))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Path) sendMail(options map[string]interface{}) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if nil != options && len(call.Arguments) == 2 {
			host := qbc.Reflect.GetString(options, "host")
			port := qbc.Reflect.GetInt(options, "port")
			secure := qbc.Reflect.GetBool(options, "secure")
			user := qbc.Reflect.GetString(qbc.Reflect.Get(options, "auth"), "user")
			pass := qbc.Reflect.GetString(qbc.Reflect.Get(options, "auth"), "pass")
			if len(host) > 0 && port > 0 && len(user) > 0 && len(pass) > 0 {
				message := call.Argument(0).Export()
				callback, _ := goja.AssertFunction(call.Argument(len(call.Arguments) - 1))
				if nil != message && nil != callback {
					from := qbc.Reflect.GetString(message, "from")
					to := qbc.Reflect.GetString(message, "to")
					subject := qbc.Reflect.GetString(message, "subject")
					text := qbc.Reflect.GetString(message, "text")
					html := qbc.Reflect.GetString(message, "html")
					attachments := qbc.Reflect.GetArray(message, "attachments")
					if len(from) > 0 && len(to) > 0 && len(subject) > 0 && (len(text) > 0 || len(html) > 0) {
						err := sendEmail(host, port, secure, user, pass, from, to, subject, text, html, attachments)
						if nil != err {
							_, _ = callback(call.This, instance.runtime.ToValue(err.Error()), goja.Undefined())
						} else {
							info := &map[string]interface{}{}
							_, _ = callback(call.This, goja.Undefined(), instance.runtime.ToValue(info))
						}
					} else {
						panic(instance.runtime.NewTypeError("Missing Parameters in Message Object."))
					}
				} else {
					panic(instance.runtime.NewTypeError("Missing Message Object."))
				}
			}
		}
		return goja.Undefined()
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

// https://nodemailer.com/message/attachments/
func sendEmail(host string, port int, secure bool, user string, pass string,
	from string, to string, subject string, text string, html string, attachments []interface{}) (err error) {
	servername := fmt.Sprintf("%v:%v", host, port)
	auth := smtp.PlainAuth("", user, pass, host)
	toList := qbc.Strings.Split(to, ";,")
	var m *qb_email.Message
	if len(html) > 0 {
		m = qbc.Email.NewHTMLMessage(subject, html)
	} else {
		m = qbc.Email.NewMessage(subject, text)
	}
	addr, err := mail.ParseAddress(from)
	if nil != err {
		m.From = &mail.Address{Address: from}
	} else {
		m.From = addr
	}
	m.To = toList
	for _, attachment := range attachments {
		if nil != attachment {
			if v, b := attachment.(string); b {
				err = addAttachmentString(m, v)
			} else if v, b := attachment.(map[string]interface{}); b {
				err = addAttachmentObject(m, v)
			}
		}
	}
	if nil != err {
		return
	}
	if secure {
		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
		err = qbc.Email.SendSecure(servername, auth, tlsconfig, m)
		return
	} else {
		err = qbc.Email.Send(servername, auth, m)
		return
	}
}

func addAttachmentString(m *qb_email.Message, attachment string) error {
	filename := qbc.Paths.FileName(attachment, true)
	return addAttachment(m, filename, attachment)
}

func addAttachmentObject(m *qb_email.Message, attachment map[string]interface{}) error {
	filename := qbc.Reflect.GetString(attachment, "filename")
	path := qbc.Reflect.GetString(attachment, "path")
	return addAttachment(m, filename, path)
}

func addAttachment(m *qb_email.Message, filename, path string) error {
	if len(filename) > 0 && len(path) > 0 {
		data, err := download(path)
		if nil != err {
			return err
		}
		return m.AddAttachmentBinary(filename, data, false)
	}
	return nil // nothing to attach
}

func download(url string) ([]byte, error) {
	if len(url) > 0 {
		if strings.Index(url, "http") > -1 {
			// HTTP
			tr := &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    15 * time.Second,
				DisableCompression: true,
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Get(url)
			if nil == err {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if nil == err {
					return body, nil
				} else {
					return []byte{}, err
				}
			} else {
				return []byte{}, err
			}
		} else {
			// FILE SYSTEM
			path := url
			return qbc.IO.ReadBytesFromFile(path)
		}
	}
	return []byte{}, qbc.Errors.Prefix(errors.New(url), "Invalid url or path: ")
}

func load(runtime *goja.Runtime, module *goja.Object, args ...interface{}) {
	instance := &Path{
		runtime: runtime,
	}

	if len(args) > 0 {
		root := qbc.Reflect.ValueOf(args[0]).String()
		if len(root) > 0 {
			instance.root = root
		}
	}

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("createTransport", instance.createTransport)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

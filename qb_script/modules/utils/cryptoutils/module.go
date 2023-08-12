package cryptoutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
	"github.com/rskvp/qb-lib/qb_script/modules/defaults/require"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "crypto-utils"

type CryptoUtils struct {
	runtime *goja.Runtime
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// crypto.encodeBase64(data)
func (instance *CryptoUtils) encodeBase64(call goja.FunctionCall) goja.Value {
	var text string
	idata := call.Argument(0).Export()
	if data, b := idata.([]uint8); b {
		text = qbc.Coding.EncodeBase64(data)
	} else if data, b := idata.([]byte); b {
		text = qbc.Coding.EncodeBase64(data)
	} else if data, b := idata.(string); b {
		text = qbc.Coding.EncodeBase64([]byte(data))
	} else {
		panic(instance.runtime.NewTypeError("invalid_data_type"))
	}
	return instance.runtime.ToValue(text)
}

// crypto.decodeBase64(data)
func (instance *CryptoUtils) decodeBase64(call goja.FunctionCall) goja.Value {
	text := commons.GetString(call, 0)
	data, err := qbc.Coding.DecodeBase64(text)
	if nil != err {
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return instance.runtime.ToValue(data)
}

// crypto.decodeBase64ToText(data)
func (instance *CryptoUtils) decodeBase64ToText(call goja.FunctionCall) goja.Value {
	text := call.Argument(0).String()
	data, err := qbc.Coding.DecodeBase64(text)
	if nil != err {
		panic(instance.runtime.NewTypeError(err.Error()))
	}
	return instance.runtime.ToValue(string(data))
}

//----------------------------------------------------------------------------------------------------------------------
//	M D 5
//----------------------------------------------------------------------------------------------------------------------

func (instance *CryptoUtils) md5(call goja.FunctionCall) goja.Value {
	text := commons.GetString(call, 0)
	data := qbc.Coding.MD5(text)
	return instance.runtime.ToValue(data)
}

//----------------------------------------------------------------------------------------------------------------------
//	H M A C
//----------------------------------------------------------------------------------------------------------------------

func (instance *CryptoUtils) encodeSha256(call goja.FunctionCall) goja.Value {
	if nil != instance {
		var secret, message string
		switch len(call.Arguments) {
		case 2:
			secret = commons.GetString(call, 0)
			message = commons.GetString(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		if len(secret) > 0 && len(message) > 0 {
			// encode
			bytes := instance.encode(sha256.New, []byte(secret), []byte(message))
			return instance.runtime.ToValue(bytes)
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *CryptoUtils) encodeSha512(call goja.FunctionCall) goja.Value {
	if nil != instance {
		var secret, message string
		switch len(call.Arguments) {
		case 2:
			secret = commons.GetString(call, 0)
			message = commons.GetString(call, 1)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		if len(secret) > 0 && len(message) > 0 {
			// encode
			bytes := instance.encode(sha512.New, []byte(secret), []byte(message))
			return instance.runtime.ToValue(bytes)
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *CryptoUtils) encode(h func() hash.Hash, secret []byte, message []byte) []byte {
	hasher := hmac.New(h, secret)
	hasher.Write(message)
	return hasher.Sum(nil)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, _ ...interface{}) {
	instance := &CryptoUtils{
		runtime: runtime,
	}

	o := module.Get("exports").(*goja.Object)

	// base64
	_ = o.Set("encodeBase64", instance.encodeBase64)
	_ = o.Set("decodeBase64", instance.decodeBase64)
	_ = o.Set("decodeBase64ToText", instance.decodeBase64ToText)
	// hmac
	_ = o.Set("encodeSha256", instance.encodeSha256)
	_ = o.Set("encodeSha512", instance.encodeSha512)
	// md5
	_ = o.Set("md5", instance.md5)

}

func Enable(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}

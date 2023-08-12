package server

import (
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_events"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

// Source @url:https://github.com/gorilla/websocket/blob/master/conn.go#L61
// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

const (
	OnDisconnectEvent = "on_disconnect"
	OnMessageEvent    = "on_message"
)

var (
	poolMux sync.Mutex
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpWebsocketConn struct {
	UUID string

	//-- pr i v a t e --//
	pool   map[string]*HttpWebsocketConn
	conn   *websocket.Conn
	events *qb_events.Emitter
	queue  []*Message
	alive  bool
}

type HttpWebsocketEventPayload struct {
	Websocket *HttpWebsocketConn // sender
	Message   *Message
	Error     error
}

type Message struct {
	Type int
	Data []byte
}

type httpServerConfigRouteWebsocket struct {
	Path    string
	Handler func(c *HttpWebsocketConn)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func NewHttpWebsocketConn(conn *websocket.Conn, pool map[string]*HttpWebsocketConn) *HttpWebsocketConn {
	instance := new(HttpWebsocketConn)
	instance.UUID = qbc.Rnd.Uuid()
	instance.conn = conn
	instance.pool = pool
	instance.events = qbc.Events.NewEmitter()

	instance.register()

	return instance
}

func (instance *HttpWebsocketConn) Conn() *websocket.Conn {
	if nil != instance && nil != instance.conn {
		return instance.conn
	}
	return nil
}

func (instance *HttpWebsocketConn) Locals(key string) interface{} {
	if nil != instance && nil != instance.conn {
		return instance.conn.Locals(key)
	}
	return nil
}

func (instance *HttpWebsocketConn) Join() {
	instance.alive = true
	instance.run()
}

func (instance *HttpWebsocketConn) IsAlive() bool {
	if nil != instance {
		return instance.alive
	}

	return false
}

func (instance *HttpWebsocketConn) Shutdown(err error) error {
	if nil != instance && nil != instance.conn {
		instance.unregister()
		// close event
		instance.events.Emit(OnDisconnectEvent, err)
		instance.events.Clear()
		// stop tickers
		instance.alive = false

		return instance.conn.Close()
	}
	return nil
}

func (instance *HttpWebsocketConn) Send(messageType int, data []byte) {
	if nil != instance {
		instance.write(messageType, data)
	}
}

func (instance *HttpWebsocketConn) SendText(text string) {
	if nil != instance {
		instance.write(TextMessage, []byte(text))
	}
}

func (instance *HttpWebsocketConn) SendData(data []byte) {
	if nil != instance {
		instance.write(TextMessage, data)
	}
}

func (instance *HttpWebsocketConn) SendTo(uuid string, messageType int, data []byte) {
	if nil != instance {
		ws := instance.clientGet(uuid)
		if nil != ws {
			ws.write(messageType, data)
		}
	}
}

func (instance *HttpWebsocketConn) SendTextTo(uuid string, text string) {
	if nil != instance {
		instance.SendTo(uuid, TextMessage, []byte(text))
	}
}

func (instance *HttpWebsocketConn) SendDataTo(uuid string, data []byte) {
	if nil != instance {
		instance.SendTo(uuid, TextMessage, data)
	}
}

func (instance *HttpWebsocketConn) ClientsCount() int {
	count := 0
	if nil != instance && nil != instance.pool {
		poolMux.Lock()
		count = len(instance.pool)
		poolMux.Unlock()
	}
	return count
}

func (instance *HttpWebsocketConn) ClientsUUIDs() []string {
	response := make([]string, 0)
	if nil != instance && nil != instance.pool {
		poolMux.Lock()
		for _, w := range instance.pool {
			response = append(response, w.UUID)
		}
		poolMux.Unlock()
	}
	return response
}

func (instance *HttpWebsocketConn) ClientByUUID(uuid string) *HttpWebsocketConn {
	if nil != instance && nil != instance.pool {
		return instance.clientGet(uuid)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	e v e n t s
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpWebsocketConn) OnDisconnect(callback func(payload *HttpWebsocketEventPayload)) {
	if nil != instance && nil != instance.events && nil != callback {
		instance.events.On(OnDisconnectEvent, func(event *qb_events.Event) {
			callback(&HttpWebsocketEventPayload{
				Websocket: instance,
				Error:     event.ArgumentAsError(0),
			})
		})
	}
}

func (instance *HttpWebsocketConn) OnMessage(callback func(payload *HttpWebsocketEventPayload)) {
	if nil != instance && nil != instance.events && nil != callback {
		instance.events.On(OnMessageEvent, func(event *qb_events.Event) {
			// fmt.Println("OnMessage", event)
			callback(&HttpWebsocketEventPayload{
				Websocket: instance,
				Message: &Message{
					Type: event.ArgumentAsInt(0),
					Data: event.ArgumentAsBytes(1),
				},
			})
		})
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpWebsocketConn) register() {
	// add to pool
	poolMux.Lock()
	instance.pool[instance.UUID] = instance
	poolMux.Unlock()
}

func (instance *HttpWebsocketConn) unregister() {
	// remove pool
	poolMux.Lock()
	delete(instance.pool, instance.UUID)
	poolMux.Unlock()
}

func (instance *HttpWebsocketConn) clientGet(uuid string) *HttpWebsocketConn {
	var ws *HttpWebsocketConn
	if nil != instance && nil != instance.pool {
		poolMux.Lock()
		if val, ok := instance.pool[uuid]; ok {
			ws = val
		}
		poolMux.Unlock()
	}
	return ws
}

func (instance *HttpWebsocketConn) run() {

	// start sending a test message to client
	go instance.pong()
	// start reading incoming messages
	go instance.read()

	// read message queue until the end and write to client stream
	for range time.Tick(1 * time.Millisecond) {
		if nil != instance {
			if instance.alive {
				if len(instance.queue) > 0 {
					// start loop on message buffer
					for _, message := range instance.queue {
						// write to client
						err := instance.conn.WriteMessage(message.Type, message.Data)
						if err != nil {
							_ = instance.Shutdown(err)
						}
					}
					// empty message buffer
					instance.queue = nil
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	// exit, no more alive
}

func (instance *HttpWebsocketConn) pong() {
	for range time.Tick(5 * time.Second) {
		if nil != instance {
			if instance.alive {
				// test client is alive
				instance.write(PongMessage, []byte{})
			} else {
				break
			}
		} else {
			break
		}
	}
}

func (instance *HttpWebsocketConn) write(mType int, data []byte) {
	instance.queue = append(instance.queue, &Message{
		Type: mType,
		Data: data,
	})
}

func (instance *HttpWebsocketConn) read() {
	for range time.Tick(10 * time.Millisecond) {
		if nil != instance {
			if instance.alive {
				t, m, e := instance.conn.ReadMessage()
				if e != nil {
					_ = instance.Shutdown(e)
					break
				}
				switch t {
				case PingMessage:
					// ping
				case CloseMessage:
					// close
					_ = instance.Shutdown(nil)
					break
				default:
					instance.events.Emit(OnMessageEvent, t, m)
				}
			} else {
				break
			}
		} else {
			break
		}
	}
}

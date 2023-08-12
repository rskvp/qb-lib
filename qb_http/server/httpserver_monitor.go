package server

import (
	"sync"
	"time"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_ticker"
)

const delay = 1
const EventOnFileChanged = "on_file_changed"

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e
// ---------------------------------------------------------------------------------------------------------------------

type ServerMonitor struct {
	files         []string // files to monitor for change
	filesChecksum []string

	events  *qb_events.Emitter
	fileMux sync.Mutex
	ticker  *qb_ticker.Ticker
}

// ---------------------------------------------------------------------------------------------------------------------
//		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewMonitor(files []string) *ServerMonitor {
	instance := new(ServerMonitor)
	instance.files = files
	instance.events = qbc.Events.NewEmitter()

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ServerMonitor) Start() {
	if nil != instance && nil == instance.ticker {
		instance.init()
	}
}

func (instance *ServerMonitor) Stop() {
	if nil != instance && nil != instance.ticker {
		instance.ticker.Stop()
		instance.ticker = nil
	}
}

func (instance *ServerMonitor) OnFileChanged(callback func(event *qb_events.Event)) {
	if nil != instance && nil != instance.events {
		instance.events.On(EventOnFileChanged, callback)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ServerMonitor) init() {
	instance.filesChecksum = make([]string, len(instance.files))
	for idx, file := range instance.files {
		if b, _ := qbc.Paths.Exists(file); b {
			text, err := qbc.IO.ReadTextFromFile(file)
			if nil != err {
				instance.filesChecksum[idx] = ""
			} else {
				instance.filesChecksum[idx] = qbc.Coding.MD5(text)
			}
		} else {
			instance.filesChecksum[idx] = ""
		}
	}
	instance.ticker = qb_ticker.NewTicker(delay*time.Second, func(t *qb_ticker.Ticker) {
		instance.check()
	})
	instance.ticker.Start()
}

func (instance *ServerMonitor) check() {
	if nil != instance {
		instance.fileMux.Lock()
		defer instance.fileMux.Unlock()

		for idx, file := range instance.files {
			if b, _ := qbc.Paths.Exists(file); b {
				text, err := qbc.IO.ReadTextFromFile(file)
				if nil == err {
					key := qbc.Coding.MD5(text)
					if key != instance.filesChecksum[idx] {
						instance.filesChecksum[idx] = key
						instance.events.EmitAsync(EventOnFileChanged, file)
						break
					}
				}
			}
		}
	}
}

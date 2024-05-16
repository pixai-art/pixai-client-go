package pixai_client

import (
	"sync"
)

type handler struct {
	handle func(*AllEvents)
	once   bool
}

type eventSource struct {
	sid            string
	ch             <-chan *AllEvents
	handlers       map[int]*handler
	handlersMutex  sync.Mutex
	handlersNextId int
	started        bool
}

func (e *eventSource) RemoveEventListener(id int) {
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	delete(e.handlers, id)
}

func (e *eventSource) AddEventListener(handle func(*AllEvents)) int {
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	e.handlersNextId++
	e.handlers[e.handlersNextId] = &handler{
		handle: handle,
		once:   false,
	}
	return e.handlersNextId
}

func (e *eventSource) Start() {
	if e.started {
		return
	}
	e.started = true

	go func() {
		for ae := range e.ch {
			if !e.started {
				return
			}
			e.handlersMutex.Lock()
			toDelete := []int{}
			for id, handler := range e.handlers {
				handler.handle(ae)
				if handler.once {
					toDelete = append(toDelete, id)
				}
			}
			for _, id := range toDelete {
				delete(e.handlers, id)
			}
			e.handlersMutex.Unlock()
		}
		e.started = false
	}()
}

func (e *eventSource) Close() {
	e.started = false
}

package socketio
import "github.com/equalll/mydebug"

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
)

type eventHandler struct {
	fn   reflect.Value
	args []reflect.Type
}

type EventEmitter struct {
	mutex  sync.Mutex
	events map[string][]*eventHandler
}

func NewEventEmitter() *EventEmitter {mydebug.INFO()
	return &EventEmitter{events: make(map[string][]*eventHandler)}
}

// global cache
var eventHandlerCache = &struct {
	sync.RWMutex
	cache map[uintptr]*eventHandler
}{cache: make(map[uintptr]*eventHandler)}

func genEventHandler(fn interface{}) (handler *eventHandler, err error) {mydebug.INFO()
	// if a handler have been generated before, use it first
	fnValue := reflect.ValueOf(fn)
	eventHandlerCache.RLock()
	if handler, ok := eventHandlerCache.cache[fnValue.Pointer()]; ok {
		eventHandlerCache.RUnlock()
		return handler, nil
	}
	eventHandlerCache.RUnlock()
	handler = new(eventHandler)
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		err = fmt.Errorf("%v is not a function", fn)
		return
	}
	handler.fn = fnValue
	fnType := fnValue.Type()
	nArgs := fnValue.Type().NumIn()
	handler.args = make([]reflect.Type, nArgs)
	if nArgs == 0 {
		err = errors.New("no arg exists")
		return
	}
	if t := fnType.In(0); t.Kind() != reflect.Ptr || t.Elem().Name() != "NameSpace" {
		err = errors.New("first argument should be of type *NameSpace")
		return
	} else {
		handler.args[0] = t
	}
	for i := 1; i < nArgs; i++ {
		handler.args[i] = fnType.In(i)
	}
	eventHandlerCache.Lock()
	eventHandlerCache.cache[fnValue.Pointer()] = handler
	eventHandlerCache.Unlock()
	return
}

func (ee *EventEmitter) On(name string, fn interface{}) error {mydebug.INFO()
	handler, err := genEventHandler(fn)
	if err != nil {
		return err
	}
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	ee.events[name] = append(ee.events[name], handler)
	return nil
}

func (ee *EventEmitter) RemoveListener(name string, fn interface{}) {mydebug.INFO()
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	for i, handler := range ee.events[name] {
		if handler.fn.Pointer() == reflect.ValueOf(fn).Pointer() {
			ee.events[name] = append(ee.events[name][0:i], ee.events[name][i+1:]...)
			break
		}
	}
	if len(ee.events[name]) == 0 {
		delete(ee.events, name)
	}
}

func (ee *EventEmitter) RemoveAllListeners(name string) {mydebug.INFO()
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	// assign nil?
	delete(ee.events, name)
}

func (ee *EventEmitter) fetchHandlers(name string) (handlers []*eventHandler) {mydebug.INFO()
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	handlers = ee.events[name]
	return
}

func (ee *EventEmitter) emit(name string, ns *NameSpace, callback func([]interface{}), args ...interface{}) {mydebug.INFO()
	handlers := ee.fetchHandlers(name)
	callArgs := make([]reflect.Value, len(args)+1)
	callArgs[0] = reflect.ValueOf(ns)
	for i, arg := range args {
		callArgs[i+1] = reflect.ValueOf(arg)
	}
	for _, handler := range handlers {
		go safeCall(handler.fn, callArgs, callback)
	}
}

func genAckCallback(ns *NameSpace, eventPacketCommon packetCommon) reflect.Value {mydebug.INFO()
	return reflect.ValueOf(func(args ...interface{}) {
		p := new(ackPacket)
		p.ackId = eventPacketCommon.id
		p.packetCommon = packetCommon{}

		var err error
		p.args, err = json.Marshal(args)
		if err != nil {
			fmt.Println(err)
		}
		err = ns.sendPacket(p)
		if err != nil {
			fmt.Println(err)
		}
	})
}

func (ee *EventEmitter) emitRaw(name string, ns *NameSpace, callback func([]interface{}), data []byte, eventPacketCommon packetCommon) error {mydebug.INFO()
	handlers := ee.fetchHandlers(name)
	var callArgs []reflect.Value
	if len(handlers) != 0 {
		handler := handlers[0]
		args := make([]interface{}, len(handler.args)-1)
		for i, arg := range handler.args[1:] {
			args[i] = reflect.New(arg).Interface()
		}
		if len(data) != 0 {
			err := json.Unmarshal(data, &args)
			if err != nil {
				return err
			}
		}
		callArgs = []reflect.Value{reflect.ValueOf(ns)}

		if args != nil && len(args) > 0 && args[0] != nil {
			for _, arg := range args {
				val := reflect.ValueOf(arg)
				if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				callArgs = append(callArgs, val)
			}
		}
	}

	if eventPacketCommon.ack {
		foundCallback := false
		for i, arg := range callArgs {
			if arg.Kind() == reflect.Func {
				callArgs[i] = genAckCallback(ns, eventPacketCommon)
				foundCallback = true
			}
		}
		if !foundCallback {
			callArgs = append(callArgs, genAckCallback(ns, eventPacketCommon))
		}
	}

	for _, handler := range handlers {
		go safeCall(handler.fn, callArgs, callback)
	}
	return nil
}

func safeCall(fn reflect.Value, args []reflect.Value, callback func([]interface{})) {mydebug.INFO()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	ret := fn.Call(args)
	if len(ret) > 0 {
		retArgs := make([]interface{}, len(ret))
		for i, arg := range ret {
			retArgs[i] = arg.Interface()
		}
		if callback != nil {
			callback(retArgs)
		}
	}
}

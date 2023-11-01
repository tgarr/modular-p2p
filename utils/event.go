package utils

import (
    "sync"
)

// event states
const (
    EVENT_STATE_PENDING         = 0
    EVENT_STATE_ABORTED         = 1
    EVENT_STATE_HANDLED         = 2
    EVENT_STATE_NOTHANDLED      = 3
)

// ==== interfaces ====

// event abstraction for the simulator
type IEvent interface {
    GetType() uint16
    GetData() interface{}
    GetDestination() IEventDestination
    GetState() uint16
    GetTime() float64

    Abort()
    SetState(state uint16) IEvent
    SetTime(time float64) IEvent
}

// An event destination is any component that can handle an event
type IEventDestination interface {
    HandleEvent(event IEvent) bool
}

// ==== concrete structures ====

/*
    Default implementation of an event.

    Implements: IEvent
*/
type Event struct {
    eventType uint16
    data interface{}
    destination IEventDestination
    state uint16
    time float64

    // sync
    stateLock sync.RWMutex
}

// ==== factories ====

func NewEvent(tp uint16,data interface{},destination IEventDestination) IEvent {
    return &Event{
        eventType:      tp,
        data:           data,
        destination:    destination,
        state:          EVENT_STATE_PENDING,
        stateLock:      sync.RWMutex{},
        time:           0, // optional
    }
}

// ==== getters ====

func (ev *Event) GetType() uint16 {
    return ev.eventType
}

func (ev *Event) GetData() interface{} {
    return ev.data
}

func (ev *Event) GetDestination() IEventDestination {
    return ev.destination
}

func (ev *Event) GetState() uint16 {
    ev.stateLock.RLock()
    defer ev.stateLock.RUnlock()
    
    return ev.state
}

func (ev *Event) GetTime() float64 {
    return ev.time
}

// ==== setters ====

func (ev *Event) SetState(state uint16) IEvent {
    ev.stateLock.Lock()
    defer ev.stateLock.Unlock()

    ev.state = state
    return ev
}

func (ev *Event) Abort() {
    ev.SetState(EVENT_STATE_ABORTED)
}

func (ev *Event) SetTime(time float64) IEvent {
    ev.time = time
    return ev
}


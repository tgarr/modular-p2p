package utils

import (
    "container/heap"
    "sync"
)

const (
    INITIAL_QUEUE_SIZE            = 1000
)

// ==== interfaces ====

// event simulation interface
type IEventSimulation interface {
    Step() bool
    Now() float64
    Schedule(event IEvent,delay float64)
    GetHooks() *SimulationHooks
}

// ==== concrete structures ====

// wrapper for the priority queue
type queueItem struct {
    event IEvent
	time float64
    id uint64
}

// priority queue
type priorityQueue []*queueItem

/*
    Event simulation

    Implements: IEventSimulation
*/
type EventSimulation struct {
    currentTime float64
    nextID uint64
    queue priorityQueue
    hooks *SimulationHooks
    lock sync.RWMutex
}

// ==== factories ====

func NewEventSimulation() IEventSimulation {
    return &EventSimulation{
        currentTime:    0.0,
        nextID:         0,
        queue:          make(priorityQueue,0,INITIAL_QUEUE_SIZE),
        hooks:          NewSimulationHooks(),
        lock:           sync.RWMutex{},
    }
}

// ==== methods ====

// interface required by heap
func (queue priorityQueue) Len() int {
	return len(queue)
}

// interface required by heap
func (queue priorityQueue) Less(i,j int) bool {
	if queue[i].time != queue[j].time {
		return queue[i].time < queue[j].time
	}

    return queue[i].id < queue[j].id
}

// interface required by heap
func (queue priorityQueue) Swap(i,j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

// interface required by heap
func (queue *priorityQueue) Push(item interface{}) {
	*queue = append(*queue, item.(*queueItem))
}

// interface required by heap
func (queue *priorityQueue) Pop() interface{} {
	n := len(*queue)
	item := (*queue)[n-1]
	*queue = (*queue)[:n-1]
	return item
}

// trigger next event
func (sim *EventSimulation) Step() bool {
    sim.lock.Lock()

    // queue is empty
    if len(sim.queue) == 0 {
        sim.lock.Unlock()
        return false
    }

    // get the next event from the priority queue
    item := heap.Pop(&sim.queue).(*queueItem)
    event := item.event
    state := event.GetState()
    if state != EVENT_STATE_ABORTED {
        sim.currentTime = item.time
    }

    sim.lock.Unlock()
   
    // handle event
    if state != EVENT_STATE_ABORTED {
        event.SetTime(sim.currentTime)
        sim.hooks.EventPreTrigger(event) // pre trigger hook

        dest := event.GetDestination()
        handled := false
        if dest != nil {
            handled = dest.HandleEvent(event)
        }

        // open: what should the state if dest == nil?
        if handled {
            event.SetState(EVENT_STATE_HANDLED)
        } else {
            event.SetState(EVENT_STATE_NOTHANDLED)
        }

        sim.hooks.EventPostTrigger(event) // post trigger hook
    }

    return true
}

// current simulation time
func (sim *EventSimulation) Now() float64 {
    sim.lock.RLock()
    defer sim.lock.RUnlock()

    return sim.currentTime
}

func (sim *EventSimulation) Schedule(event IEvent,delay float64) {
    sim.lock.Lock()

    item := queueItem {
        event:      event,
        time:       sim.currentTime + delay,
        id:         sim.nextID,
    }
    
    heap.Push(&sim.queue,&item)
    sim.nextID++
    
    sim.lock.Unlock()

    event.SetTime(item.time)
    sim.hooks.EventScheduled(event) // scheduled hook
}

func (sim *EventSimulation) GetHooks() *SimulationHooks {
    return sim.hooks
}


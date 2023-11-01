package utils

// ==== interfaces ====

type IEventScheduledHandler interface {
    EventScheduled(ev IEvent)
}

type IEventPreTriggerHandler interface {
    EventPreTrigger(ev IEvent)
}

type IEventPostTriggerHandler interface {
    EventPostTrigger(ev IEvent)
}

// ==== concrete structures ====

/*
    Manage handlers for when events are scheduled, before they are triggered and
    after they are triggered. This is not thread-safe. The recommendation is
    that all hooks should be registered before starting the simulation.
*/
type SimulationHooks struct {
    eventScheduled map[uint16][]IEventScheduledHandler
    eventScheduledAll []IEventScheduledHandler

    preTrigger map[uint16][]IEventPreTriggerHandler
    preTriggerAll []IEventPreTriggerHandler

    postTrigger map[uint16][]IEventPostTriggerHandler
    postTriggerAll []IEventPostTriggerHandler
}

// ==== factories ====

func NewSimulationHooks() *SimulationHooks {
    return &SimulationHooks{
        eventScheduled:         make(map[uint16][]IEventScheduledHandler),
        eventScheduledAll:      make([]IEventScheduledHandler,0,4),
        preTrigger:             make(map[uint16][]IEventPreTriggerHandler),
        preTriggerAll:          make([]IEventPreTriggerHandler,0,4),
        postTrigger:            make(map[uint16][]IEventPostTriggerHandler),
        postTriggerAll:         make([]IEventPostTriggerHandler,0,4),
    }
}

// ==== methods ====

func (hooks *SimulationHooks) EventScheduled(ev IEvent) {
    for _,handler := range hooks.eventScheduledAll {
        handler.EventScheduled(ev)
    }

    tp := ev.GetType()
    if handlerList, ok := hooks.eventScheduled[tp]; ok {
        for _, handler := range handlerList {
            handler.EventScheduled(ev)
        }
    }
}

func (hooks *SimulationHooks) EventPreTrigger(ev IEvent) {
    for _,handler := range hooks.preTriggerAll {
        handler.EventPreTrigger(ev)
    }

    tp := ev.GetType()
    if handlerList, ok := hooks.preTrigger[tp]; ok {
        for _, handler := range handlerList {
            handler.EventPreTrigger(ev)
        }
    }
}

func (hooks *SimulationHooks) EventPostTrigger(ev IEvent) {
    for _,handler := range hooks.postTriggerAll {
        handler.EventPostTrigger(ev)
    }

    tp := ev.GetType()
    if handlerList, ok := hooks.postTrigger[tp]; ok {
        for _, handler := range handlerList {
            handler.EventPostTrigger(ev)
        }
    }
}

func (hooks *SimulationHooks) IsRegisteredScheduled(tp uint16,handler IEventScheduledHandler) bool {
    if handlerList, ok := hooks.eventScheduled[tp]; ok {
        for _, h := range handlerList {
            if h == handler {
                return true
            }
        }
    }

    return false
}

func (hooks *SimulationHooks) IsRegisteredScheduledAll(handler IEventScheduledHandler) bool {
    for _,h := range hooks.eventScheduledAll {
        if h == handler {
            return true
        }
    }

    return false
}

func (hooks *SimulationHooks) IsRegisteredPreTrigger(tp uint16,handler IEventPreTriggerHandler) bool {
    if handlerList, ok := hooks.preTrigger[tp]; ok {
        for _, h := range handlerList {
            if h == handler {
                return true
            }
        }
    }

    return false
}

func (hooks *SimulationHooks) IsRegisteredPreTriggerAll(handler IEventPreTriggerHandler) bool {
    for _,h := range hooks.preTriggerAll {
        if h == handler {
            return true
        }
    }

    return false
}

func (hooks *SimulationHooks) IsRegisteredPostTrigger(tp uint16,handler IEventPostTriggerHandler) bool {
    if handlerList, ok := hooks.postTrigger[tp]; ok {
        for _, h := range handlerList {
            if h == handler {
                return true
            }
        }
    }

    return false
}

func (hooks *SimulationHooks) IsRegisteredPostTriggerAll(handler IEventPostTriggerHandler) bool {
    for _,h := range hooks.postTriggerAll {
        if h == handler {
            return true
        }
    }

    return false
}

func (hooks *SimulationHooks) RegisterScheduled(tp uint16,handler IEventScheduledHandler) {
    if hooks.IsRegisteredScheduled(tp,handler) {
        return
    }
    
    if _, ok := hooks.eventScheduled[tp]; !ok {
        hooks.eventScheduled[tp] = make([]IEventScheduledHandler,0,16)
    }

    hooks.eventScheduled[tp] = append(hooks.eventScheduled[tp],handler)
}

func (hooks *SimulationHooks) RegisterScheduledAll(handler IEventScheduledHandler) {
    if hooks.IsRegisteredScheduledAll(handler) {
        return
    }
    
    hooks.eventScheduledAll = append(hooks.eventScheduledAll,handler)
}

func (hooks *SimulationHooks) RegisterPreTrigger(tp uint16,handler IEventPreTriggerHandler) {
    if hooks.IsRegisteredPreTrigger(tp,handler) {
        return
    }
    
    if _, ok := hooks.preTrigger[tp]; !ok {
        hooks.preTrigger[tp] = make([]IEventPreTriggerHandler,0,16)
    }

    hooks.preTrigger[tp] = append(hooks.preTrigger[tp],handler)
}

func (hooks *SimulationHooks) RegisterPreTriggerAll(handler IEventPreTriggerHandler) {
    if hooks.IsRegisteredPreTriggerAll(handler) {
        return
    }
    
    hooks.preTriggerAll = append(hooks.preTriggerAll,handler)
}

func (hooks *SimulationHooks) RegisterPostTrigger(tp uint16,handler IEventPostTriggerHandler) {
    if hooks.IsRegisteredPostTrigger(tp,handler) {
        return
    }
    
    if _, ok := hooks.postTrigger[tp]; !ok {
        hooks.postTrigger[tp] = make([]IEventPostTriggerHandler,0,16)
    }

    hooks.postTrigger[tp] = append(hooks.postTrigger[tp],handler)
}

func (hooks *SimulationHooks) RegisterPostTriggerAll(handler IEventPostTriggerHandler) {
    if hooks.IsRegisteredPostTriggerAll(handler) {
        return
    }
    
    hooks.postTriggerAll = append(hooks.postTriggerAll,handler)
}


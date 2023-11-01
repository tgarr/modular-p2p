package core

import (
    "blockchainlab/simulator/utils"
)

// ==== interfaces ====

// Basic interface followed by all components. Provide syntax sugar to common operations related to the simulation.
type ISimulationComponent interface {
    utils.IEventDestination
    
    // component life cycle
    Init(sim ISimulation,components ...ISimulationComponent)
    Finish()
    IsInitialized() bool

    // syntax sugar
    ScheduleEvent(event utils.IEvent,delay float64)
    GetSimulation() ISimulation
    GetTime() float64
    GetName() string
}

// ==== concrete sctructures ====

// Default component implementation. Most components should just incorporate it.
type DefaultComponent struct {
    sim ISimulation
    initialized bool
}

// ==== methods ====

func (comp *DefaultComponent) ScheduleEvent(event utils.IEvent,delay float64) {
    comp.GetSimulation().ScheduleEvent(event,delay)
}

func (comp *DefaultComponent) Init(sim ISimulation,components ...ISimulationComponent) {
    comp.sim = sim
    comp.initialized = true
}

func (comp *DefaultComponent) Finish() {
    comp.initialized = false
}

func (comp *DefaultComponent) HandleEvent(event utils.IEvent) bool {
    switch event.GetType() {
    case GLOBAL_NETWORK_EVENT_INIT:
        event.GetDestination().(ISimulationComponent).Init(event.GetData().(ISimulation)) 
        return true
    case NODE_EVENT_INIT:
        args := event.GetData().([]interface{})
        if len(args) == 2 {
            sim := args[0].(ISimulation)
            gnet := args[1].(IGlobalNetwork)
            event.GetDestination().(ISimulationComponent).Init(sim,gnet) 
            return true
        }
    case NODE_EVENT_FINISH:
        event.GetDestination().(ISimulationComponent).Finish()
        return true
    }

    return false
}

// ==== getters ====

func (comp *DefaultComponent) GetSimulation() ISimulation {
    return comp.sim
}

func (comp *DefaultComponent) GetTime() float64 {
    return comp.GetSimulation().GetTime()
}

func (comp *DefaultComponent) GetName() string {
    return "name_not_set"
}

func (comp *DefaultComponent) IsInitialized() bool {
    return comp.initialized
}


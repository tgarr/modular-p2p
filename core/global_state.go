package core

import (
    "blockchainlab/simulator/utils"
    "sync"
)

const (
    DEFAULT_GLOBAL_STATE_TAG                    = "default_global_state"
)

// ==== interfaces ====

// the global state implements a block registry and a K/V store
type ISimulationGlobalState interface {
    utils.IEventScheduledHandler
    utils.IEventPreTriggerHandler
    utils.IEventPostTriggerHandler

    // TODO block registry
    //PutBlock(block IBlock)
    //GetBlock(id uint64) IBlock

    // key/value store
    Put(key string,value interface{})
    Get(key string) interface{}

    Init(sim ISimulation)
    Finish()
}

// ==== concrete structures ====

/*
    Default implementation of the global state. Registers every new block by
    hooking into the new block event. New implementations of the gobal state
    should include this default implementation, unless they do not want to
    register all blocks.

    Implements: ISimulationGlobalState
*/
type SimulationGlobalState struct {
    sim ISimulation

    stateLock sync.RWMutex
    state map[string]interface{}

    //blockRegistry map[uint64]IBlock
    //blockLock sync.RWMutex
}

// ==== factories ====

var globalStateRegistry map[string]func() ISimulationGlobalState = make(map[string]func() ISimulationGlobalState)
var globalStateLogger utils.ISimulationLogger

func init(){
    RegisterGlobalState(DEFAULT_GLOBAL_STATE_TAG,NewSimulationGlobalState)
}

func RegisterGlobalState(key string, factory func() ISimulationGlobalState) {
    if _, ok := globalStateRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }   

    globalStateRegistry[key] = factory
}

func NewSimulationGlobalState() ISimulationGlobalState {
    if globalStateLogger == nil {
        globalStateLogger = utils.GetSimulationLogger(DEFAULT_GLOBAL_STATE_TAG)
    }

    return &SimulationGlobalState{
        sim:                nil,
        state:              make(map[string]interface{}),
        stateLock:          sync.RWMutex{},
        //blockRegistry:      make(map[uint64]IBlock),
        //blockLock:          sync.RWMutex{},
    }
}

func NewGlobalStateFromRegistry(key string) ISimulationGlobalState {
    if factory, ok := globalStateRegistry[key]; ok {
        return factory()
    }

    return nil
}

// ==== methods ====

func (global *SimulationGlobalState) Init(sim ISimulation){
    globalStateLogger.Debug("initializing: registering to new block events")
    global.sim = sim
    sim.GetHooks().RegisterPreTrigger(BLOCK_EVENT_NEW,global)    
}

func (global *SimulationGlobalState) Finish() {
}

func (global *SimulationGlobalState) EventScheduled(ev utils.IEvent) {
}

func (global *SimulationGlobalState) EventPreTrigger(ev utils.IEvent) {
    switch ev.GetType() {
    case BLOCK_EVENT_NEW:
        //block := ev.GetData().(IBlock)
        //global.PutBlock(block)
    }
}

func (global *SimulationGlobalState) EventPostTrigger(ev utils.IEvent) {
}

func (global *SimulationGlobalState) Put(key string,value interface{}) {
    global.stateLock.Lock()
    defer global.stateLock.Unlock()

    global.state[key] = value
}

func (global *SimulationGlobalState) Get(key string) interface{} {
    global.stateLock.RLock()
    defer global.stateLock.RUnlock()

    return global.state[key]
}

/*
func (global *SimulationGlobalState) PutBlock(block IBlock) {
    global.blockLock.Lock()
    defer global.blockLock.Unlock()

    hash := block.GetHash()
    if oldBlock, ok := global.blockRegistry[hash]; ok {
        if oldBlock != block {
            globalStateLogger.Warn("hash collision: old block created at %v by %v, new block created at %v by %v",oldBlock.GetTime(),oldBlock.GetCreator(),block.GetTime(),block.GetCreator())
        }
    } 

    globalStateLogger.Debug("registering new block at %v by %v",block.GetTime(),block.GetCreator())
    global.blockRegistry[hash] = block
}

func (global *SimulationGlobalState) GetBlock(hash uint64) IBlock {
    global.blockLock.RLock()
    defer global.blockLock.RUnlock()

    if block, ok := global.blockRegistry[hash]; ok {
        return block
    }

    return nil
}
*/


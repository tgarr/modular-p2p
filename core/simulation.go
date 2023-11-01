package core

import (
    "blockchainlab/simulator/utils"
    "sync"
    "errors"
    "fmt"
    "math/rand"
)

const (
    SIMULATION_TAG              = "simulation"
    DEFAULT_SIMULATION_NAME     = "defaultSim"
    DEFAULT_SEED                = 0
)

// ==== interfaces ====

/* 
    Interface for the simulation: not necessary in most cases, but it makes it
    possible to define a new simulation structure with different behavior.
*/
type ISimulation interface {
    utils.IEventDestination

    Run() error                                                 // run the simulation
    Stop()                                                      // request simulation to stop
    AddNode(node INode) ISimulation                             // add a node to the simulation
    RemoveNode(node_id uint32) error                            // remove a node from the simulation
    ScheduleEvent(event utils.IEvent,delay float64)             // schedule an event

    GetGlobalNetwork() IGlobalNetwork                           // get the global network for the simulation
    GetGlobalState() ISimulationGlobalState                     // get the global state
    IsRunning() bool                                            // check if the simulation is running
    GetEndCondition() IEndCondition                             // get the end condition
    GetNumNodes() uint32                                        // get the number of nodes in the simulation
    GetNode(node_id uint32) INode                               // get the node with the given id
    GetHooks() *utils.SimulationHooks                           // get hook manager
    GetTime() float64                                           // get simulation time
    GetName() string                                            // get simulation name
    GetRNG() *rand.Rand                                         // get random number generator

    SetGlobalNetwork(net IGlobalNetwork) ISimulation            // set the global network for the simulation
    SetGlobalState(state ISimulationGlobalState) ISimulation    // set the global network for the simulation
    SetEndCondition(end IEndCondition) ISimulation              // set the simulation end condition
}

// ==== concrete structures  ====

/* 
    Simulation implementation: main loop of the simulation, connect all nodes in
    a network and schedule events using a discreve event simulation. 

    Implements: ISimulation and IEventDestination interfaces.
*/
type Simulation struct {
    evSimulation utils.IEventSimulation
    network IGlobalNetwork
    state ISimulationGlobalState
    nodeMap map[uint32]INode
    running bool
    endCondition IEndCondition
    name string
    rng *rand.Rand
    
    nodeMapLock sync.RWMutex
    runningLock sync.RWMutex
}

// ==== factories ====

func init(){
    // config
    utils.ConfigSetDefault(SIMULATION_TAG + ".name", DEFAULT_SIMULATION_NAME)
    utils.ConfigSetDefault(SIMULATION_TAG + ".seed", DEFAULT_SEED)
}

var simLogger utils.ISimulationLogger

// basic factory for Simulation 
func NewSimulation() ISimulation {
    if simLogger == nil {
        simLogger = utils.GetSimulationLogger(SIMULATION_TAG)
    }

    config := utils.GetSimulationConfig()
    seed := config.GetInt64(SIMULATION_TAG + ".seed")
    return &Simulation {
        evSimulation:   utils.NewEventSimulation(),
        nodeMap:        make(map[uint32]INode),
        network:        nil,
        state:          nil,
        running:        false,
        endCondition:   nil,
        nodeMapLock:    sync.RWMutex{},
        runningLock:    sync.RWMutex{},
        name:           config.GetString(SIMULATION_TAG + ".name"),
        rng:            rand.New(rand.NewSource(seed)),
    }
}

// ==== methods ====

// run the simulation
func (sim *Simulation) Run() error {
    // check if everything is set, otherwise log and return an error
    if sim.GetEndCondition() == nil {
        simLogger.Error("end condition is not set")    
        return errors.New("end condition is not set")
    } else if sim.GetGlobalNetwork() == nil {
        simLogger.Error("global network is not set")
        return errors.New("global network is not set")
    }

    simLogger.Info("starting simulation %s with %d nodes",sim.GetName(),sim.GetNumNodes())

    // initialize components by scheduling init events
    // nodes are expected to connect to the global network and create the first non-init events
    sim.runningLock.Lock()
    sim.nodeMapLock.RLock()

    // initialize global state
    if sim.state != nil {
        sim.state.Init(sim)
    }

    // initialize global network
    sim.ScheduleEvent(utils.NewEvent(GLOBAL_NETWORK_EVENT_INIT,sim,sim.GetGlobalNetwork()),0)

    // initialize nodes
    args := []interface{}{sim,sim.GetGlobalNetwork()}
    for _, node := range sim.nodeMap {
        sim.ScheduleEvent(utils.NewEvent(NODE_EVENT_INIT,args,node),0)
    }
    
    sim.running = true
    sim.nodeMapLock.RUnlock()
    sim.runningLock.Unlock()

    // main loop
    var eventProcessed bool
    for sim.IsRunning() {
        // advance simulation
        eventProcessed = sim.evSimulation.Step()

        /*
            Check if simulation reached the end. It will continue only if all of the conditions below are true:
                - stop was not requested
                - an event was processed (event queue is not empty)
                - the end condition was not met
        */
        if !eventProcessed || sim.GetEndCondition().Check(sim) {
            sim.Stop()
        } 
    }

    // finish global state
    if sim.state != nil {
        sim.state.Finish()
    }

    simLogger.Info("simulation %s finished",sim.GetName())
    simLogger.Sync()
    return nil
}

// request simulation to stop
func (sim *Simulation) Stop(){
    simLogger.Info("requesting simulation %s to stop",sim.GetName())
    sim.ScheduleEvent(utils.NewEvent(SIMULATION_EVENT_STOP,nil,sim),0)
}

// add a new node to the simulation
func (sim *Simulation) AddNode(node INode) ISimulation {
    // mutex for concurrent access to the node map
    sim.nodeMapLock.Lock()
    defer sim.nodeMapLock.Unlock()

    new_id := uint32(len(sim.nodeMap) + 1) // node ids start from 1
    sim.nodeMap[new_id] = node
    node.SetID(new_id)

    // if already running, initialize the node
    if sim.IsRunning() {
        args := []interface{}{sim,sim.GetGlobalNetwork()}
        sim.ScheduleEvent(utils.NewEvent(NODE_EVENT_INIT,args,node),0)
    }

    return sim
}

func (sim *Simulation) RemoveNode(node_id uint32) error {
    sim.nodeMapLock.Lock()
    defer sim.nodeMapLock.Unlock()
    
    if node, ok := sim.nodeMap[node_id]; ok {
        sim.ScheduleEvent(utils.NewEvent(NODE_EVENT_FINISH,nil,node),0)
        delete(sim.nodeMap,node_id)
        return nil
    } else {
        return fmt.Errorf("node %d does not exist",node_id)
    }
}

func (sim *Simulation) ScheduleEvent(event utils.IEvent,delay float64) {
    sim.evSimulation.Schedule(event,delay)
}

func (sim *Simulation) HandleEvent(event utils.IEvent) bool {
    switch event.GetType() {
    case SIMULATION_EVENT_STOP:
        sim.runningLock.Lock()
        sim.running = false
        sim.runningLock.Unlock()
        return true
    case SIMULATION_EVENT_ADD_NODE:
        sim.AddNode(event.GetData().(INode))
        return true
    case SIMULATION_EVENT_REMOVE_NODE:
        sim.RemoveNode(event.GetData().(uint32))
        return true
    }

    return false
}

// ==== getters ====

func (sim *Simulation) GetHooks() *utils.SimulationHooks {
    return sim.evSimulation.GetHooks()
}

func (sim *Simulation) GetName() string {
    return sim.name
}

func (sim *Simulation) GetNumNodes() uint32 {
    sim.nodeMapLock.RLock()
    defer sim.nodeMapLock.RUnlock()

    return uint32(len(sim.nodeMap))
}

func (sim *Simulation) GetGlobalNetwork() IGlobalNetwork {
    return sim.network
}

func (sim *Simulation) GetGlobalState() ISimulationGlobalState {
    return sim.state
}

func (sim *Simulation) IsRunning() bool {
    sim.runningLock.RLock()
    defer sim.runningLock.RUnlock()

    return sim.running
}
    
func (sim *Simulation) GetEndCondition() IEndCondition {
    return sim.endCondition
}

func (sim *Simulation) GetTime() float64 {
    return sim.evSimulation.Now()
}

// returns the node with the given id
func (sim *Simulation) GetNode(node_id uint32) INode {
    sim.nodeMapLock.RLock()
    defer sim.nodeMapLock.RUnlock()
    
    if node, ok := sim.nodeMap[node_id]; ok {
         return node
    }

    return nil
}

func (sim *Simulation) GetRNG() *rand.Rand {
    return sim.rng
}

// ==== setters ====

func (sim *Simulation) SetGlobalNetwork(net IGlobalNetwork) ISimulation {
    sim.network = net
    return sim
}
func (sim *Simulation) SetGlobalState(state ISimulationGlobalState) ISimulation {
    sim.state = state
    return sim
}

func (sim *Simulation) SetEndCondition(end IEndCondition) ISimulation {
    sim.endCondition = end
    return sim
}


package behavior

import (
    "blockchainlab/simulator/core"
    "blockchainlab/simulator/utils"
)

const (
    DEFAULT_BEHAVIOR_TAG                        = "default_behavior"
)

// ==== concrete structures ====

/*
    Simple node behavior that just relays messages to all layers.
*/
type DefaultBehavior struct {
    core.DefaultComponent

    node core.INode
}

// ==== factories ====

var bLogger utils.ISimulationLogger =  nil

func NewNodeBehavior() core.INodeBehavior {
    if bLogger == nil {
        bLogger = utils.GetSimulationLogger(DEFAULT_BEHAVIOR_TAG)
    }

    return &DefaultBehavior {
        node:           nil,
    }
}

func init() {
    // register factory
    core.RegisterNodeBehavior(DEFAULT_BEHAVIOR_TAG,NewNodeBehavior)
}

// ==== methods ====

func (behavior *DefaultBehavior) Init(sim core.ISimulation,components ...core.ISimulationComponent){
    behavior.DefaultComponent.Init(sim)
    behavior.node = components[0].(core.INode)

    bLogger.Debug("node %d behavior initializing",behavior.node.GetID())
}

func (behavior *DefaultBehavior) MessageReceived(msg core.IMessage) bool {
    // TODO message received
    return true
}

// ==== getters ====

func (behavior *DefaultBehavior) GetName() string { 
    return DEFAULT_BEHAVIOR_TAG
}


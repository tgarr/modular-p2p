package node

import (
    "blockchainlab/simulator/core"
    "blockchainlab/simulator/utils"
    "fmt"
)

const (
    DEFAULT_NODE_TAG                 = "default_node"
)

// ==== concrete structures ====

/* 
    A simple node that follows the given layer implementations

    Implements: INode 
*/
type DefaultNode struct {
    core.DefaultComponent

    nodeID uint32
    nodeType uint16

    // layers
    nodeNetwork core.INodeNetwork
    behavior core.INodeBehavior

    // TODO other layers
    //state core.INodeState
    //consensus core.IConsensusProtocol
    applications []core.IApplication
}

// ==== factories ====

func init() {
    // config
    utils.ConfigSetDefault(DEFAULT_NODE_TAG + ".default_node_network",nil)
    utils.ConfigSetDefault(DEFAULT_NODE_TAG + ".default_behavior",nil)

    /* TODO other layers config
    utils.ConfigSetDefault(DEFAULT_NODE_TAG + ".default_state",nil)
    utils.ConfigSetDefault(DEFAULT_NODE_TAG + ".default_consensus",nil)
    */

    // register factory
    core.RegisterNode(DEFAULT_NODE_TAG,NewDefaultNode)
}

var nLogger utils.ISimulationLogger = nil

// factory for DefaultNode
func NewDefaultNode() core.INode {
    if nLogger == nil {
        nLogger = utils.GetSimulationLogger(DEFAULT_NODE_TAG)
    }

    return &DefaultNode{
        nodeID:         0,
        nodeType:       core.NODE_TYPE_FULL,
        nodeNetwork:    nil,
        behavior:       nil,
        // TODO other layers
        //state:          nil,
        //consensus:      nil,
        applications:   make([]core.IApplication, 0, 10),
    }
}

// ==== methods ====

// add a new application to the node
func (node *DefaultNode) AddApplication(app core.IApplication) core.INode {
    if app != nil {
        node.applications = append(node.applications,app)
        if node.IsInitialized() {
            sim := node.GetSimulation()
            n := sim.GetNode(node.GetID())
            app.Init(sim,n)
        }
    }
    return node
}

func (node *DefaultNode) HandleEvent(event utils.IEvent) bool {
    if node.DefaultComponent.HandleEvent(event) {
        return true
    }

    dest := event.GetDestination().(core.INode)
    switch event.GetType() {
    case core.NODE_ADD_APPLICATION:
        app := event.GetData().(core.IApplication)
        dest.AddApplication(app)
        return true
    }

    return false
}

func (node *DefaultNode) Init(sim core.ISimulation,components ...core.ISimulationComponent) {
    if node.IsInitialized() {
        nLogger.Error("node %d already initialized, doing it again",node.nodeID)
    }

    node.DefaultComponent.Init(sim)
    nLogger.Debug("node %d initializing",node.nodeID)

    // check if global network was provided
    if len(components) == 0 {
        panic("node initialization requires a global network")
    }

    gnet := components[0].(core.IGlobalNetwork)
    if gnet == nil {
        panic(fmt.Sprintf("node %d cannot initialize with a 'nil' global network",node.GetID()))
    }
    
    // set up stack: behavior and node network are mandatory, others are optional
    var layer core.ISimulationComponent
    config := utils.GetSimulationConfig()

    // node network
    layer = node.GetNodeNetwork()
    if layer == nil {
        nnetConf := config.GetString(DEFAULT_NODE_TAG + ".default_node_network")

        nnet := core.NewNodeNetworkFromRegistry(nnetConf)
        if nnet == nil {
            panic(fmt.Sprintf("node %d network not set: %v not registered",node.GetID(),nnetConf))
        } else {
            nLogger.Debug("node %d is using node network %v",node.GetID(),nnetConf)
            node.SetNodeNetwork(nnet)
            layer = nnet
        }
    }
    layer.Init(sim,node,gnet)
    
    // behavior
    layer = node.GetBehavior()
    if layer == nil {
        behaviorConf := config.GetString(DEFAULT_NODE_TAG + ".default_behavior")

        behavior := core.NewNodeBehaviorFromRegistry(behaviorConf)
        if behavior == nil {
            panic(fmt.Sprintf("node %d behavior not set: %v not registered",node.GetID(),behaviorConf))
        } else {
            nLogger.Debug("node %d is using behavior %v",node.GetID(),behaviorConf)
            node.SetBehavior(behavior)
            layer = behavior
        }
    }
    layer.Init(sim,node)

    /* TODO layers below are not implemented yet
    
    // state
    layer = node.GetNodeState()
    if layer == nil {
        stateConf := config.GetString(DEFAULT_NODE_TAG + ".default_state")

        ledger := core.NewLedgerFromRegistry(ledgerConf)
        if ledger == nil {
            nLogger.Debug(no )
            panic(fmt.Sprintf("node %d ledger not set: %v not registered",node.GetID(),ledgerConf))
        } else {
            nLogger.Debug("node %d is using ledger %v",node.GetID(),ledgerConf)
            node.SetLedger(ledger)
            layer = ledger
        }
    }
    layer.Init(sim,node)

    // consensus
    layer = node.GetConsensusProtocol()
    if layer == nil {
        consensusConf := config.GetString(DEFAULT_NODE_TAG + ".default_consensus")

        consensus := core.NewConsensusProtocolFromRegistry(consensusConf)
        if consensus == nil {
            panic(fmt.Sprintf("node %d consensus not set: %v not registered",node.GetID(),consensusConf))
        } else {
            nLogger.Debug("node %d is using consensus %v",node.GetID(),consensusConf)
            node.SetConsensusProtocol(consensus)
            layer = consensus
        }
    }
    layer.Init(sim,node)
    
    */

    // application
    for _, app := range node.GetApplications() {
        app.Init(sim,node)
    }
}

func (node *DefaultNode) Finish() {
    nLogger.Debug("node %d finishing",node.nodeID)
    node.GetNodeNetwork().Finish()
    node.DefaultComponent.Finish()
}

// ==== getters ====

func (node *DefaultNode) GetID() uint32 {
    return node.nodeID
}

func (node *DefaultNode) GetType() uint16 {
    return node.nodeType
}


func (node *DefaultNode) GetNodeNetwork() core.INodeNetwork {
    return node.nodeNetwork
}

func (node *DefaultNode) GetBehavior() core.INodeBehavior {
    return node.behavior
}

/* TODO getters for other layers
func (node *DefaultNode) GetLedger() core.ILedger {
    return node.ledger
}

func (node *DefaultNode) GetConsensusProtocol() core.IConsensusProtocol {
    return node.consensus
}
*/

func (node *DefaultNode) GetApplications() []core.IApplication {
    return node.applications
}

func (node *DefaultNode) GetName() string {
    return DEFAULT_NODE_TAG
}

// ==== setters ====

func (node *DefaultNode) SetID(id uint32) core.INode {
    node.nodeID = id
    return node
}

func (node *DefaultNode) SetType(tp uint16) core.INode {
    node.nodeType = tp
    return node
}

func (node *DefaultNode) SetNodeNetwork(net core.INodeNetwork) core.INode {
    node.nodeNetwork = net
    return node
}

func (node *DefaultNode) SetBehavior(behavior core.INodeBehavior) core.INode {
    node.behavior = behavior
    return node
}

/* TODO setters for other layers
func (node *DefaultNode) SetLedger(ledger core.ILedger) core.INode {
    node.ledger = ledger
    return node
}

func (node *DefaultNode) SetConsensusProtocol(consensus core.IConsensusProtocol) core.INode {
    node.consensus = consensus
    return node
}
*/


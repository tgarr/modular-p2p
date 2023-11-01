package node_network

import (
    "blockchainlab/simulator/core"
    "blockchainlab/simulator/utils"
    "sync"
)

const (
    DEFAULT_NODE_NETWORK_TAG                        = "default_node_network"
)

// ==== concrete structures ====

/*
    Simple implementation of a node network layer. It keeps a list of neighbors
    that can be added or removed by the node behavior (no protocol for this is
    implemented). Any message received is relayed to the node behavior.

    Implements: INodeNetwork
*/
type DefaultNodeNetwork struct {
    core.DefaultComponent

    node core.INode
    globalNet core.IGlobalNetwork
    neighbors []uint32
    neighborLock sync.RWMutex
}

// ==== factories ====

var nnetLogger utils.ISimulationLogger =  nil

func NewNodeNetwork() core.INodeNetwork {
    if nnetLogger == nil {
        nnetLogger = utils.GetSimulationLogger(DEFAULT_NODE_NETWORK_TAG)
    }

    return &DefaultNodeNetwork {
        node:           nil,
        globalNet:      nil,
        neighbors:      make([]uint32,0,10),
        neighborLock:   sync.RWMutex{},
    }
}

func init() {
    // register factory
    core.RegisterNodeNetwork(DEFAULT_NODE_NETWORK_TAG,NewNodeNetwork)
}

// ==== methods ====

func (net *DefaultNodeNetwork) Init(sim core.ISimulation,components ...core.ISimulationComponent){
    net.DefaultComponent.Init(sim)
    
    if len(components) < 2 {
        panic("DefaultNodeNetwork requires a node and a global network to initialize")
    }

    net.node = components[0].(core.INode)
    nnetLogger.Debug("node %d network initializing",net.node.GetID())

    gnet := components[1].(core.IGlobalNetwork)
    if gnet == nil {
        panic("cannot connect to <nil> global network")
    }
   
    // connect
    net.ScheduleEvent(utils.NewEvent(core.NODE_NETWORK_EVENT_CONNECT,gnet,net),0)
}

func (net *DefaultNodeNetwork) Finish() {
    net.ScheduleEvent(utils.NewEvent(core.NODE_NETWORK_EVENT_DISCONNECT,nil,net),0)
    net.DefaultComponent.Finish()
}

func (net *DefaultNodeNetwork) HandleEvent(event utils.IEvent) bool {
    if net.DefaultComponent.HandleEvent(event) {
        return true
    }

    dest := event.GetDestination().(core.INodeNetwork)
    switch event.GetType() {
    case core.NODE_NETWORK_EVENT_MESSAGE_RECEIVED:
        return dest.MessageReceived(event.GetData().(core.IMessage))
    case core.NODE_NETWORK_EVENT_CONNECT:
        gnet := event.GetData().(core.IGlobalNetwork)
        dest.Disconnect()
        dest.Connect(gnet)
        return true
    case core.NODE_NETWORK_EVENT_DISCONNECT:
        dest.Disconnect()
        return true
    default:
        nnetLogger.Debug("unknown event %d",event.GetType())
    }

    return false
}

func (net *DefaultNodeNetwork) Connect(gnet core.IGlobalNetwork) {
    net.globalNet = gnet
    net.globalNet.Connect(net.node)
}

func (net *DefaultNodeNetwork) Disconnect() {
    gnet := net.GetGlobalNetwork()
    if gnet != nil {
        net.globalNet.Disconnect(net.node)
    }
}

func (net *DefaultNodeNetwork) IsConnected() bool {
    return net.GetGlobalNetwork().IsConnected(net.node)
}

func (net *DefaultNodeNetwork) SendMessage(msg core.IMessage) {
    if net.IsConnected() {
        gnet := net.GetGlobalNetwork()
        net.ScheduleEvent(utils.NewEvent(core.GLOBAL_NETWORK_EVENT_SEND_MESSAGE,msg,gnet),0)
    } else {
        nnetLogger.Debug("cannot send message: node %d is disconneted",net.node.GetID())
    }
}

func (net *DefaultNodeNetwork) MessageReceived(msg core.IMessage) bool {
    nnetLogger.Debug("node %d received message %d from %d",net.node.GetID(),msg.GetTag(),msg.GetSender())
    return net.node.GetBehavior().MessageReceived(msg)
}

func (net *DefaultNodeNetwork) SendNode(tag int32,data interface{},target uint32) {
    msg := core.NewP2PMessage(data,net.node.GetID(),target)
    msg.SetTag(tag)
    
    nnetLogger.Debug("node %d sending message %d to node %d",net.node.GetID(),tag,target)
    net.SendMessage(msg)
}

func (net *DefaultNodeNetwork) SendBroadcast(tag int32,data interface{}) {
    msg := core.NewBroadcastMessage(data,net.node.GetID())
    msg.SetTag(tag)
    
    nnetLogger.Debug("node %d broadcasting message %d",net.node.GetID(),tag)
    net.SendMessage(msg)
}

func (net *DefaultNodeNetwork) SendNeighbors(tag int32,data interface{}) {
    msg := core.NewP2PMessageNodes(data,net.node.GetID(),net.GetNeighbors())
    msg.SetTag(tag)

    nnetLogger.Debug("node %d sending message %d to neighbors",net.node.GetID(),tag)
    net.SendMessage(msg)
}

func (net *DefaultNodeNetwork) AddNeighbor(nodeID uint32){
    if net.IsNeighbor(nodeID) {
        return
    }
    
    net.neighborLock.Lock()
    defer net.neighborLock.Unlock()

    net.neighbors = append(net.neighbors,nodeID)
}

func (net *DefaultNodeNetwork) RemoveNeighbor(nodeID uint32){
    if !net.IsNeighbor(nodeID) {
        return
    }

    net.neighborLock.Lock()
    defer net.neighborLock.Unlock()

    last := len(net.neighbors) - 1
    for i,n := range net.neighbors {
        if n == nodeID {
            net.neighbors[i] = net.neighbors[last]
            net.neighbors = net.neighbors[:last]
            break
        }
    }
}

func (net *DefaultNodeNetwork) IsNeighbor(nodeID uint32) bool {
    net.neighborLock.RLock()
    defer net.neighborLock.RUnlock()
    
    for _, n := range net.neighbors {
        if nodeID == n {
            return true
        }
    }

    return false
}

// ==== getters ====

func (net *DefaultNodeNetwork) GetNeighbors() []uint32 {
    net.neighborLock.RLock()
    defer net.neighborLock.RUnlock()

    nList := make([]uint32,len(net.neighbors))
    copy(nList,net.neighbors)
    return nList
}

func (net *DefaultNodeNetwork) GetNumNeighbors() uint32 {
    net.neighborLock.RLock()
    defer net.neighborLock.RUnlock()

    return uint32(len(net.neighbors))
}

func (net *DefaultNodeNetwork) GetGlobalNetwork() core.IGlobalNetwork {
    return net.globalNet
}

func (net *DefaultNodeNetwork) GetName() string {
    return DEFAULT_NODE_NETWORK_TAG
}


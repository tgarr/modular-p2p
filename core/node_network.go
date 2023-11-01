package core

// ==== interfaces ====

/* 
    This implements a node network layer. Each node has its own instance of INodelNetwork, which are
    connected to a single IGlobalNetwork.
*/
type INodeNetwork interface {
    ISimulationComponent
    
    SendNeighbors(tag int32,data interface{})
    SendNode(tag int32,data interface{},target uint32)
    SendBroadcast(tag int32,data interface{})
    
    SendMessage(msg IMessage)
    MessageReceived(msg IMessage) bool
   
    Connect(net IGlobalNetwork)
    Disconnect()
    IsConnected() bool

    GetGlobalNetwork() IGlobalNetwork
    
    AddNeighbor(nodeID uint32)
    RemoveNeighbor(nodeID uint32)
    GetNeighbors() []uint32
    GetNumNeighbors() uint32
    IsNeighbor(nodeID uint32) bool
}

// ==== factories ====

var nodeNetworkRegistry map[string]func() INodeNetwork = make(map[string]func() INodeNetwork)

func RegisterNodeNetwork(key string, factory func() INodeNetwork) {
    if _, ok := nodeNetworkRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }   

    nodeNetworkRegistry[key] = factory
}

func NewNodeNetworkFromRegistry(key string) INodeNetwork {
    if factory, ok := nodeNetworkRegistry[key]; ok {
        return factory()
    }   

    return nil 
}


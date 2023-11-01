package core

// ==== interfaces ====

/*
    Each simulation should have only one instance of IGlobalNetwork, which
    interconnects the network layers (INodeNetwork) of each node. A message can
    be sent in two modes: broadcast and p2p. In broadcast mode, a message is
    sent to multiple nodes (all or a subset), while in p2p mode a message is
    sent to just one specific node. The mode is determined by IMessageDelivery.
*/
type IGlobalNetwork interface {
    ISimulationComponent

    SendMessage(msg IMessage) IGlobalNetwork        // send a message
    Connect(node INode) IGlobalNetwork              // connect the given node to the global network
    Disconnect(node INode) IGlobalNetwork           // disconnect the given node
    IsConnected(node INode) bool                    // check if a node is connected
    EnableBroadcast(node INode) IGlobalNetwork      // start sending global broadcasts to given node
    DisableBroadcast(node INode) IGlobalNetwork     // stop sending global broadcasts to given node
    IsBroadcastEnabled(node INode) bool             // check if the given node is receiving global broadcast
}

// ==== factories ====

var globalNetworkRegistry map[string]func() IGlobalNetwork = make(map[string]func() IGlobalNetwork)

func RegisterGlobalNetwork(key string, factory func() IGlobalNetwork) {
    if _, ok := globalNetworkRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }

    globalNetworkRegistry[key] = factory
}

func NewGlobalNetworkFromRegistry(key string) IGlobalNetwork {
    if factory, ok := globalNetworkRegistry[key]; ok {
        return factory()
    }

    return nil 
}


package core

// node types
const (
    NODE_TYPE_FULL                  = 0 // full node
    NODE_TYPE_CLIENT                = 1 // client node
    NODE_TYPE_LIGHT                 = 2 // lightweight node
    NODE_TYPE_ARCHIVE               = 3 // archive node
    NODE_TYPE_IOT                   = 4 // Internet of Things node
    NODE_TYPE_WALLET                = 5 // Wallet software
)

// ==== interfaces ====

/*
    Nodes are the basic building block of a simulation. The blockchain stack is implemented within
    each node.  The main logic of a node is implemented by INodeBehavior. Ideally, behavior consists
    of the main protocol high-level logic, while an IApplication uses the main protocol to do something.
    
    For example, if implementing a Bitcoin full node, INodeBehavior should implement the whole
    protocol followed by a full node. In a client node, INodeBehavior should just route the events
    received, and the client logic should be in an IApplication.
*/
type INode interface {
    ISimulationComponent

    GetID() uint32                                              // get node id
    GetType() uint16                                            // get node type
    
    SetID(id uint32) INode                                      // set node id (which is decided by the simulation)
    SetType(tp uint16) INode                                    // set node type

    // blockchain layers
    
    // node local network layer (connected to the global)
    SetNodeNetwork(network INodeNetwork) INode
    GetNodeNetwork() INodeNetwork

    // node behavior
    SetBehavior(behavior INodeBehavior) INode                   // set node behavior
    GetBehavior() INodeBehavior                                 // get node behavior
    
    /* TODO omitting stuff not implemented yet
    SetNodeState(state INodeState) INode                        // set node state
    GetNodeState() INodeState                                   // get node state
    
    SetConsensusProtocol(consensus IConsensusProtocol) INode    // set consensus protocol
    GetConsensusProtocol() IConsensusProtocol                   // get consensus protocol
    */

    // applications running on the node
    AddApplication(app IApplication) INode
    GetApplications() []IApplication
}

// ==== factories ====

var nodeRegistry map[string]func() INode = make(map[string]func() INode)

func RegisterNode(key string, factory func() INode) {
    if _, ok := nodeRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }

    nodeRegistry[key] = factory
}

func NewNodeFromRegistry(key string) INode {
    if factory, ok := nodeRegistry[key]; ok {
        return factory()
    }

    return nil 
}


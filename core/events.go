package core

/*
    Event types: different packages should use different ranges of values:
        0-1000:         core
        1001-2000:      layers.global_network 
        2001-3000:      layers.node
        3001-4000:      layers.node.node_network
        5001-6000:      layers.node.behavior
        6001-7000:      layers.node.consensus
        7001-8000:      
        8001-9000:      
        9001-10000:      
        9001-10000:      
        10001-11000:      
        11001-12000:      

        30001-65535:    user-defined packages
*/

const (
    // general
    EVENT_GENERIC                                       = 0     // when type is not important

    // simulation
    SIMULATION_EVENT_STOP                               = 1     // stop simulation
    SIMULATION_EVENT_ADD_NODE                           = 2     // add a node
    SIMULATION_EVENT_REMOVE_NODE                        = 3     // remove a node

    // node
    NODE_EVENT_INIT                                     = 10    // init node
    NODE_EVENT_FINISH                                   = 11    // finalize node
    NODE_ADD_APPLICATION                                = 12    // start an application on the node
    
    // global network
    GLOBAL_NETWORK_EVENT_INIT                           = 20    // init global network
    GLOBAL_NETWORK_EVENT_SEND_MESSAGE                   = 21    // send message
    
    // note network
    NODE_NETWORK_EVENT_MESSAGE_RECEIVED                 = 30    // message received from global network
    NODE_NETWORK_EVENT_CONNECT                          = 31    // connect to global network
    NODE_NETWORK_EVENT_DISCONNECT                       = 32    // disconnect from global network

    // block generation
    BLOCK_EVENT_NEW                                     = 40    // new block created
)


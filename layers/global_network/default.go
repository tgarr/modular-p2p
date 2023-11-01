package global_network

import (
    "blockchainlab/simulator/core"
    "blockchainlab/simulator/utils"
    "strconv"
    "sync"
    "math/rand"
)

// general
const (
    DEFAULT_GNET_TAG                            = "default_global_network"          // tag for registry and log
)

// Default distributions. See utils/sampler.go for configuration parameters.
const (
    DEFAULT_BROADCAST_DISTRIBUTION              = "exponential"                     // default distribution for broadcast messages
    DEFAULT_P2P_DISTRIBUTION                    = "normal"                          // default distribution for p2p messages
)

var DEFAULT_EXPONENTIAL_CONFIG                  = []float64{0.109,0.01,-1.0}        // [average,min,max]
var DEFAULT_UNIFORM_CONFIG                      = []float64{0.01,0.5}               // [min,max]
var DEFAULT_NORMAL_CONFIG                       = []float64{0.05,0.05,0.01,0.5}     // [average,stddev,min,max]
var DEFAULT_ZIPF_CONFIG                         = []float64{}                       // # TODO

// ==== concrete structures  ====

/*
    Simple global network that uses statistical distributions to compute the
    propagation delay of p2p and broadcast messages. 

    Implements: IGlobalNetwork
*/
type DefaultGlobalNetwork struct {
    core.DefaultComponent

    nodeMap map[uint32]core.INode
    nodeTypeMap map[uint16][]core.INode
    nodeMapLock sync.RWMutex
    globalBroadcastActive map[uint32]bool

    broadcastSampler utils.ISimulationSampler
    p2pSampler utils.ISimulationSampler

    broadcastDist string
    broadcastConfig []string
    p2pDist string
    p2pConfig []string
}

// ==== factories ====

func init() {
    // config
    utils.ConfigSetDefault(DEFAULT_GNET_TAG + ".broadcast_distribution", DEFAULT_BROADCAST_DISTRIBUTION)
    utils.ConfigSetDefault(DEFAULT_GNET_TAG + ".p2p_distribution", DEFAULT_P2P_DISTRIBUTION)
    
    utils.ConfigSetDefault(DEFAULT_GNET_TAG + ".broadcast_config", nil)
    utils.ConfigSetDefault(DEFAULT_GNET_TAG + ".p2p_config", nil)

    // register factory
    core.RegisterGlobalNetwork(DEFAULT_GNET_TAG,NewDefaultGlobalNetwork)
}

// build a sampler accoring to configuration
func buildSampler(distName string, distConfig []string,rng *rand.Rand) utils.ISimulationSampler {
    var configValues []float64 = nil

    if distConfig != nil {
        configValues = make([]float64,0,4)
        for _, str := range distConfig {
            f, err := strconv.ParseFloat(str,64)
            if err != nil {
                panic(err)
            }
            configValues = append(configValues,f)
        }
    } else {
        switch distName {
        case "exponential":
            configValues = DEFAULT_EXPONENTIAL_CONFIG
        case "uniform":
            configValues = DEFAULT_UNIFORM_CONFIG
        case "normal":
            configValues = DEFAULT_NORMAL_CONFIG
        case "zipf":
            configValues = DEFAULT_ZIPF_CONFIG
        default:
            configValues = nil
        }
    }

    return utils.NewSampler(distName,configValues,rng)
}

var gnetLogger utils.ISimulationLogger = nil

func NewDefaultGlobalNetwork() core.IGlobalNetwork {
    config := utils.GetSimulationConfig()

    broadcastDist := config.GetString(DEFAULT_GNET_TAG + ".broadcast_distribution") 
    broadcastConfig := config.GetStringSlice(DEFAULT_GNET_TAG + ".broadcast_config")

    p2pDist := config.GetString(DEFAULT_GNET_TAG + ".p2p_distribution")
    p2pConfig := config.GetStringSlice(DEFAULT_GNET_TAG + ".p2p_config")

    // gnetLogger
    if gnetLogger == nil {
        gnetLogger = utils.GetSimulationLogger(DEFAULT_GNET_TAG)
    }
   
    return &DefaultGlobalNetwork{
        nodeMap:                    make(map[uint32]core.INode),
        nodeTypeMap:                make(map[uint16][]core.INode),
        nodeMapLock:                sync.RWMutex{},
        globalBroadcastActive:      make(map[uint32]bool),
        broadcastSampler:           nil,
        p2pSampler:                 nil,
        broadcastDist:              broadcastDist,
        broadcastConfig:            broadcastConfig,
        p2pDist:                    p2pDist,
        p2pConfig:                  p2pConfig,
    }
}

// ==== methods ====

func (net *DefaultGlobalNetwork) Init(sim core.ISimulation,components ...core.ISimulationComponent){
    net.DefaultComponent.Init(sim)

    rng := sim.GetRNG()
    net.broadcastSampler = buildSampler(net.broadcastDist,net.broadcastConfig,rng)
    net.p2pSampler = buildSampler(net.p2pDist,net.p2pConfig,rng)

    gnetLogger.Debug("initializing with p2pSampler=%v and broadcastSampler=%v",net.p2pSampler.GetDistName(),net.broadcastSampler.GetDistName())
}

func (net *DefaultGlobalNetwork) HandleEvent(event utils.IEvent) bool {
    if net.DefaultComponent.HandleEvent(event) {
        return true
    }

    dest := event.GetDestination().(core.IGlobalNetwork)
    switch event.GetType() {
    case core.GLOBAL_NETWORK_EVENT_SEND_MESSAGE:
        msg := event.GetData().(core.IMessage)
        dest.SendMessage(msg)
        return true
    }

    return false
}

func (net *DefaultGlobalNetwork) SendMessage(msg core.IMessage) core.IGlobalNetwork {
    var sampler utils.ISimulationSampler
    isBroadcast := false
    dtype := 0 // 0: specific nodes, 1: types, 2: excl. types
    
    msg.SetTime(net.GetTime())
    delivery := msg.GetDelivery()
    switch delivery.GetDeliveryType() {
    case core.MESSAGE_DELIVERY_TYPE_P2P_NODES,core.MESSAGE_DELIVERY_TYPE_P2P_NODE_TYPES,core.MESSAGE_DELIVERY_TYPE_P2P_NODE_TYPES_EXCEPT:
        sampler = net.p2pSampler
        isBroadcast = false
    default:
        sampler = net.broadcastSampler
        isBroadcast = true
    }

    switch delivery.GetDeliveryType() {
    case core.MESSAGE_DELIVERY_TYPE_P2P_NODES,core.MESSAGE_DELIVERY_TYPE_BROADCAST_NODES:
        dtype = 0
    case core.MESSAGE_DELIVERY_TYPE_P2P_NODE_TYPES,core.MESSAGE_DELIVERY_TYPE_BROADCAST_NODE_TYPES:
        dtype = 1
    case core.MESSAGE_DELIVERY_TYPE_P2P_NODE_TYPES_EXCEPT,core.MESSAGE_DELIVERY_TYPE_BROADCAST_NODE_TYPES_EXCEPT:
        dtype = 2
    default:
        panic("delivery type not supported")
    }

    net.nodeMapLock.RLock()
    targets := delivery.GetDeliveryTargets()
    if targets == nil { // all nodes
        for _, node := range net.nodeMap {
            if node.GetID() == msg.GetSender() { // no loopback
                continue
            }

            if isBroadcast && !net.globalBroadcastActive[node.GetID()] { // global broadcast
                continue
            }
    
            ev := utils.NewEvent(core.NODE_NETWORK_EVENT_MESSAGE_RECEIVED,msg,node.GetNodeNetwork())
            delay := sampler.Sample()
            net.ScheduleEvent(ev,delay)
        }
    } else {
        switch dtype {
        case 0: // specific nodes
            for _, nodeID := range targets {
                if node,ok := net.nodeMap[nodeID]; ok {
                    ev := utils.NewEvent(core.NODE_NETWORK_EVENT_MESSAGE_RECEIVED,msg,node.GetNodeNetwork())
                    delay := sampler.Sample()
                    net.ScheduleEvent(ev,delay)
                } else {
                    gnetLogger.Debug("node %d not connected",nodeID)
                }
            }
        case 1: // nodes of specified types
            for _, tp := range targets {
                if nodeList, ok := net.nodeTypeMap[uint16(tp)]; ok {
                    for _, node := range nodeList {
                        if node.GetID() == msg.GetSender() { // no loopback
                            continue
                        }

                        ev := utils.NewEvent(core.NODE_NETWORK_EVENT_MESSAGE_RECEIVED,msg,node.GetNodeNetwork())
                        delay := sampler.Sample()
                        net.ScheduleEvent(ev,delay)
                    }
                }
            }
        default: // nodes not of specified types
            for tp, nodeList := range net.nodeTypeMap {
                found := false
                for _, targetType := range targets {
                    if uint32(tp) == targetType {
                        found = true
                        break
                    }
                }

                if !found {
                    for _, node := range nodeList {
                        if node.GetID() == msg.GetSender() { // no loopback
                            continue
                        }

                        ev := utils.NewEvent(core.NODE_NETWORK_EVENT_MESSAGE_RECEIVED,msg,node.GetNodeNetwork())
                        delay := sampler.Sample()
                        net.ScheduleEvent(ev,delay)
                    }
                }
            }
        }
    }
    net.nodeMapLock.RUnlock()

    return net
}

func (net *DefaultGlobalNetwork) Connect(node core.INode) core.IGlobalNetwork {
    net.nodeMapLock.Lock()
    defer net.nodeMapLock.Unlock()

    if _, ok := net.nodeMap[node.GetID()]; !ok {
        // node map
        net.nodeMap[node.GetID()] = node

        // node type map
        tp := node.GetType()
        if _, ok = net.nodeTypeMap[tp]; !ok {
            net.nodeTypeMap[tp] = make([]core.INode,0,100)
        }
        net.nodeTypeMap[tp] = append(net.nodeTypeMap[tp],node)

        // active broadcast map: default value based on type
        switch tp {
        case core.NODE_TYPE_FULL,core.NODE_TYPE_ARCHIVE:
            net.globalBroadcastActive[node.GetID()] = true
        default:
            net.globalBroadcastActive[node.GetID()] = false
        }

        gnetLogger.Debug("node %d connected",node.GetID())
    } else {
        gnetLogger.Debug("node %d already connected",node.GetID())
    }

    return net
}

func (net *DefaultGlobalNetwork) Disconnect(node core.INode) core.IGlobalNetwork {
    net.nodeMapLock.Lock()
    defer net.nodeMapLock.Unlock()
    
    if _, ok := net.nodeMap[node.GetID()]; ok {
        delete(net.nodeMap,node.GetID())
       
        tp := node.GetType()
        last := len(net.nodeTypeMap[tp]) - 1
        for i,n := range net.nodeTypeMap[tp] {
            if n == node {
                net.nodeTypeMap[tp][i] = net.nodeTypeMap[tp][last]
                net.nodeTypeMap[tp] = net.nodeTypeMap[tp][:last]
                break
            }
        }

        delete(net.globalBroadcastActive,node.GetID())
        gnetLogger.Debug("node %d disconnected",node.GetID())
    } else {
        gnetLogger.Debug("node %d not connected",node.GetID())
    }

    return net
}

func (net *DefaultGlobalNetwork) IsConnected(node core.INode) bool {
    net.nodeMapLock.RLock()
    defer net.nodeMapLock.RUnlock()

    _, ok := net.nodeMap[node.GetID()]

    return ok
}

func (net *DefaultGlobalNetwork) EnableBroadcast(node core.INode) core.IGlobalNetwork {
    net.nodeMapLock.RLock()
    defer net.nodeMapLock.RUnlock()

    if _, ok := net.nodeMap[node.GetID()]; ok {
        net.globalBroadcastActive[node.GetID()] = true
    } else {
        gnetLogger.Debug("node %d not connected",node.GetID())
    }
    
    return net
}

func (net *DefaultGlobalNetwork) DisableBroadcast(node core.INode) core.IGlobalNetwork {
    net.nodeMapLock.RLock()
    defer net.nodeMapLock.RUnlock()

    if _, ok := net.nodeMap[node.GetID()]; ok {
        net.globalBroadcastActive[node.GetID()] = false
    } else {
        gnetLogger.Debug("node %d not connected",node.GetID())
    }

    return net
}

func (net *DefaultGlobalNetwork) IsBroadcastEnabled(node core.INode) bool {
    active, ok := net.globalBroadcastActive[node.GetID()]
    return active && ok
}

// ==== getters ====

func (net *DefaultGlobalNetwork) GetName() string {
    return DEFAULT_GNET_TAG
}


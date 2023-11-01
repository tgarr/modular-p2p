package core

// ==== interfaces ====

// behavior of a node
type INodeBehavior interface {
    ISimulationComponent

    MessageReceived(msg IMessage) bool
    //BlockCreated(block IBlock)
}

// ==== factories ====

var behaviorRegistry map[string]func() INodeBehavior = make(map[string]func() INodeBehavior)

func RegisterNodeBehavior(key string, factory func() INodeBehavior) {
    if _, ok := behaviorRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }

    behaviorRegistry[key] = factory
}

func NewNodeBehaviorFromRegistry(key string) INodeBehavior {
    if factory, ok := behaviorRegistry[key]; ok {
        return factory()
    }

    return nil 
}


package core

// TODO consensus layer

// ==== interfaces ====

// consensus protocol followed by nodes
type IConsensusProtocol interface {
    ISimulationComponent
}

// ==== factories ====

var consensusRegistry map[string]func() IConsensusProtocol = make(map[string]func() IConsensusProtocol)

func RegisterConsensusProtocol(key string, factory func() IConsensusProtocol) {
    if _, ok := consensusRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }

    consensusRegistry[key] = factory
}

func NewConsensusProtocolFromRegistry(key string) IConsensusProtocol {
    if factory, ok := consensusRegistry[key]; ok {
        return factory()
    }

    return nil 
}


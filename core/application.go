package core

// ==== interfaces ====

// interface for an application running on a node
type IApplication interface {
    ISimulationComponent
}

// ==== factories ====

var applicationRegistry map[string]func() IApplication = make(map[string]func() IApplication)

func RegisterApplication(key string, factory func() IApplication) {
    if _, ok := applicationRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }

    applicationRegistry[key] = factory
}

func NewApplicationFromRegistry(key string) IApplication {
    if factory, ok := applicationRegistry[key]; ok {
        return factory()
    }

    return nil
}


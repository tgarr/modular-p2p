package core

import (
    "sync"
    "blockchainlab/simulator/utils"
)

const (
    SIMULATION_MEASUREMENTS_TAG                 = "measurements"

    SIMULATION_MEASUREMENT_TYPE_EVENT           = 0
    SIMULATION_MEASUREMENT_TYPE_TAG             = 1

    SIMULATION_MEASUREMENT_INITIAL_SIZE         = 65535
)

// ==== interfaces ====

/* 
    A measurement module takes custom measurements during and/or after the
    simulation. Each module writes its results to a separate json file.
*/
type ISimulationMeasurementModule interface {
    Init(sim ISimulation)                           // initialize measurement module
    Tag(id uint64,time float64,extra string)        //  
    GetFinalResult() interface{}                    // struct with final result of the module
    GetOutputPath() string                          // path of the output json file
}

// ==== concrete structures ====

type RawMeasurementEntry struct {
    tp uint8
    id uint64
    time float64
    extra string
}

type SimulationMeasurements struct {
    sim ISimulation
    lock sync.RWMutex
    outputPath string
    entries []RawMeasurementEntry
}

// ==== factories ====

var measurementModuleRegistry map[string]func() ISimulationMeasurementModule = make(map[string]func() ISimulationMeasurementModule)
var measurementLogger utils.ISimulationLogger

func init(){
    RegisterGlobalState(DEFAULT_GLOBAL_STATE_TAG,NewSimulationGlobalState)
}

func RegisterMeasurementModule(key string, factory func() ISimulationMeasurementModule) {
    if _, ok := measurementModuleRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }

    measurementModuleRegistry[key] = factory
}

func NewSimulationMeasurements() *SimulationMeasurements {
    if measurementLogger == nil {
        measurementLogger = utils.GetSimulationLogger(DEFAULT_GLOBAL_STATE_TAG)
    }

    return &SimulationMeasurements{
        sim:                nil,
        lock:               sync.RWMutex{},
        state:              make(map[string]interface{}),
        stateLock:          sync.RWMutex{},
        //blockRegistry:      make(map[uint64]IBlock),
        //blockLock:          sync.RWMutex{},
    }
}

func NewMeasurementModuleFromRegistry(key string) ISimulationMeasurementModule {
    if factory, ok := measurementModuleRegistry[key]; ok {
        return factory()
    }

    return nil
}

// ==== methods ====



package main

import (
    "fmt"
    "blockchainlab/simulator"
    "blockchainlab/simulator/utils"
)

func main() {
    logger := utils.GetSimulationLogger("main")
    logger.Info("starting main file")

    // create simulation form config file (section 'setup')
    sim := simulator.NewSimulationFromConfig()

    // run
    err := sim.Run()
    if err != nil {
        fmt.Println(err)
    }
    
    logger.Info("finishing main file")
}


package simulator

import (
    "blockchainlab/simulator/core"
    "blockchainlab/simulator/utils"
    _ "blockchainlab/simulator/layers/global_network"
    _ "blockchainlab/simulator/layers/node"
    _ "blockchainlab/simulator/layers/node/application"
    _ "blockchainlab/simulator/layers/node/node_network"
    _ "blockchainlab/simulator/layers/node/behavior"
    // TODO _ "blockchainlab/simulator/layers/node/consensus"
    // TODO _ "blockchainlab/simulator/layers/node/ledger"
    "fmt"
)

const (
    CONFIG_SETUP_TAG                                = "setup"
)

func init() {
    // set default layers ("setup" section)
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".end_condition",[]string{"time","600.0"})
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".global_network","default_global_network")
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".global_state","default_global_state")
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".node_list",[]string{"default_node"})
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".node_count_list",[]int{2})
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".node_network_list",[]string{"default_node_network"})
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".node_behavior_list",[]string{"default_behavior"})
    
    /* TODO config for other layers
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".node_ledger_list",[]string{"default_ledger"})
    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".node_consensus_list",[]string{"default_consensus"})
    */

    utils.ConfigSetDefault(CONFIG_SETUP_TAG + ".node_applications_list",[][]string{[]string{}})
}

// create a simulation from config file, section "setup"
func NewSimulationFromConfig() core.ISimulation {
    sim := core.NewSimulation()
    config := utils.GetSimulationConfig()

    // end condition
    endConf := config.GetStringSlice(CONFIG_SETUP_TAG + ".end_condition")
    if len(endConf) != 2 {
        panic("end_condition should be int the format: [\"name\",\"parameter\"]")
    }
    endCondition := core.NewEndConditionFromRegistry(endConf[0],endConf[1])

    // global network
    gnetConf := config.GetString(CONFIG_SETUP_TAG + ".global_network")
    gnet := core.NewGlobalNetworkFromRegistry(gnetConf)

    // global state
    gstateConf := config.GetString(CONFIG_SETUP_TAG + ".global_state")
    gstate := core.NewGlobalStateFromRegistry(gstateConf)

    // node implementations
    nodeConf := config.GetStringSlice(CONFIG_SETUP_TAG + ".node_list")
    if len(nodeConf) == 0 {
        panic("cannot create simulation: no node implementation given")
    }
   
    // node counts
    nodeCounts := config.GetIntSlice(CONFIG_SETUP_TAG + ".node_count_list")
    if len(nodeCounts) != len(nodeConf) {
        panic("cannot create simulation: node_count_list must have the same length of node_list")
    }

    // node network 
    nnetConf := config.GetStringSlice(CONFIG_SETUP_TAG + ".node_network_list")
    if len(nnetConf) != len(nodeConf) {
        panic("cannot create simulation: node_network_list must have the same length of node_list")
    }

    // behavior
    behaviorConf := config.GetStringSlice(CONFIG_SETUP_TAG + ".node_behavior_list")
    if len(behaviorConf) != len(nodeConf) {
        panic("cannot create simulation: node_behavior_list must have the same length of node_list")
    }

    /* TODO other layers
    // ledger
    ledgerConf := config.GetStringSlice(CONFIG_SETUP_TAG + ".node_ledger_list")
    if len(ledgerConf) != len(nodeConf) {
        panic("cannot create simulation: node_ledger_list must have the same length of node_list")
    }

    // consensus protocol 
    consensusConf := config.GetStringSlice(CONFIG_SETUP_TAG + ".node_consensus_list")
    if len(consensusConf) != len(nodeConf) {
        panic("cannot create simulation: node_consensus_list must have the same length of node_list")
    }
    */

    // applications (optional)
    applicationsConf := config.GetSliceStringSlice(CONFIG_SETUP_TAG + ".node_applications_list")

    // set up simulation
    sim.SetEndCondition(endCondition).SetGlobalNetwork(gnet).SetGlobalState(gstate)

    // create and add nodes
    for idx, _ := range nodeConf {
        count := nodeCounts[idx]
        for i := 0; i < count; i++ {
            nodeInstance := core.NewNodeFromRegistry(nodeConf[idx])
            if nodeInstance == nil {
                panic(fmt.Sprintf("cannot create simulation: no factory registered for %v",nodeConf[idx]))
            }

            nnetInstance := core.NewNodeNetworkFromRegistry(nnetConf[idx])
            if nnetInstance == nil {
                panic(fmt.Sprintf("cannot create simulation: no factory registered for %v",nnetConf[idx]))
            }
            
            behaviorInstance := core.NewNodeBehaviorFromRegistry(behaviorConf[idx])
            if behaviorInstance == nil {
                panic(fmt.Sprintf("cannot create simulation: no factory registered for %v",behaviorConf[idx]))
            }

            /* TODO instantiate other layers
            ledgerInstance := core.NewLedgerFromRegistry(ledgerConf[idx])
            if ledgerInstance == nil {
                panic(fmt.Sprintf("cannot create simulation: no factory registered for %v",ledgerConf[idx]))
            }

            consensusInstance := core.NewConsensusProtocolFromRegistry(consensusConf[idx])
            if consensusInstance == nil {
                panic(fmt.Sprintf("cannot create simulation: no factory registered for %v",consensusConf[idx]))
            }
            */

            nodeInstance.SetNodeNetwork(nnetInstance).
                SetBehavior(behaviorInstance)//.
                // TODO SetLedger(ledgerInstance).
                //SetConsensusProtocol(consensusInstance).

            if len(applicationsConf) > idx {
                appList := applicationsConf[idx]
                for _, appName := range appList {
                    app := core.NewApplicationFromRegistry(appName)
                    if app != nil {
                        nodeInstance.AddApplication(app)
                    } else {
                        panic(fmt.Sprintf("cannot create simulation: no factory registered for %v",appName))
                    }
                }
            }

            sim.AddNode(nodeInstance)
        }
    }

    return sim 
}


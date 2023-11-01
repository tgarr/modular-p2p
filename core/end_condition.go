package core

import (
    "strconv"
)

// ==== interfaces ====

// interface for a generic end condition checker
type IEndCondition interface {
    Check(sim ISimulation) bool
}

// ==== concrete structures  ====

/*
    End condition based on simulation time.

    Implements: IEndCondition
*/
type TimeEndCondition struct {
    endTime float64
}

// ==== factories ====

var endConditionRegistry map[string]func(arg string) IEndCondition = make(map[string]func(arg string) IEndCondition)

func init(){
    // register end condition "time"
    RegisterEndCondition("time",NewTimeEndConditionFromConfig)
}

func RegisterEndCondition(key string, factory func(arg string) IEndCondition) {
    if _, ok := endConditionRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }

    endConditionRegistry[key] = factory
}

func NewEndConditionFromRegistry(key string, arg string) IEndCondition {
    if factory, ok := endConditionRegistry[key]; ok {
        return factory(arg)
    }

    return nil
}

// factory for TimeEndCondition
func NewTimeEndCondition(endTime float64) IEndCondition {
    return &TimeEndCondition{
        endTime:    endTime,
    }
}

func NewTimeEndConditionFromConfig(arg string) IEndCondition {
    endTime, err := strconv.ParseFloat(arg,64)
    if err != nil {
        panic(err)
    }

    return &TimeEndCondition{
        endTime:    endTime,
    }
}

// ==== methods ====

// check if end condition is met
func (end *TimeEndCondition) Check(sim ISimulation) bool {
    return sim.GetTime() >= end.endTime
}


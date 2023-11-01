package utils

import (
    "math"
    "math/rand"
)

// ==== interfaces ====

type ISimulationSampler interface {
    Sample() float64
    GetDistName() string
}

// ==== concrete structures ====

/*
    Simple sampler that supports exponential, uniform, normal, and zipf
    distributions.

    Implements: ISimulationSampler
*/
type DefaultSampler struct {
    samplingFunc func() float64
    distName string
}

// ==== factories ====

/*
    Creates a new sampler based on the given distribution name, configuration,
    and random number generator.  Each distribution requires different
    parameters:
        exponential:        [average,min,max]
        uniform:            [min,max]
        normal:             [average,stddev,min,max]
        zipf:               # TODO zipf paramaters

        Note: all values are in seconds. Negative values for min or max indicate
        no limit.
*/
func NewSampler(distName string, distConfig []float64, rng *rand.Rand) ISimulationSampler {
    switch distName {
    case "exponential":
        var avg float64
        var min float64
        var max float64

        if len(distConfig) == 3 {
            avg = distConfig[0]
            min = distConfig[1]
            max = distConfig[2]
        } else {
            panic("distribution exponential requires three parameters: [average,min,max]")
        }

        return NewExponentialSampler(avg,min,max,rng)
    case "uniform":
        var min float64
        var max float64

        if len(distConfig) == 2 {
            min = distConfig[0]
            max = distConfig[1]
        } else {
            panic("distribution uniform requires two parameters: [min,max]")
        }

        return NewUniformSampler(min,max,rng)
    case "normal":
        var avg float64
        var std float64
        var min float64
        var max float64

        if len(distConfig) == 4 {
            avg = distConfig[0]
            std = distConfig[1]
            min = distConfig[2]
            max = distConfig[3]
        } else {
            panic("distribution normal requires four parameters: [average,stddev,min,max]")
        }

        return NewNormalSampler(avg,std,min,max,rng)
    case "zipf":
        // TODO build zipf lambda
        panic("distribution zipf not supported yet")
    default:
        panic("distribution " + distName + " not supported")
    }

    return nil
}

func NewNormalSampler(avg,std,min,max float64,rng *rand.Rand) ISimulationSampler {
    if max < 0 {
        max = +math.MaxFloat64
    }

    if min < 0 {
        min = 0
    }

    samplingFunc := func() float64 {
        sample := rng.NormFloat64() * std + avg
        sample = math.Min(sample,max)
        sample = math.Max(sample,min)
        return sample
    }

    return &DefaultSampler{
        samplingFunc:   samplingFunc,
        distName:       "normal",
    }
}

func NewUniformSampler(min,max float64,rng *rand.Rand) ISimulationSampler {
    if max < 0 {
        max = +math.MaxFloat64
    }

    if min < 0 {
        min = 0
    }

    samplingFunc := func() float64 {
        sample := min + rng.Float64() * (max-min)
        return sample
    }

    return &DefaultSampler{
        samplingFunc:   samplingFunc,
        distName:       "uniform",
    }
}

func NewExponentialSampler(avg,min,max float64,rng *rand.Rand) ISimulationSampler {
    if max < 0 {
        max = +math.MaxFloat64
    }

    if min < 0 {
        min = 0
    }

    samplingFunc := func() float64 {
        sample := rng.ExpFloat64() * avg
        sample = math.Min(sample,max)
        sample = math.Max(sample,min)
        return sample
    }

    return &DefaultSampler{
        samplingFunc:   samplingFunc,
        distName:       "exponential",
    }
}

// ==== methods ====

func (sampler *DefaultSampler) Sample() float64 {
    return sampler.samplingFunc()
}

// ==== getters ====

func (sampler *DefaultSampler) GetDistName() string {
    return sampler.distName
}


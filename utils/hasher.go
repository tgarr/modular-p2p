package utils

import (
    "github.com/zeebo/xxh3"
)

// ==== interfaces ====

type IHashable interface {
    GetHash() uint64
}

// ==== concrete structures ====

// the hasher used in the simulator is a wrapper on top of xxh3
type SimulationHasher struct {
    hasher *xxh3.Hasher
}

// ==== factories ====

func NewHasher() *SimulationHasher {
    return &SimulationHasher{
        hasher:     xxh3.New(),
    }
}

// ==== methods ====

func (hasher *SimulationHasher) WriteBytes(buffer []byte) {
    hasher.hasher.Write(buffer)
}

func (hasher *SimulationHasher) WriteString(str string) {
    hasher.hasher.WriteString(str)
}

func (hasher *SimulationHasher) Hash() uint64 {
    return hasher.hasher.Sum64()
}

func HashBytes(buffer []byte) uint64 {
    return xxh3.Hash(buffer)
}

func HashString(str string) uint64 {
    return xxh3.HashString(str)
}


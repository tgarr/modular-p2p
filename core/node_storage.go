package core

import (
    "blockchainlab/simulator/utils"
)

// node state fields
const (
    NODE_STATE_UTXO                                 = "utxo"
    NODE_STATE_ACCOUNTS                             = "accounts"
)

// block types
const (
    BLOCK_STANDARD                                  = 0
    BLOCK_PRIMARY                                   = 1
    BLOCK_SECONDARY                                 = 2
)

// transaction types
const (
    TX_STANDARD                                     = 0
    TX_UTXO_SPEND                                   = 1
    TX_CREATE_CONTRACT                              = 2
    TX_INVOKE_CONTRACT                              = 3
    TX_CREATE_ASSET                                 = 4
    TX_TRANSFER_ASSET                               = 5
)

// block reference types
const (
    BREF_STANDARD                                   = 0
    BREF_UNCLE                                      = 1
    BREF_PRIMARY                                    = 2
    BREF_SECONDARY                                  = 3
)

// metadata fields
const (
    METADATA_HEIGHT                                 = "height"
    METADATA_WEIGHT                                 = "weight"
    METADATA_CONFIDENCE                             = "confidence"
)

type INodeStorage interface {
    ISimulationComponent

    GetState() map[string]interface{}
    GetLedger() ILedger

    VerifyBlock(block IBlock) bool
    VerifyTransaction(tx ITransaction) bool
}

type ILedger interface {
    GetMetadata() map[uint64]map[string]interface{}             // [hash][field] -> data
    GetBlockMetadata(blockID uint64) map[string]interface{}     // [field] -> data
    GetTxMetadata(txID uint64) map[string]interface{}           // [field] -> data
    
    GetBlocks() map[uint64]IBlock
    GetBlock(uint64) IBlock
    GetTransaction(txID uint64) ITransaction
    
    GetSize() uint64                                            // size of ledger (e.g. max height)
    
    SetGenesisBlock(block IBlock)
    AddBlock(block IBlock)
    RemoveBlock(block IBlock)
}

type IBlock interface {
    ILedgerElement

    GetReferences() map[uint16][]uint64
    GetTransactions() map[uint16][]ITransaction
}

type ITransaction interface {
    ILedgerElement
}

type ILedgerElement interface {
    utils.IHashable
    
    GetType() uint16
    GetTime() float64
    GetCreator() uint32
    GetSize() uint64
    Verify() bool
}

// ==== factories ====

var nodeStorageRegistry map[string]func() INodeStorage = make(map[string]func() INodeStorage)

func RegisterNodeStorage(key string, factory func() INodeStorage) {
    if _, ok := nodeStorageRegistry[key]; ok {
        panic("factory for " + key + " already registered!")
    }

    nodeStorageRegistry[key] = factory
}

func NewNodeStorageFromRegistry(key string) INodeStorage {
    if factory, ok := nodeStorageRegistry[key]; ok {
        return factory()
    }

    return nil 
}


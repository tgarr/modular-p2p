package core

import (
    "reflect"
)

// message delivery types
const (
    MESSAGE_DELIVERY_TYPE_P2P_NODES                         = 0 // p2p message to the given nodes
    MESSAGE_DELIVERY_TYPE_P2P_NODE_TYPES                    = 1 // p2p message to nodes of the given types
    MESSAGE_DELIVERY_TYPE_P2P_NODE_TYPES_EXCEPT             = 2 // p2p message to nodes of types different than the given types
    MESSAGE_DELIVERY_TYPE_BROADCAST_NODES                   = 3 // broadcast message to the given nodes
    MESSAGE_DELIVERY_TYPE_BROADCAST_NODE_TYPES              = 4 // broadcast message to nodes of the given types
    MESSAGE_DELIVERY_TYPE_BROADCAST_NODE_TYPES_EXCEPT       = 5 // broadcast message to nodes of types different than the given types
)

// ==== interfaces ====

// A network message sent from a node to one or more nodes
type IMessage interface {
    GetData() interface{}                   // message data (receiver must cast to correct type)
    GetSize() uint64                        // size of the message
    GetSender() uint32                      // node id of the sender
    GetDelivery() IMessageDelivery          // delivery mode
    GetTag() int32                          // custom tag
    GetTime() float64                       // time message was sent

    SetTag(tag int32) IMessage              // set custom tag
    SetTime(time float64) IMessage          // set message time
    SetSize(size uint64) IMessage           // set message size
}

// indicates the type of delivery the message requires 
type IMessageDelivery interface {
    GetDeliveryType() uint16                // delivery type
    GetDeliveryTargets() []uint32           // delivery targets

    SetDeliveryType(tp uint16) IMessageDelivery
    SetDeliveryTargets(targets []uint32) IMessageDelivery
}

// ==== concrete structures ====

/*
    A generic message implementation that can hold any type of data.

    Implements: IMessage interface.
*/
type DefaultMessage struct {
    data interface{}
    size uint64
    sender uint32
    delivery IMessageDelivery
    tag int32
    time float64
}

/*
    A simple delivery implementation that cover p2p and broadcast modes. Targets can be node IDs or node types. If target is nil, send to all nodes using the selected mode.

    Implements: IMessageDelivery
*/
type DefaultDelivery struct {
    tp uint16
    targets []uint32
}

// ==== factories ====

// broadcast message to all nodes that receive global broadcasts
func NewBroadcastMessage(data interface{},sender uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_BROADCAST_NODES,nil)
    return NewMessage(data,sender,delivery)
}

// broadcast message to a subset of nodes
func NewBroadcastMessageNodes(data interface{},sender uint32,nodes []uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_BROADCAST_NODES,nodes)
    return NewMessage(data,sender,delivery)
}

// broadcast message to a subset of node types
func NewBroadcastMessageNodeTypes(data interface{},sender uint32,nodes []uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_BROADCAST_NODE_TYPES,nodes)
    return NewMessage(data,sender,delivery)
}

// broadcast message to all node types except the types given
func NewBroadcastMessageNodeTypesExcept(data interface{},sender uint32,nodes []uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_BROADCAST_NODE_TYPES_EXCEPT,nodes)
    return NewMessage(data,sender,delivery)
}

// p2p message to one specific node
func NewP2PMessage(data interface{},sender uint32,destination uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_P2P_NODES,[]uint32{destination})
    return NewMessage(data,sender,delivery)
}

// p2p message to a subset of nodes
func NewP2PMessageNodes(data interface{},sender uint32,nodes []uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_P2P_NODES,nodes)
    return NewMessage(data,sender,delivery)
}

// p2p message to all nodes
func NewP2PMessageAll(data interface{},sender uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_P2P_NODES,nil)
    return NewMessage(data,sender,delivery)
}

// p2p message to a subset of node types
func NewP2PMessageNodeTypes(data interface{},sender uint32,nodes []uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_P2P_NODE_TYPES,nodes)
    return NewMessage(data,sender,delivery)
}

// p2p message to all node types except the types given
func NewP2PMessageNodeTypesExcept(data interface{},sender uint32,nodes []uint32) IMessage {
    delivery := NewMessageDelivery(MESSAGE_DELIVERY_TYPE_P2P_NODE_TYPES_EXCEPT,nodes)
    return NewMessage(data,sender,delivery)
}

// complete factory for DefaultMessage
func NewMessage(data interface{},sender uint32,delivery IMessageDelivery) IMessage {
    return &DefaultMessage{
        data:           data,
        size:           uint64(reflect.TypeOf(data).Size()),
        sender:         sender,
        delivery:       delivery,
        tag:            0, // optional
        time:           0, // optional
    }
}

// factory for message delivery
func NewMessageDelivery(tp uint16,targets []uint32) IMessageDelivery {
    return &DefaultDelivery{
        tp:         tp,
        targets:    targets,
    }
}

// ==== getters ====

func (msg *DefaultMessage) GetData() interface{} {
    return msg.data
}

func (msg *DefaultMessage) GetSize() uint64 {
    return msg.size
}

func (msg *DefaultMessage) GetSender() uint32 {
    return msg.sender
}

func (msg *DefaultMessage) GetDelivery() IMessageDelivery {
    return msg.delivery
}

func (msg *DefaultMessage) GetTag() int32 {
    return msg.tag
}

func (msg *DefaultMessage) GetTime() float64 {
    return msg.time
}

func (del *DefaultDelivery) GetDeliveryType() uint16 {
    return del.tp
}

func (del *DefaultDelivery) GetDeliveryTargets() []uint32 {
    return del.targets
}

// ==== setters ====

func (msg *DefaultMessage) SetTag(tag int32) IMessage {
    msg.tag = tag
    return msg
}

func (msg *DefaultMessage) SetTime(time float64) IMessage {
    msg.time = time
    return msg
}

func (msg *DefaultMessage) SetSize(size uint64) IMessage {
    msg.size = size
    return msg
}

func (del *DefaultDelivery) SetDeliveryType(tp uint16) IMessageDelivery {
    del.tp = tp
    return del
}

func (del *DefaultDelivery) SetDeliveryTargets(targets []uint32) IMessageDelivery {
    del.targets = targets
    return del
}


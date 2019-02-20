package socketio
import "github.com/equalll/mydebug"

import (
	"encoding/json"
)

const (
	PACKET_DISCONNECT = iota
	PACKET_CONNECT
	PACKET_HEARTBEAT
	PACKET_MESSAGE
	PACKET_JSONMESSAGE
	PACKET_EVENT
	PACKET_ACK
	PACKET_ERROR
	PACKET_NOOP
)

type MessageType uint8

type Packet interface {
	Id() int
	Type() MessageType
	EndPoint() string
	Ack() bool
}

type packetCommon struct {
	id       int
	endPoint string
	ack      bool
}

func (p *packetCommon) Id() int {mydebug.INFO()
	return p.id
}

func (p *packetCommon) EndPoint() string {mydebug.INFO()
	return p.endPoint
}

func (p *packetCommon) Ack() bool {mydebug.INFO()
	return p.ack
}

type disconnectPacket struct {
	packetCommon
}

func (*disconnectPacket) Type() MessageType {mydebug.INFO()
	return PACKET_DISCONNECT
}

type connectPacket struct {
	packetCommon
	query string
}

func (*connectPacket) Type() MessageType {mydebug.INFO()
	return PACKET_CONNECT
}

type heartbeatPacket struct {
	packetCommon
}

func (*heartbeatPacket) Type() MessageType {mydebug.INFO()
	return PACKET_HEARTBEAT
}

type messageMix interface {
	Packet
	Data() []byte
}

type messagePacket struct {
	packetCommon
	data []byte
}

func (*messagePacket) Type() MessageType {mydebug.INFO()
	return PACKET_MESSAGE
}

func (p *messagePacket) Data() []byte {mydebug.INFO()
	return p.data
}

type jsonPacket struct {
	packetCommon
	data []byte
}

func (*jsonPacket) Type() MessageType {mydebug.INFO()
	return PACKET_JSONMESSAGE
}

func (p *jsonPacket) Data() []byte {mydebug.INFO()
	return p.data
}

type eventPacket struct {
	packetCommon
	name string
	args json.RawMessage
}

func (*eventPacket) Type() MessageType {mydebug.INFO()
	return PACKET_EVENT
}

type ackPacket struct {
	packetCommon
	ackId int
	args  json.RawMessage
}

func (*ackPacket) Type() MessageType {mydebug.INFO()
	return PACKET_ACK
}

type errorPacket struct {
	packetCommon
	reason string
	advice string
}

func (*errorPacket) Type() MessageType {mydebug.INFO()
	return PACKET_ERROR
}

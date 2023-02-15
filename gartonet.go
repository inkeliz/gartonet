package gartonet

import (
	"context"
	"errors"
	"net"
	"unsafe"
)

// Client is an Art-Net client.
type Client struct {
	ctx    context.Context
	cancel func()

	seq [256]uint8

	udp *net.UDPConn
}

// NewClient creates a new Art-Net client.
func NewClient(addr *net.UDPAddr) (*Client, error) {
	if addr.Port == 0 {
		addr.Port = 6454
	}

	udp, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{ctx: ctx, cancel: cancel, udp: udp}

	return c, nil
}

// NewClientString creates a new Art-Net client, from a string IP address,
// that assumes port 6454.
func NewClientString(ip string) (*Client, error) {
	return NewClient(&net.UDPAddr{IP: net.ParseIP(ip), Port: 6454})
}

// Send sends the Art-Net packet.
func (c *Client) Send(packet *Packet) error {
	if c == nil || packet == nil {
		return errors.New("invalid packet or client")
	}
	c.seq[packet.Header[HeaderOffsetUniverse]]++
	packet.Header[HeaderOffsetSeq] = c.seq[packet.Header[HeaderOffsetUniverse]]
	_, _, err := c.udp.WriteMsgUDP(packet.Bytes(), nil, nil)
	return err
}

// Close closes the client.
func (c *Client) Close() error {
	c.cancel()
	return c.udp.Close()
}

var (
	_Header = [...]byte{
		'A', 'r', 't', '-',
		'N', 'e', 't', 0,
		0x00, 0x50, 0, 0x0e,
		0, 0, 0, 0,
		512 >> 8, 512 & 0xFF,
	}
)

const (
	// HeaderOffsetOpCode is the offset of the OpCode in the header.
	HeaderOffsetOpCode = 8
	// HeaderOffsetSeq is the offset of the sequence number in the header.
	HeaderOffsetSeq = 12
	// HeaderOffsetNet is the offset of the net and subnet in the header.
	HeaderOffsetNet = 15
	// HeaderOffsetUniverse is the offset of the universe in the header.
	HeaderOffsetUniverse = 14
	// HeaderOffsetLength is the offset of the length in the header.
	HeaderOffsetLength = 16
)

// Packet is the Art-Net packet.
// Usually you want to use NewPacket to create a new packet, and
// re-use it for multiple updates.
//
// The packet is not thread-safe, you should not update the packet
// while it's being sent.
//
// Usually you don't change Header, but you can change the DMX data.
type Packet struct {
	// Header is the Art-Net header. Usually, you don't need to change
	// this, if you create a new packet using NewPacket.
	Header [18]uint8
	// DMX is the DMX data. The length of the DMX data is always 512.
	DMX [512]uint8
}

// NewPacket creates a new Art-Net packet. Assuming that the packet is
// going to be sent to the specified universe. The net and sub are
// usually 0. The universe is the universe number upto 255.
//
// It's highly recommended to re-use the packet for multiple sends,
// the Client.Send updates the sequence number.
func NewPacket(net uint8, universe uint8) *Packet {
	p := Packet{}
	copy(p.Header[:], _Header[:])

	p.Header[HeaderOffsetUniverse] = universe
	p.Header[HeaderOffsetNet] = net

	return &p
}

// Bytes returns the bytes of the packet, including the header and the DMX data.
//
// That is not a copy, so you should not change the packet while it's being sent.
func (p *Packet) Bytes() []byte {
	return ((*[512 + 18]uint8)(unsafe.Pointer(&p.Header[0])))[:]
}

package rcon

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	quake "github.com/icedream/go-q3net"
)

type Client struct {
	address    *net.UDPAddr
	addressStr string
	password   string

	socket       *net.UDPConn
	socketBuffer []byte
}

func NewRconClient() *Client {
	socketBuffer := make([]byte, 64*1024)
	return &Client{
		socketBuffer: socketBuffer,
	}
}

func (c *Client) SetSocketAddr(addr string) (err error) {
	newAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}

	c.address, c.addressStr = newAddr, addr

	return
}

func (c *Client) SetPassword(pw string) {
	c.password = pw
}

func (c *Client) Address() *net.UDPAddr {
	return c.address
}

func (c *Client) AddressString() string {
	return c.addressStr
}

func (c *Client) Password() string {
	return c.password
}

func (c *Client) InitSocket() (err error) {
	c.socket, err = net.ListenUDP("udp", nil)
	if err != nil {
		return
	}

	return
}

func (c *Client) udpReceiveAndUnmarshal() (msg *quake.Message, err error) {
	length, _, err := c.socket.ReadFromUDP(c.socketBuffer)
	if err != nil {
		return
	}

	msg, err = quake.UnmarshalMessage(c.socketBuffer[0:length])
	if err != nil {
		return
	}
	return
}

func (c *Client) Receive() (msg *quake.Message, err error) {
	msg, err = c.udpReceiveAndUnmarshal()
	if err != nil {
		return
	}
	if !strings.EqualFold(msg.Name, "print") {
		err = errors.New("rcon: Unexpected response from server: " + msg.Name)
	}
	return
}

func (c *Client) Send(input string) (err error) {
	buf := new(bytes.Buffer)
	msg := &quake.Message{
		Header: quake.OOBHeader,
		Name:   "rcon",
		Data:   []byte(fmt.Sprintf("%s %s", c.password, input)),
	}
	if err = msg.Marshal(buf); err != nil {
		return
	}
	if _, err = c.socket.WriteToUDP(buf.Bytes(), c.address); err != nil {
		return
	}
	return
}

func (c *Client) Release() {
	if c.socket != nil {
		c.socket.Close()
		c.socket = nil
	}
}

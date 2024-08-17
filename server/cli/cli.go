package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

const (
	connectLimit = time.Second * 5
)

type clientChannel chan string

type ClientModel struct {
	Channel      clientChannel
	Conn         net.Conn
	IsTerminated atomic.Bool
	Name         string
	Reader       io.Reader
	ReadBuffer  []byte
	Timer		*time.Timer
}

func NewClientModel(conn net.Conn, name string) *ClientModel {
	return &ClientModel{
		Channel: make(clientChannel,10),
		Conn:    conn,
		Name:    name,
		IsTerminated: atomic.Bool{},
		Reader:  bufio.NewReader(conn),
		ReadBuffer: make([]byte, 1024),
		Timer: time.NewTimer(connectLimit),
	}
}

func (c *ClientModel) DetectTimeout(leaving chan *ClientModel) {
	for range c.Timer.C {
		leaving<-c
	}
}

func (c *ClientModel) ReadMsg() (string,error){
		if c.IsTerminated.Load() {
			return "",errors.New("client is terminated")
		}
		n,err:= c.Reader.Read(c.ReadBuffer)
		if err!=nil {
			return "",err
		}
		c.Timer.Reset(connectLimit)
		return string(c.ReadBuffer[:n]),nil

}

func (c *ClientModel) SendMsg(msg string) {
	c.Channel <- msg
}

func (c *ClientModel) WriteMsg(){
	for msg:= range c.Channel {
		fmt.Fprintln(c.Conn,msg)
	}
}

func (c *ClientModel) Close() {
	if c.IsTerminated.Load() {
		return
	}
	c.IsTerminated.Store(true)
	c.Conn.Close()
}
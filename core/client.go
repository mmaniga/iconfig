package core

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// Client :
type Client struct {
	clientID      string
	socket        *websocket.Conn
	configSend    chan *ConfigInfo // This is for future in case we have a scenario if config or other updated from client
	configReceive chan *ConfigInfo
	send          chan *ConfigInfo
	crequest      chan *ConfigRequest
	st            *store
	connected     bool
}

func (c *Client) read() {
	var msg = ConfigInfo{}
	var roMsg = ConfigRequest{}
	defer func() {
		c.st.unregister <- c
		c.socket.Close()
	}()
	for c.connected {
		if _, byt, eror := c.socket.ReadMessage(); eror == nil {
			if err := json.Unmarshal(byt, &msg); err == nil && msg.MType == ConfigSend {
				fmt.Println("HAndling ConfigSend")
				c.configReceive <- &msg
			} else if err != nil {
				fmt.Println("Read Error = %+v\n", err.Error())
				c.connected = false
				break
			} else {
				fmt.Println("Entering alternate path configRequest")
				if err1 := json.Unmarshal(byt, &roMsg); err1 == nil {
					fmt.Println("Entering for ConfigRequest client =  ", roMsg.ClientID, roMsg.Key)
					c.crequest <- &roMsg
				} else if err1 != nil {
					fmt.Printf("Read Error = %+v\n", err.Error())
					c.connected = false
					break
				}
			}
		} else {
			fmt.Println(eror)
			c.connected = false
			break
		}
	}
}

func (c *Client) write() {
	for c.connected {
		for msg := range c.configSend {
			fmt.Println("Client Write ", msg)
			if err := c.socket.WriteJSON(msg); err != nil {
				c.connected = false
				break
			}
		}
	}
	c.socket.Close()
}

func (c *Client) writerfid() {
	for c.connected {
		for msg := range c.send {
			fmt.Println("What happens in writefid ", msg)
			if err := c.socket.WriteJSON(msg); err != nil {
				c.connected = false
				break
			}
		}
	}
	c.socket.Close()
}

func (c *Client) writeConfig() {
	for c.connected {
		for msg := range c.configReceive {
			fmt.Println("What happens in writeConfig ", msg)
			if err := c.socket.WriteJSON(msg); err != nil {
				c.connected = false
				break
			}
		}
	}
	c.socket.Close()
}

func (c *Client) handleCRequest() {
	for c.connected {
		for msg := range c.crequest {
			fmt.Println("in handleCRequest ", msg)
			m := ConfigInfo{Key: msg.Key, Value: c.st.cmap[msg.Key]}
			c.configReceive <- &m
		}
	}
	c.socket.Close()
}
func (c *Client) handleMessage() {
	for {
		select {
		case msg := <-c.configReceive:
			fmt.Println("Getting into configReceive ", msg)
			c.st.config <- msg
		}
	}
}

//NewClient :
func NewClient(socket *websocket.Conn, h *store, cID string) *Client {
	return &Client{
		socket:        socket,
		configSend:    make(chan *ConfigInfo),
		configReceive: make(chan *ConfigInfo),
		send:          make(chan *ConfigInfo),
		crequest:      make(chan *ConfigRequest),
		st:            h,
		clientID:      cID,
		connected:     true,
	}
}

package core

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type store struct {
	crequest   chan *ConfigRequest // Request for configuration from client
	config     chan *ConfigInfo    // A channel of ConfigInfo
	register   chan *Client        // To do it in non blocking manner we do this
	unregister chan *Client        // To do it in non blocking manner we do this
	clients    map[*Client]bool    // Map of all the clients in the network
	cmap       map[string]string   // This is the actual configuration map

}

func newStore() *store {
	return &store{
		crequest:   make(chan *ConfigRequest),
		config:     make(chan *ConfigInfo),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		cmap:       make(map[string]string),
	}
}

//SendDataToWS :
func SendDataToWS(con *websocket.Conn, str []byte) error {
	return con.WriteMessage(websocket.TextMessage, str)
}

func (r *store) run() {
	fmt.Println("Starting the Config Server....")
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			delete(r.clients, client)
			close(client.configSend)
		case msg := <-r.config: // This channel is used only for all client communication and update central store
			go func() {
				r.cmap[msg.Key] = msg.Value // No need for lock, its in channel
				for c := range r.clients {  // Iterate all clients and send config message to them
					c.send <- msg
				}
			}()
		case msg := <-r.crequest:
			go func() {
				c := GetClient(msg.ClientID)
				msgs := ConfigInfo{Key: msg.Key, Value: r.cmap[msg.Key]}
				fmt.Println("Sending info from r.crequest for key value")
				c.send <- &msgs
			}()

		}
	}
}

const (
	socketBufferSize = 12288
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// format: /join/{client-id}

func (r *store) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	args := strings.Split(req.URL.Path, "/")
	clientID := args[2]
	fmt.Println("Client ID is ", clientID)

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Println("ServeHTTP:", err)
		return
	}
	client := NewClient(socket, r, clientID)
	r.register <- client
	AddClient(clientID, client)

	defer func() { r.unregister <- client }()
	go client.write()

	go client.writerfid()
	go client.handleMessage()
	go client.writeConfig()
	go client.handleCRequest()
	client.read()
}

func (r *store) Get(k string) string {
	return r.cmap[k]
}

var clients = struct {
	sync.RWMutex
	m map[string]*Client
}{m: make(map[string]*Client)}

// AddClient :
func AddClient(clientID string, c *Client) {
	clients.Lock()
	clients.m[clientID] = c
	clients.Unlock()
}

// GetClient :
func GetClient(clientID string) (c *Client) {
	clients.RLock()
	c = clients.m[clientID]
	clients.RUnlock()
	return c
}

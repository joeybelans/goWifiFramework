package webSocketHandler

import "github.com/gorilla/websocket"

// Upgrade HTTP connection to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Page/Data mappings
type dataMapKey struct {
	component string
	message   string
	page      string
}

// Package variables
var (
	wsconn   *websocket.Conn
	curPage  string
	dataMap  map[dataMapKey]bool
	packages map[string]string
)

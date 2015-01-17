package kismet

import (
	"database/sql"
	"net"

	"github.com/gorilla/websocket"
)

type kismetServer struct {
	version   string
	starttime string
	name      string
	dumpfiles []string
	uid       int
	protocols []string
}

type networkInterface struct {
	pname       string
	lname       string
	hwaddr      string
	channel     string
	hop         string
	velocity    string
	channellist string
	bssids      map[string]string
}

type kismetClient struct {
	host       string
	port       int
	db         *sql.DB
	debug      bool
	conn       net.Conn
	server     kismetServer
	index      int
	pipeline   map[int]string
	interfaces map[string]networkInterface
}

type accessPoint struct {
	//uuid string
	apType       string
	ssid         string
	manuf        string
	channel      int
	firsttime    string
	lasttime     string
	atype        int
	rangeip      string
	netmaskip    string
	gatewayip    string
	signalDBM    int
	minSignalDBM int
	maxSignalDBM int
	numPackets   int
}

type network struct {
	//type???
	cloaked    int
	firsttime  string
	lasttime   string
	maxrate    int
	encryption string
	bssids     map[string]string
}

type client struct {
	bssid        string
	firsttime    string
	lasttime     string
	signalDBM    int
	minSignalDBM int
	maxSignalDBM int
	numPackets   int
}

var (
	kismet       kismetClient
	networks     map[string]network
	accessPoints map[string]accessPoint
	clients      map[string]client
	tstamp       string
	packets      int
	packetRate   int
	filtered     bool
	wsconn       *websocket.Conn
	curPage      string
	ssids        []string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

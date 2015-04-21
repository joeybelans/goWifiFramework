package kismet

import (
	"database/sql"
	"net"

	"github.com/gorilla/websocket"
)

// Packet statistics
type packetStats struct {
	total    int
	crypt    int
	dropped  int
	rate     int
	filtered int
	mgmt     int
	data     int
}

// Kismet server object
type kismetServer struct {
	version   string
	starttime string
	name      string
	dumpfiles []string
	uid       int
	stats     packetStats
}

// Network interface object
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

// Kismet client object
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

// Access point object
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

// Wireless network object
type network struct {
	inscope    bool
	cloaked    int
	firsttime  string
	lasttime   string
	maxrate    int
	encryption string
	bssids     map[string]string
}

// Wireless client object
type client struct {
	bssid        string
	firsttime    string
	lasttime     string
	signalDBM    int
	minSignalDBM int
	maxSignalDBM int
	numPackets   int
}

// Package variables
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

// Upgrade HTTP connection to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Kismet capabilities
type capability struct {
	function interface{}
	fields   []string
}

var capabilities = map[string]capability{
	"KISMET":   capability{parseKISMET, []string{"version", "starttime", "servername", "dumpfiles", "uid"}},
	"TIME":     capability{parseTIME, []string{"timesec"}},
	"ACK":      capability{parseACK, []string{"cmdid", "text"}},
	"INFO":     capability{parseINFO, []string{"packets", "crypt", "dropped", "rate", "filtered", "llcpackets", "datapackets"}},
	"STATUS":   capability{parseSTATUS, []string{"text", "flags"}},
	"ALERT":    capability{parseALERT, []string{"sec", "usec", "header", "bssid", "source", "dest", "other", "channel", "text"}},
	"ERROR":    capability{parseERROR, []string{"cmdid", "text"}},
	"BSSIDSRC": capability{parseBSSIDSRC, []string{"bssid", "uuid", "lasttime"}},
	"BSSID": capability{parseBSSID,
		[]string{"bssid", "type", "manuf", "channel", "firsttime", "lasttime", "atype", "rangeip", "netmaskip", "gatewayip", "signal_dbm", "minsignal_dbm", "maxsignal_dbm", "datapackets"}},
	"SSID": capability{parseSSID,
		[]string{"mac", "checksum", "type", "ssid", "beaconinfo", "cryptset", "cloaked", "firsttime", "lasttime", "maxrate", "beaconrate", "packets", "beacons", "dot11d"}},
	"CLISRC": capability{parseCLISRC, []string{"bssid", "mac", "uuid", "lasttime", "numpackets", "signal_dbm", "minsignal_dbm", "maxsignal_dbm"}},
	"CLIENT": capability{parseCLIENT, []string{"bssid", "mac", "type", "firsttime", "lasttime", "manuf", "signal_dbm", "minsignal_dbm", "maxsignal_dbm", "datapackets"}},
	"SOURCE": capability{parseSOURCE,
		[]string{"interface", "type", "username", "channel", "uuid", "packets", "hop", "velocity", "dwell", "hop_time_sec", "hop_time_usec", "channellist", "error", "warning"}},
	"TERMINATE": capability{parseTERMINATE, []string{"text"}},
}

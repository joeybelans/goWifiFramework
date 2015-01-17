package kismet

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func Send(cmd string, params string) {
	cmdstr := cmd + " " + params
	kismet.pipeline[kismet.index] = cmdstr
	cmdstr = fmt.Sprintf("!%d ", kismet.index) + cmdstr
	if kismet.debug {
		fmt.Println("SEND: " + cmdstr)
	}
	kismet.conn.Write([]byte(cmdstr + "\n"))
	kismet.index += 1
}

func Run(host string, port int, db *sql.DB, debug bool, ssids []string) {
	// Set package variables
	kismet.host = host
	kismet.port = port
	kismet.db = db
	kismet.debug = debug
	kismet.index = 1
	kismet.pipeline = map[int]string{}
	kismet.interfaces = map[string]networkInterface{}
	networks = map[string]network{}
	accessPoints = map[string]accessPoint{}
	clients = map[string]client{}
	filtered = false
	ssids = ssids

	// Establish connection to kismet server
	var err error
	kismet.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", kismet.host, kismet.port))
	if err != nil {
		log.Fatal("Can't connect to kismet server")
	}

	// Listen for responses
	go listen()

	// Setup the client
	setupClient()
}

func ServerVersion() string {
	return kismet.server.version
}

func ServerName() string {
	return kismet.server.name
}

func ServerStart() string {
	return kismet.server.starttime
}

func Connected() bool {
	tstampInt, _ := strconv.ParseInt(tstamp, 10, 64)
	if time.Now().Unix()-tstampInt < 5 {
		return true
	}
	return false
}

func Stats() (int, int, int, int, int, int) {
	fInt := 0
	if filtered {
		fInt = 1
	}

	rCount := 0
	for network := range networks {
		rogue := true
		for _, ssid := range ssids {
			if ssid == network {
				rogue = false
				break
			}
		}
		if rogue {
			rCount++
		}
	}

	return len(networks), len(clients), rCount, packets, packetRate, fInt
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
	var err error
	wsconn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, p, err := wsconn.ReadMessage()
	if err != nil {
		curPage = ""
		return
	}
	curPage = string(p)

	for {
		_, message, err := wsconn.ReadMessage()
		if err != nil {
			break
		}
		processWSCommand(string(message))
	}
}

func GetNetworks() map[string]network {
	return networks
}

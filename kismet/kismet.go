package kismet

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joeybelans/gokismet/kdb"
)

type kismetServer struct {
	version   string
	starttime string
	name      string
	dumpfiles []string
	uid       int
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

func listen() {
	// Create mapping to parse kismet responses
	responses := map[string]interface{}{
		"KISMET":     parseKISMET,
		"PROTOCOLS":  parsePROTOCOLS,
		"CAPABILITY": parseCAPABILITY,
		"TIME":       parseTIME,
		"ACK":        parseACK,
		"INFO":       parseINFO,
		"STATUS":     parseSTATUS,
		"ALERT":      parseALERT,
		"ERROR":      parseERROR,
		"BSSIDSRC":   parseBSSIDSRC,
		"BSSID":      parseBSSID,
		"SSID":       parseSSID,
		"CLISRC":     parseCLISRC,
		"CLIENT":     parseCLIENT,
		"SOURCE":     parseSOURCE,
		"TERMINATE":  parseTERMINATE,
	}

	// Continuously read data
	scanner := bufio.NewScanner(kismet.conn)
	re := regexp.MustCompile(`([^ \001]+|\001[^\001]*\001)`)
	for scanner.Scan() {
		status := scanner.Text()

		// Determine response type
		status = strings.TrimSpace(status)
		matches := re.FindAllStringSubmatch(status, -1)
		fields := parseFields(matches)
		f, found := responses[fields[0]]
		if found {
			f.(func([]string))(fields)
		} else if kismet.debug {
			fmt.Println("UNKNOWN CMD: " + fields[0])
			fmt.Println(status)
		}
	}
	fmt.Println("QUITING")
}

func setupClient() {
	// Disable TIME
	//Send("REMOVE", "TIME")

	// Enable default capabilities
	capabilities := []string{
		"INFO networks,packets,rate,filtered,llcpackets,numsources,numerrorsources",
		"STATUS text,flags",
		"SOURCE interface,type,username,channel,uuid,packets,hop,velocity,dwell,hop_time_sec,hop_time_usec,channellist,error,warning",
		"ALERT sec,usec,header,bssid,source,dest,other,channel,text",
		"ERROR cmdid,text",
		"BSSIDSRC bssid,uuid,lasttime",
		"BSSID bssid,type,manuf,channel,firsttime,lasttime,atype,rangeip,netmaskip,gatewayip,signal_dbm,minsignal_dbm,maxsignal_dbm,datapackets",
		"SSID mac,checksum,type,ssid,beaconinfo,cryptset,cloaked,firsttime,lasttime,maxrate,beaconrate,packets,beacons,dot11d,",
		"NETTAG bssid,tag,value",
		"CLISRC bssid,mac,uuid,lasttime,numpackets,signal_dbm,minsignal_dbm,maxsignal_dbm",
		"CLIENT bssid,mac,type,firsttime,lasttime,manuf,signal_dbm,minsignal_dbm,maxsignal_dbm,datapackets",
	}

	for _, txt := range capabilities {
		Send("ENABLE", txt)
	}

	_, ok := kismet.pipeline[kismet.index]
	for ok || len(kismet.pipeline) > 0 {
		time.Sleep(1 * time.Second)
		_, ok = kismet.pipeline[kismet.index]
	}
}

func parseFields(matches [][]string) []string {
	var params []string

	cmd := strings.TrimRight(strings.TrimLeft(matches[0][0], "*"), ":")
	params = append(params, cmd)

	for i := 1; i < len(matches); i++ {
		txt := matches[i][0]
		txt = strings.Trim(txt, "\001")
		txt = strings.TrimSpace(txt)
		params = append(params, txt)
	}

	return params
}

func parseKISMET(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	kismet.server.version = fields[1]
	kismet.server.starttime = fields[2]
	kismet.server.name = fields[3]
	kismet.server.dumpfiles = strings.Split(fields[4], ",")
	kismet.server.uid, _ = strconv.Atoi(fields[5])
}

func parsePROTOCOLS(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}
}

func parseCAPABILITY(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}
}

func parseTIME(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}
	tstamp = fields[1]
}

func parseACK(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	i, _ := strconv.Atoi(fields[1])
	delete(kismet.pipeline, i)
	if kismet.debug {
		fmt.Println("ACK: " + fields[1])
		fmt.Println("PIPELINE:", kismet.pipeline)
	}
}

func parseINFO(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	packets, _ = strconv.Atoi(fields[2])
	packetRate, _ = strconv.Atoi(fields[3])
	f, _ := strconv.Atoi(fields[4])
	if f != 0 {
		filtered = true
	}

	fStr := "No"
	if filtered {
		fStr = "Yes"
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

	wsSend("kismetParseInfo", strconv.Itoa(len(networks))+":"+strconv.Itoa(len(clients))+":"+strconv.Itoa(rCount)+":"+fields[2]+":"+fields[3]+":"+fStr)
}

func parseSTATUS(fields []string) {
	if kismet.debug {
		fmt.Println("STATUS: " + strings.Join(fields, ","))
	}

	if len(fields) == 3 {
		switch {
		case fields[2] == "2":
			fmt.Println("STATUS-INFO: " + fields[1])

			re := regexp.MustCompile(`Created source (.*) with UUID (.*)$`)
			matches := re.FindStringSubmatch(fields[1])
			if matches != nil {
				fmt.Println(matches[1] + " : " + matches[2])
				iface, _ := net.InterfaceByName(matches[1])
				kismet.interfaces[matches[2]] = networkInterface{matches[1], "", iface.HardwareAddr.String(), "", "", "", "", map[string]string{}}
			} else if fields[1] == "Saved data files" {
				saveDB()
			} else {
				re = regexp.MustCompile(`Added source '(.*):name=(.*)' from client ADDSOURCE`)
				matches = re.FindStringSubmatch(fields[1])
				if matches != nil {
					fmt.Println(matches[1] + " : " + matches[2])
					for uid, ele := range kismet.interfaces {
						if ele.pname == matches[1] {
							ele.lname = matches[2]
							kismet.interfaces[uid] = ele
							break
						}
					}
				}
			}
		case fields[2] == "4":
			fmt.Println("STATUS-ERROR: " + fields[1])
		default:
			if kismet.debug {
				fmt.Println("STATUS-OTHER: (" + fields[2] + ") " + fields[1])
			}
		}
	}
}

func parseSOURCE(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	_, exists := kismet.interfaces[fields[5]]
	if !exists {
		iface, _ := net.InterfaceByName(fields[1])
		kismet.interfaces[fields[5]] = networkInterface{fields[1], fields[3], iface.HardwareAddr.String(), fields[4], fields[7], fields[8], fields[12], map[string]string{}}
	} else {
		nic := kismet.interfaces[fields[5]]
		nic.channel = fields[4]
		nic.hop = fields[7]
		nic.velocity = fields[8]
		nic.channellist = fields[12]
		kismet.interfaces[fields[5]] = nic
	}
	if fields[14] != "" {
		fmt.Println("NIC-WARNING: (" + fields[1] + ") " + fields[14])
	}
	wsSend("kismetParseSOURCE", fields[1]+";"+fields[3]+";"+fields[4]+";"+fields[7]+";"+fields[8]+";"+fields[12])
}

func parseALERT(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	fmt.Println("ALERT: " + fields[3] + " " + fields[9])
}

func parseERROR(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	i, _ := strconv.Atoi(fields[1])
	delete(kismet.pipeline, i)
	if kismet.debug {
		fmt.Println("ERROR: " + fields[2])
		fmt.Println("PIPELINE:", kismet.pipeline)
	}
}

func parseBSSIDSRC(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	_, exists := kismet.interfaces[fields[2]]
	if exists {
		kismet.interfaces[fields[2]].bssids[fields[1]] = fields[3]
	}
}

func parseBSSID(fields []string) {
	//"BSSID bssid,type,llcpackets,datapackets,cryptpackets,manuf,channel,firsttime,lasttime,atype,rangeip,netmaskip,gatewayip,gpsfixed,minlat,minlon,minalt,minspd,maxlat,maxlon,maxalt,maxspd,signal_dbm,noise_dbm,minsignal_dbm,minnoise_dbm,maxsignal_dbm,maxnoise_dbm,signal_rssi,noise_rssi,minsignal_rssi,minnoise_rssi,maxsignal_rssi,maxnoise_rssi,bestlat,bestlon,bestalt,agglat,agglon,aggalt,aggpoints,datasize,turbocellnid,turbocellmode,turbocellsat,carrierset,maxseenrate,encodingset,decrypted,dupeivpackets,bsstimestamp,cdpdevice,cdpport,fragments,retries,newpackets,freqmhz,datacryptset",
	if kismet.debug {
		fmt.Println(fields)
	}

	channel, _ := strconv.Atoi(fields[4])
	aType, _ := strconv.Atoi(fields[7])
	signalDBM, _ := strconv.Atoi(fields[11])
	minSignalDBM, _ := strconv.Atoi(fields[12])
	maxSignalDBM, _ := strconv.Atoi(fields[13])
	numPackets, _ := strconv.Atoi(fields[14])

	_, exists := accessPoints[fields[1]]
	if !exists {
		accessPoints[fields[1]] = accessPoint{fields[2], "", fields[3], channel, fields[5], fields[6], aType, fields[8], fields[9], fields[10], signalDBM, minSignalDBM, maxSignalDBM, numPackets}
	} else {
		ap := accessPoints[fields[1]]
		ap.apType = fields[2]
		ap.manuf = fields[3]
		ap.lasttime = fields[6]
		ap.atype = aType
		ap.rangeip = fields[8]
		ap.netmaskip = fields[9]
		ap.gatewayip = fields[10]
		ap.signalDBM = signalDBM
		ap.minSignalDBM = minSignalDBM
		ap.maxSignalDBM = maxSignalDBM
		ap.numPackets = numPackets
		accessPoints[fields[1]] = ap
	}

	clientCount := 0
	for _, client := range clients {
		if client.bssid == fields[1] {
			clientCount++
		}
	}

	lInt, _ := strconv.Atoi(fields[6])
	wsSend("kismetParseBSSID", fields[1]+";"+fields[4]+";"+time.Unix(int64(lInt), 0).Local().String()+";"+fields[11]+";"+strconv.Itoa(clientCount)+";"+fields[13]+";"+fields[14])
}

func parseSSID(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	if fields[4] != "" {
		cloaked, _ := strconv.Atoi(fields[7])
		maxrate, _ := strconv.Atoi(fields[10])

		_, exists := networks[fields[4]]
		if !exists {
			networks[fields[4]] = network{cloaked, fields[8], fields[9], maxrate, "", map[string]string{}}
		} else {
			network := networks[fields[4]]
			network.cloaked = cloaked
			network.lasttime = fields[9]
			network.maxrate = maxrate
			networks[fields[4]] = network
		}

		_, exists = networks[fields[4]].bssids[fields[1]]
		if !exists {
			networks[fields[4]].bssids[fields[1]] = fields[9]
		}
		wsSend("kismetParseSSID", fields[4]+";"+fields[1]+";"+fields[9])
	}
}

func parseCLISRC(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}
}

func parseCLIENT(fields []string) {
	if kismet.debug {
		fmt.Println(fields)
	}

	signalDBM, _ := strconv.Atoi(fields[7])
	minSignalDBM, _ := strconv.Atoi(fields[8])
	maxSignalDBM, _ := strconv.Atoi(fields[9])
	numPackets, _ := strconv.Atoi(fields[10])

	_, exists := clients[fields[2]]
	if !exists && fields[2] != fields[1] {
		clients[fields[2]] = client{fields[1], fields[4], fields[5], signalDBM, minSignalDBM, maxSignalDBM, numPackets}
	} else {
		cli := clients[fields[2]]
		cli.lasttime = fields[4]
		cli.signalDBM = signalDBM
		cli.minSignalDBM = minSignalDBM
		cli.maxSignalDBM = maxSignalDBM
		cli.numPackets = numPackets
		clients[fields[2]] = cli
	}
}

func parseTERMINATE(fields []string) {
	if kismet.debug {
		fmt.Println("TERMINATE")
	}
	wsSend("kismetParseTerminate", "DISCONNECTED")
}

func wsSend(msgType string, msg string) {
	if wsconn != nil {
		switch msgType {
		case "kismetParseInfo":
			if curPage == "/" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		case "kismetParseSSID":
			if curPage == "/discover" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		case "kismetParseBSSID":
			if curPage == "/discover" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		case "kismetParseSOURCE":
			if curPage == "/discover" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		default:
			wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
		}
	}
}

func processWSCommand(msg string) {
	s := strings.SplitN(msg, ":", 2)

	switch s[0] {
	case "kismetDISCONNECT":
		kismet.conn.Close()

	case "kismetCONNECT":
		Run(kismet.host, kismet.port, kismet.db, kismet.debug, ssids)

	case "statsNIC":
		iface, _ := net.InterfaceByName(s[1])
		active := "0"
		stats := ""
		for _, ele := range kismet.interfaces {
			if stats == "" && ele.hwaddr == iface.HardwareAddr.String() {
				active = "2"
				stats = ele.pname + ";" + ele.lname + ";" + ele.channel + ";" + ele.hop + ";" + ele.velocity + ";" + ele.channellist
			}
			if ele.pname == s[1] {
				active = "1"
				break
			}
		}
		stats = s[1] + ";" + active + ";" + stats
		wsSend("nicINFO", stats)

	case "statsSSID":
		wlan := networks[s[1]]

		channels := make(map[int]string)
		clientCount := 0
		min := 0
		max := -99
		for bssid := range wlan.bssids {
			ap := accessPoints[bssid]
			if ap.minSignalDBM < min {
				min = ap.minSignalDBM
			}
			if ap.maxSignalDBM > max {
				max = ap.maxSignalDBM
			}
			channels[ap.channel] = string(ap.channel)

			for clientMAC := range clients {
				if clientMAC == bssid {
					clientCount++
				}
			}
		}

		cloaked := "No"
		if wlan.cloaked == 1 {
			cloaked = "Yes"
		}

		keys := []string{}
		for k := range channels {
			keys = append(keys, strconv.Itoa(k))
		}

		i, _ := strconv.ParseInt(wlan.firsttime, 10, 64)
		firstTime := time.Unix(i, 0)

		i, _ = strconv.ParseInt(wlan.lasttime, 10, 64)
		lastTime := time.Unix(i, 0)

		wsSend("ssidINFO", s[1]+";"+cloaked+";"+firstTime.String()+";"+lastTime.String()+";"+strconv.Itoa(wlan.maxrate)+";"+strconv.Itoa(min)+";"+
			strconv.Itoa(max)+";"+strconv.Itoa(clientCount)+";"+strconv.Itoa(len(wlan.bssids))+";"+wlan.encryption+";"+strings.Join(keys, ","))

	case "nicADDSOURCE":
		t := strings.Split(s[1], ":")
		txt := t[0]
		if t[1] != "" {
			txt = txt + ":name=" + t[1]
		}
		Send("ADDSOURCE", txt)

	case "nicDELSOURCE":
		for uid, ele := range kismet.interfaces {
			if ele.pname == s[1] {
				Send("DELSOURCE", uid)
				delete(kismet.interfaces, uid)
				break
			}
		}
	}
}

func saveDB() {
	for ssid, network := range networks {
		kdb.InsertNetwork(kismet.db, ssid, network.cloaked, network.firsttime, network.lasttime, network.maxrate, network.encryption)
	}

	for bssid, ap := range accessPoints {
		kdb.InsertAP(kismet.db, bssid, ap.ssid, ap.channel, ap.firsttime, ap.lasttime, ap.minSignalDBM, ap.maxSignalDBM)
	}

	for mac, client := range clients {
		kdb.InsertClient(kismet.db, mac, client.firsttime, client.lasttime, client.signalDBM, client.minSignalDBM, client.maxSignalDBM)
	}
}

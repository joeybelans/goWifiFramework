// Manages server connection and data related to kismet server
package kismet

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/joeybelans/gokismet/kdb"
)

var (
	Version   string
	StartTime string
	Name      string
	dumpfiles []string
	uid       int
)

/*
// Setup the handler
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
	//setupClient()
}
*/

// Listen and process kismet data
func listen() {
	// Continuously read data
	scanner := bufio.NewScanner(kismet.conn)
	re := regexp.MustCompile(`([^ \001]+|\001[^\001]*\001)`)
	for scanner.Scan() {
		status := scanner.Text()

		// Determine response type
		status = strings.TrimSpace(status)
		matches := re.FindAllStringSubmatch(status, -1)
		fields := parseFields(matches)
		c, found := capabilities[fields[0]]
		if found {
			c.function.(func([]string))(fields)
		} else if fields[0] == "PROTOCOLS" {
			parsePROTOCOLS(fields)
		} else if fields[0] == "CAPABILITY" {
			parseCAPABILITY(fields)
		} else if kismet.debug {
			fmt.Println("UNKNOWN CMD: " + fields[0])
			fmt.Println(status)
		}
	}
	fmt.Println("QUITING")
}

// Send/Receive initial kismet commands/responses to setup the client communication
func setupClient() {
	// Disable TIME
	//Send("REMOVE", "TIME")

	// Enable default capabilities
	//KISMET,ERROR,ACK,PROTOCOLS,CAPABILITY,TERMINATE,TIME,PACKET,STATUS,PLUGIN,SOURCE,ALERT,COMMON,TRACKINFO,WEPKEY,STRING,GPS,BSSID,SSID,CLIENT,BSSIDSRC,CLISRC,NETTAG,CLITAG,REMOVE,CHANNEL,INFO,BATTERY,CRITFAIL
	/*
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
	*/

	// Process the pipeline of initial commands/resonses
	_, ok := kismet.pipeline[kismet.index]
	for ok || len(kismet.pipeline) > 0 {
		time.Sleep(1 * time.Second)
		_, ok = kismet.pipeline[kismet.index]
	}
}

// Dump the wireless data to sqlite file
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

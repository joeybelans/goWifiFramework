package kismet

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/joeybelans/gokismet/kdb"
)

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

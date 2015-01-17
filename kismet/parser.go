package kismet

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
		_, exists = accessPoints[fields[1]]
		if !exists {
			accessPoints[fields[1]] = accessPoint{}
		}
		ap := accessPoints[fields[1]]
		ap.ssid = fields[4]
		accessPoints[fields[1]] = ap
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

	if fields[1] != fields[2] {
		_, exists := clients[fields[2]]
		if !exists {
			clients[fields[2]] = client{fields[1], fields[4], fields[5], signalDBM, minSignalDBM, maxSignalDBM, numPackets}
		} else {
			cli := clients[fields[2]]
			cli.lasttime = fields[5]
			cli.signalDBM = signalDBM
			cli.minSignalDBM = minSignalDBM
			cli.maxSignalDBM = maxSignalDBM
			cli.numPackets = numPackets
			clients[fields[2]] = cli
		}

		lInt, _ := strconv.Atoi(fields[5])

		wsSend("kismetParseCLIENT", fields[2]+";"+time.Unix(int64(lInt), 0).Local().String()+";"+fields[7]+";"+fields[8]+";"+fields[9]+";"+fields[10])
	}
}

func parseTERMINATE(fields []string) {
	if kismet.debug {
		fmt.Println("TERMINATE")
	}
	wsSend("kismetParseTerminate", "DISCONNECTED")
}

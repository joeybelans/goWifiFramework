package webSocketHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	wsconn = nil
	curPage = ""
	packages = make(map[string]string)

	/*
		// Create component/page/data map
		dataMap = map[dataMapKey]bool{
			dataMapKey{"kismet", "INFO", "/kismet"}:   true,
			dataMapKey{"kismet", "BSSID", "/kismet"}:  true,
			dataMapKey{"kismet", "SSID", "/kismet"}:   true,
			dataMapKey{"kismet", "CLIENT", "/kismet"}: true,
		}
	*/
}

func AddPage(url string, pkg string) {
	packages[url] = pkg
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

	var data map[string]interface{}
	for {
		_, message, err := wsconn.ReadMessage()
		if err != nil {
			break
		}

		if err := json.Unmarshal(message, &data); err != nil {
			panic(err)
		}

		processWSCommand(data)
	}
}

func OKToSend(component string, message string) bool {
	/*
		_, exists := dataMap[dataMapKey{component, message, curPage}]
		return exists
	*/
	return true
}

func SendString(component string, message string, obj map[string]string) {
	msg := map[string]interface{}{"component": component, "message": message, "obj": obj}
	if wsconn != nil {
		wsconn.WriteJSON(msg)
	}
}

func SendInterface(component string, message string, obj map[string]interface{}) {
	msg := map[string]interface{}{"component": component, "message": message, "obj": obj}
	if wsconn != nil {
		wsconn.WriteJSON(msg)
	}
}

func wsSend(component string, message string, obj map[string]interface{}) {
	msg := map[string]interface{}{"component": component, "message": message, "obj": obj}
	if wsconn != nil {
		wsconn.WriteJSON(msg)
	}
}

func processWSCommand(data map[string]interface{}) {
	/*
		fmt.Println(packages[curPage])
		fmt.Println(data)
		 fields := parseFields(matches)
		                 c, found := parsers[fields[0]]
				                 if found {
						                            c.(func([]string))(fields)
									    var parsers = map[string]interface{}{
									               "KISMET":    parseKISMET,
	*/

	/*
		if data["type"].(string) == "home" {
			home.ProcessCommand(data)
		}
	*/

	/*
		switch data["type"].(string) {
		case "kismet":
			fmt.Println("PROCESS")
			//processKismet(data)
		}
			switch data["message"].(string) {
			case "kismetDISCONNECT":
				kismet.conn.Close()

			case "kismetCONNECT":
				Run(kismet.host, kismet.port, kismet.db, kismet.debug, ssids)

			case "statsNIC":
				nic := data["nic"].(string)
				stats := map[string]interface{}{
					"message":  "nicINFO",
					"nic":      nic,
					"active":   0,
					"physical": "",
					"alias":    "",
					"channel":  0,
					"hop":      0,
					"velocity": 0,
					"cList":    "",
				}
				iface, _ := net.InterfaceByName(nic)
				for _, ele := range kismet.interfaces {
					if stats["active"] == 0 && ele.hwaddr == iface.HardwareAddr.String() {
						stats["active"] = "2"
						stats["physical"] = ele.pname
						stats["alias"] = ele.lname
						stats["channel"] = ele.channel
						stats["hop"] = ele.hop
						stats["velocity"] = ele.velocity
						stats["cList"] = ele.channellist
					}
					if ele.pname == nic {
						stats["active"] = "1"
						break
					}
				}
				wsSend(stats)

			case "statsSSID":
				ssid := data["ssid"].(string)
				wlan := networks[ssid]

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

				wsSend(map[string]interface{}{
					"message":   "ssidINFO",
					"ssid":      ssid,
					"cloaked":   cloaked,
					"firstseen": firstTime.String(),
					"lastseen":  lastTime.String(),
					"maxrate":   wlan.maxrate,
					"min":       min,
					"max":       max,
					"clients":   clientCount,
					"aps":       len(wlan.bssids),
					"crypt":     wlan.encryption,
					"keys":      strings.Join(keys, ","),
				})

			case "nicADDSOURCE":
				txt := data["nic"].(string)
				if data["name"].(string) != "" {
					txt = txt + ":name=" + data["name"].(string)
				}
				Send("ADDSOURCE", txt)

			case "nicDELSOURCE":
				for uid, ele := range kismet.interfaces {
					if ele.pname == data["nic"].(string) {
						Send("DELSOURCE", uid)
						delete(kismet.interfaces, uid)
						break
					}
				}
			case "nicLOCK":
				for uid, ele := range kismet.interfaces {
					if ele.pname == data["nic"].(string) {
						Send("HOPSOURCE", uid+" LOCK "+data["channel"].(string))
						break
					}
				}
			case "nicHOP":
				for uid, ele := range kismet.interfaces {
					if ele.pname == data["nic"].(string) {
						Send("HOPSOURCE", uid+" HOP "+data["velocity"].(string))
						break
					}
				}
			case "nicCHANSOURCE":
				for uid, ele := range kismet.interfaces {
					if ele.pname == data["nic"].(string) {
						Send("CHANSOURCE", uid+" "+data["cList"].(string))
						break
					}
				}
			case "getAPDetails":
				ap := accessPoints[data["bssid"].(string)]
				wlan := networks[ap.ssid]

				i, _ := strconv.ParseInt(ap.firsttime, 10, 64)
				apFirstTime := time.Unix(i, 0)

				i, _ = strconv.ParseInt(ap.lasttime, 10, 64)
				apLastTime := time.Unix(i, 0)

				i, _ = strconv.ParseInt(wlan.firsttime, 10, 64)
				wlanFirstTime := time.Unix(i, 0)

				i, _ = strconv.ParseInt(wlan.lasttime, 10, 64)
				wlanLastTime := time.Unix(i, 0)

				wsSend(map[string]interface{}{
					"message":     "apDetails",
					"ssid":        ap.ssid,
					"cloaked":     wlan.cloaked,
					"firstseen":   wlanFirstTime.String(),
					"lastseen":    wlanLastTime.String(),
					"maxrate":     wlan.maxrate,
					"crypt":       wlan.encryption,
					"aps":         len(wlan.bssids),
					"bssid":       data["bssid"].(string),
					"aptype":      ap.apType,
					"manuf":       ap.manuf,
					"channel":     ap.channel,
					"apfirstseen": apFirstTime.String(),
					"aplastseen":  apLastTime.String(),
					"atype":       ap.atype,
					"ip":          ap.rangeip,
					"netmask":     ap.netmaskip,
					"gateway":     ap.gatewayip,
					"power":       ap.signalDBM,
					"min":         ap.minSignalDBM,
					"max":         ap.maxSignalDBM,
					"packets":     ap.numPackets,
				})
			}
	*/
}

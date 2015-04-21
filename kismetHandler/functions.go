// Manages server connection and data related to kismet server
package kismetHandler

func ServerVersion() string {
	return ""
}

func ServerName() string {
	return ""
}

func ServerStart() string {
	return ""
}

func Connected() bool {
	/*tstampInt, _ := strconv.ParseInt(tstamp, 10, 64)
	if time.Now().Unix()-tstampInt < 5 {
		return true
	}
	*/
	return false
}

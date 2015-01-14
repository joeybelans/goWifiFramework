package main

// Create networks structure
type bstatus struct {
	timestamp uint32
	power     uint8
}

type bssid struct {
	hidden  bool
	channel uint8
	status  []bstatus
}

type network struct {
	inscope    bool
	encryption string
	handshake  bool
	key        string
	eap        string
	bssids     map[string]bssid
}

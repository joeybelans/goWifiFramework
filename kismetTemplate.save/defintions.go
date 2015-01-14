package kismetTemplate

import "html/template"

type page struct {
	URL   string
	Title string
	iface interface{}
}

type header struct {
	Title     string
	Path      string
	Connected string
}

type home struct {
	ServerVersion string
	ServerName    string
	StartTxt      string
	DBFile        string
	SSIDs         []string
	Stats         map[string]int
	Interfaces    []string
}

type discover struct {
	DBFile string
	SSIDs  []string
}

var templates map[string]*template.Template

package httpTemplate

// Header
type header struct {
	Title string
}

func tmplHeader(pages map[int]page) string {
	txt := `
<html>
<head>
<title>GoKismet - {{.Title}}</title>
<link rel="stylesheet" type="text/css" href="/css/header.css">
<link rel="stylesheet" type="text/css" href="/css/nav.css">
<link rel="stylesheet" type="text/css" href="/css/home.css">
<link rel="stylesheet" type="text/css" href="/css/kismet.css">
<script src="/js/global.js"></script>
<script src="/js/kismet.js"></script>
<script src="/js/kismetParser.js"></script>
<script src="/js/home.js"></script>
</head>
<body>
<table border="0" width="100%" cellspacing="0" cellpadding="0"><tr><td>
<div id="hmenu"> 
<ul>`
	for i := 0; i < len(pages); i++ {
		txt = txt + "<li{{ if eq .Title \"" + pages[i].Title + "\" }} class=\"current-page\">" + pages[i].Title + "{{ else }}><a href='" + pages[i].URL + "'>" +
			pages[i].Title + "</a>{{ end }}</li>\n"
	}
	txt = txt + "</ul></div></td></tr></table><p>\n"
	return txt
}

/*
func tmplHeader(pages map[int]page) string {
	txt := `
<html>
<head>
<title>GoKismet - {{.Title}}</title>
<link rel="stylesheet" type="text/css" href="/css/header.css">
<link rel="stylesheet" type="text/css" href="/css/nav.css">
<link rel="stylesheet" type="text/css" href="/css/home.css">
<script src="/js/header.js"></script>
<script src="/js/kismetParser.js"></script>
<script src="/js/home.js"></script>
<script src="/js/discover.js"></script>
</head>
<body>
<table border="0" width="100%" cellspacing="0" cellpadding="0"><tr><td width="5%">
<div class="switch">
<input id="cmn-toggle-7" class="cmn-toggle cmn-toggle-on-off" type="checkbox" onClick="kismetOnOff()" {{.Connected}}>
<label for="cmn-toggle-7" data-on="Connected" data-off="Disconnected"></label>
</div>
</td><td>
<div id="hmenu">
<ul>`
	for i := 0; i < len(pages); i++ {
		txt = txt + "<li{{ if eq .Title \"" + pages[i].Title + "\" }} class=\"current-page\">" + pages[i].Title + "{{ else }}><a href='" + pages[i].URL + "'>" + pages[i].Title + "</a>{{ end }}</li>\n"
	}
	txt = txt + "</ul></div></td></tr></table><p>\n"
	return txt
}
*/

// Generates the header source code
// HTML template
package header

// Template structure
type templateData struct {
	Title string
	File  string
}

// Template source code
func templateSource(pages map[int]page) string {
	txt := `
<html>
<head>
<title>GoKismet - {{.Title}}</title>
<link rel="stylesheet" type="text/css" href="/css/global.css">
<link rel="stylesheet" type="text/css" href="/css/header.css">
<link rel="stylesheet" type="text/css" href="/css/{{.File}}.css">
<script src="/js/global.js"></script>
<script src="/js/{{.File}}.js"></script>
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

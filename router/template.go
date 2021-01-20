package router

import (
	"html/template"
	"time"
)

var defaultFuncs = template.FuncMap{
	"defTitle": func(ip interface{}) string {
		v, ok := ip.(string)
		if !ok || (ok && v == "") {
			return "Tasks - Golang Std HTTP"
		}
		return v
	},
	"dtestr": func(ip time.Time, fmt string) string {
		return ip.Format(fmt)
	},
}
var templateFiles = []string{
	"./web/templates/base.gohtml",
	"./web/templates/statusTag.gohtml",
}

func tmplLayout(files ...string) []string {
	return append(templateFiles, files...)
}

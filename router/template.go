package router

import "html/template"

var defaultFuncs = template.FuncMap{
	"defTitle": func(ip interface{}) string {
		v, ok := ip.(string)
		if !ok || (ok && v == "") {
			return "Tasks - Golang Std HTTP"
		}
		return v
	},
}
var templateFiles = []string{
	"./web/templates/base.gohtml",
}

func tmplLayout(files ...string) []string {
	return append(templateFiles, files...)
}

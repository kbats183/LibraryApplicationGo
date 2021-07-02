package templater

import (
	"fmt"
	"html/template"
	"io"
)

var templateCache = make(map[string]*template.Template)

type PageAlert struct {
	Status  string
	Content string
}

type TemplateData struct {
	PageName  string
	PageAlert PageAlert
	Data      interface{}
}

func ExecuteTemplate(templateName string, pageName string, writer io.Writer, data interface{}, alert *PageAlert) {
	tmpl, inCache := templateCache[templateName]
	if !inCache {
		var err error
		tmpl, err = template.ParseFiles("./template/"+templateName+".html", "./template/base.html")
		if err != nil {
			fmt.Println("Template parsing error:", err)
			return
		}
		templateCache[templateName] = tmpl
	}

	var pageAlert PageAlert
	if alert != nil {
		pageAlert = *alert
	}
	err := tmpl.Execute(writer, TemplateData{pageName, pageAlert, data})
	if err != nil {
		fmt.Println("Template executing  error:", err)
	}
}

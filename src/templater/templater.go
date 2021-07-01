package templater

import (
	"fmt"
	"html/template"
	"io"
)

type PageAlert struct {
	Status  string
	Content string
}

type TemplateData struct {
	BaseUrl   string
	PageName  string
	PageAlert PageAlert
	Data      interface{}
}

var baseUrl = "localhost/"

func ExecuteTemplate(templateName string, pageName string, writer io.Writer, data interface{}, alert *PageAlert) {
	tmpl, err := template.ParseFiles("./template/"+templateName+".html", "./template/base.html")
	if err != nil {
		fmt.Println("Template parsing error:", err)
	}

	var pageAlert PageAlert
	if alert != nil {
		pageAlert = *alert
	}
	err = tmpl.Execute(writer, TemplateData{baseUrl, pageName, pageAlert, data})
	if err != nil {
		fmt.Println("Template executing  error:", err)
	}
}

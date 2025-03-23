package views

import "html/template"

type Views struct {
}

func (this *Views) GetTemplates() *template.Template {
	templates, err := template.ParseGlob("views/*.html")
	if err != nil {
		panic(err.Error())
	}
	newTemplates := template.Must(templates, err)
	return newTemplates
}

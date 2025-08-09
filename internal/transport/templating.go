package transport

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templateFS embed.FS

func SetTemplateFS(router *gin.Engine) {
	templ := template.Must(template.ParseFS(templateFS, "templates/*.html"))
	router.SetHTMLTemplate(templ)
}

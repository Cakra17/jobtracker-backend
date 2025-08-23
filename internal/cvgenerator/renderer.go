package cvgenerator

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

type TemplateRenderer struct {
	templatePath string
}

func NewTemplateRenderer(templatePath string) *TemplateRenderer {
	return &TemplateRenderer{templatePath: templatePath}
}

func (tr *TemplateRenderer) RenderHTML(resume *Resume) (string, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/template.html",tr.templatePath))
  if err != nil {
    return "", err 
  }
  
  var buf bytes.Buffer
  err = tmpl.Execute(&buf, resume)
  if err != nil {
    return "", err
  }

  return buf.String(), nil
}

func (tr *TemplateRenderer) RenderHTMLWithCSS(resume *Resume) (string, error) {
  html, err := tr.RenderHTML(resume)
  if err != nil {
    return "", nil
  }

  cssContent, err := os.ReadFile(fmt.Sprintf("%s/style.css", tr.templatePath))
  html = strings.Replace(html, `<link rel="stylesheet" href="./style.css`, fmt.Sprintf("<style>%s</style>", cssContent), 1)
  return html, nil
}


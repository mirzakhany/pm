package engine

import (
	"os"
	"path"
	"strings"
	"text/template"
)

type Package struct {
	Name        string
	NamePlural  string
	EntityAlias string
}

type Data struct {
	Pkg Package
}

func Render(tpl, outputDir, fileName string, data Data) error {

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmpl, err := template.New(fileName).Funcs(funcMap).Parse(tpl)
	if err != nil {
		return err
	}

	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	outputFile := path.Join(outputDir, fileName)
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

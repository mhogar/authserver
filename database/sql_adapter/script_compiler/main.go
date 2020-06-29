package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

type tmplData struct {
	Name   string
	Script string
}

func main() {
	inDir := "database/sql_adapter/postgres/scripts"

	//load the template
	tmpl := template.Must(template.ParseFiles("database/sql_adapter/script_compiler/script_repository.go.tmpl"))

	//get all files in input dir
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}

	//read all sql files and create data objects
	var data []tmplData
	for _, file := range files {
		if path.Ext(file.Name()) == ".sql" {
			script, err := ioutil.ReadFile(path.Join(inDir, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			data = append(data, tmplData{
				Name:   strings.SplitN(path.Base(file.Name()), ".", 2)[0],
				Script: string(script),
			})
		}
	}

	//create the output file
	f, err := os.Create(path.Join(inDir, "script_repository.go"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//execute the template
	tmpl.Execute(f, data)
}

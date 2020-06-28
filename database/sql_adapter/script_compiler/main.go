package main

import (
	"fmt"
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
	inDir := "database/postgres_adapter/scripts"

	//load the template
	tmpl := template.Must(template.ParseFiles("database/sql_adapter/script_compiler/scripts.go.tmpl"))

	//get all files in input dir
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}

	//read all sql files and create data objects
	var datas []tmplData
	for _, file := range files {
		if path.Ext(file.Name()) == ".sql" {
			script, err := ioutil.ReadFile(path.Join(inDir, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			datas = append(datas, tmplData{
				Name:   strings.SplitN(path.Base(file.Name()), ".", 2)[0],
				Script: string(script),
			})
		}
	}

	//create the output file
	f, err := os.Create(path.Join(inDir, "scripts.go"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//write the header
	fmt.Fprintln(f, "// Auto generated. DO NOT EDIT.\n\npackage scripts")

	//loop through each data struct and execute the template
	for _, data := range datas {
		tmpl.Execute(f, data)
	}
}

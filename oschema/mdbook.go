package main

import (
	"bytes"
	"path/filepath"
)

type mdWriter struct {
	target string
	model  *YamlModel
}

func writeMarkdownBook(args *args) {
	model, err := ReadYamlModel(args.yamlDir)
	check(err, "could not read YAML model")
	target := args.target
	mkdir(target)
	writer := &mdWriter{model: model, target: target}
	writer.writeBook()
}

func (w *mdWriter) writeBook() {

	w.file("book.toml", `[book]
language = "en"
multilingual = false
src = "src"
title = "openLCA Schema"
`)

	w.dir("src")
	w.file("src/SUMMARY.md", w.summary())

	w.dir("src/classes")
	for _, t := range w.model.Types {
		if t.IsEnum() {
			continue
		}
		w.file("src/classes/"+t.Name()+".md", w.class(t.Class))
	}

	w.dir("src/enums")
	for _, t := range w.model.Types {
		if t.IsClass() {
			continue
		}
		w.file("src/enums/"+t.Name()+".md", w.enum(t.Enum))
	}

}

func (w *mdWriter) dir(path string) string {
	fullPath := filepath.Join(w.target, path)
	mkdir(fullPath)
	return fullPath
}

func (w *mdWriter) file(path, content string) {
	fullPath := filepath.Join(w.target, path)
	writeFile(fullPath, content)
}

func (w *mdWriter) summary() string {
	var buff bytes.Buffer
	buff.WriteString("# Classes\n\n")
	for _, t := range w.model.Types {
		if t.IsEnum() {
			continue
		}
		buff.WriteString(" - [" + t.Name() + "](./classes/" + t.Name() + ".md)\n")
	}

	buff.WriteString("\n# Enumerations\n\n")
	for _, t := range w.model.Types {
		if t.IsClass() {
			continue
		}
		buff.WriteString(" - [" + t.Name() + "](./enums/" + t.Name() + ".md)\n")
	}

	return buff.String()
}

func (w *mdWriter) class(class *YamlClass) string {
	return "# " + class.Name + "\n\n"
}

func (w *mdWriter) enum(enum *YamlEnum) string {
	return "# " + enum.Name + "\n\n"
}

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type mdWriter struct {
	target string
	model  *YamlModel
	args   *args
}

func writeMarkdownBook(args *args) {
	model, err := ReadYamlModel(args.yamlDir)
	check(err, "could not read YAML model")
	target := args.target
	mkdir(target)
	writer := &mdWriter{
		model:  model,
		target: target,
		args:   args}
	writer.writeBook()
}

func (w *mdWriter) writeBook() {

	w.file("book.toml", `[book]
language = "en"
multilingual = false
src = "src"
title = "openLCA Schema"

[output.html]
mathjax-support = true
`)

	w.dir("src")
	w.file("src/SUMMARY.md", w.summary())

	// try to copy the schema README
	readme := filepath.Join(filepath.Dir(w.args.yamlDir), "README.md")
	if _, err := os.Stat(readme); err == nil {
		if text, err := ioutil.ReadFile(readme); err == nil {
			w.file("src/Intro.md", string(text))
		} else {
			log.Println("WARNING: failed to copy", readme)
		}
	}

	w.dir("src/classes")
	for _, t := range w.model.Types {
		if t.IsEnum() {
			continue
		}
		w.file("src/classes/"+t.Name()+".md", w.docClassOf(t.Class))
	}

	w.dir("src/enums")
	for _, t := range w.model.Types {
		if t.IsClass() {
			continue
		}
		w.file("src/enums/"+t.Name()+".md", w.docEnumOf(t.Enum))
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

	buff := NewBuffer()
	buff.Writeln("# Summary\n")
	buff.Writeln("[Introduction](./Intro.md)")

	buff.Writeln("# Classes\n")
	innerTypes := w.innerTypes()
	for _, t := range w.model.Types {
		if t.IsEnum() || innerTypes[t.Name()] != "" {
			continue
		}

		buff.Writeln(" - [" + t.Name() + "](./classes/" + t.Name() + ".md)")
		for _, inner := range w.model.Types {
			if inner.IsEnum() {
				continue
			}
			if innerTypes[inner.Name()] == t.Name() {
				buff.Writeln("   - [" + inner.Name() + "](./classes/" +
					inner.Name() + ".md)\n")
			}
		}
	}

	buff.Writeln("\n# Enumerations\n")
	for _, t := range w.model.Types {
		if t.IsClass() {
			continue
		}
		buff.Writeln(" - [" + t.Name() + "](./enums/" + t.Name() + ".md)")
	}

	return buff.String()
}

func (w *mdWriter) docClassOf(class *YamlClass) string {
	var buff bytes.Buffer
	buff.WriteString("# " + class.Name + "\n\n")
	buff.WriteString(class.Doc + "\n\n")

	buff.WriteString("## Properties\n\n")

	parents := make([]*YamlClass, 0)
	parent := w.model.ParentOf(class)
	for {
		if parent == nil {
			break
		}
		parents = append([]*YamlClass{parent}, parents...)
		parent = w.model.ParentOf(parent)
	}

	for _, p := range parents {
		for _, prop := range p.Props {
			buff.WriteString("### `" + prop.Name + "`\n\n")
			buff.WriteString("Inherited from [" + p.Name + "." + prop.Name +
				"](./" + p.Name + ".md#" + prop.Name + ")\n\n")
			buff.WriteString(w.docPropOf(prop))
		}
	}

	for _, prop := range class.Props {
		buff.WriteString("### `" + prop.Name + "`\n\n")
		if prop.Doc != "" {
			buff.WriteString(prop.Doc + "\n\n")
		}
		buff.WriteString(w.docPropOf(prop))
	}

	buff.WriteString("## Python class stub\n\n")
	buff.WriteString("\n\n```python\n\n")
	buff.WriteString(w.model.ToPyClass(class))
	buff.WriteString("\n```\n")

	return buff.String()
}

func (w *mdWriter) docPropOf(prop *YamlProp) string {
	var buff bytes.Buffer
	if prop.Required {
		buff.WriteString("* _is required_\n")
	} else {
		buff.WriteString("* _is optional_\n")
	}
	buff.WriteString("* _Type:_ " + w.docTypeOf(prop.Type) + "\n")
	buff.WriteString("* _Proto-Index:_ " + strconv.Itoa(prop.Index) + "\n")
	return buff.String()
}

func (w *mdWriter) docEnumOf(enum *YamlEnum) string {
	var buff bytes.Buffer
	buff.WriteString("# " + enum.Name + "\n\n")
	buff.WriteString(enum.Doc + "\n\n")

	buff.WriteString("## Items\n\n")

	for _, item := range enum.Items {
		buff.WriteString("### `" + item.Name + "`\n\n")
		if item.Doc != "" {
			buff.WriteString(item.Doc + "\n\n")
		}
		buff.WriteString("* _Proto-Index:_ " + strconv.Itoa(item.Index) + "\n")
	}

	return buff.String()
}

func (w *mdWriter) docTypeOf(yamlType string) string {

	if yamlType == "" {
		return "__ERROR! EMPTY__"
	}

	if strings.HasPrefix(yamlType, "List[") {
		unpacked := strings.TrimPrefix(strings.TrimSuffix(yamlType, "]"), "List[")
		return "`List` of " + w.docTypeOf(unpacked)
	}

	if strings.HasPrefix(yamlType, "Ref[") {
		unpacked := strings.TrimPrefix(strings.TrimSuffix(yamlType, "]"), "Ref[")
		return "[Ref](./Ref.md) of " + w.docTypeOf(unpacked)
	}

	if yamlType == "GeoJSON" {
		return "`GeoJSON` ([external doc](https://tools.ietf.org/html/rfc7946))"
	}

	if startsWithLower(yamlType) {
		return "`" + yamlType +
			"` ([external doc](http://www.w3.org/TR/xmlschema-2/#" + yamlType + "))"
	}

	t := w.model.TypeMap[yamlType]
	if t == nil {
		log.Println("WARNING: unknown type:", yamlType)
		return "`" + yamlType + "`"
	}
	if t.IsEnum() {
		return "[" + yamlType + "](../enums/" + yamlType + ".md)"
	} else {
		return "[" + yamlType + "](./" + yamlType + ".md)"
	}

}

// Returns a map `inner type -> outer type` of types that are only used in
// in a specific outer type (like Exchange in Processes).
func (w *mdWriter) innerTypes() map[string]string {
	m := make(map[string]string)
	for _, inner := range w.model.Types {
		if inner.IsEnum() {
			continue
		}
		parent := w.model.ParentOf(inner.Class)
		if parent == nil || parent.Name == "RootEntity" {
			continue
		}

		matches := func(outer *YamlClass) bool {
			for _, prop := range outer.Props {
				propType := prop.Type
				if strings.HasPrefix(propType, "List[") {
					propType = strings.TrimPrefix(
						strings.TrimSuffix(propType, "]"), "List[")
				}
				if strings.HasPrefix(propType, "Ref[") {
					propType = strings.TrimPrefix(
						strings.TrimSuffix(propType, "]"), "Ref[")
				}
				if propType == inner.Name() {
					return true
				}
			}
			return false
		}

		candidate := ""
		for _, outer := range w.model.Types {
			if outer.IsEnum() {
				continue
			}
			if !matches(outer.Class) {
				continue
			}
			if candidate == "" {
				candidate = outer.Name()
			} else {
				candidate = ""
				break
			}
		}

		if candidate != "" {
			m[inner.Name()] = candidate
		}

	}

	return m
}

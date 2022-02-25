package main

import (
	"bytes"
	"fmt"
	"log"
)

type pyWriter struct {
	buff  *bytes.Buffer
	model *YamlModel
}

func writePythonModule(args *args) {
	model, err := ReadYamlModel(args.yamlDir)
	check(err, "could not read YAML model")

	var buffer bytes.Buffer
	writer := pyWriter{
		buff:  &buffer,
		model: model,
	}
	writer.writeModel()

	if args.target != "" {
		writeFile(args.target, buffer.String())
	} else {
		fmt.Println(buffer.String())
	}
}

func (w *pyWriter) writeModel() {

	w.writeln("# DO NOT CHANGE THIS CODE AS THIS IS GENERATED AUTOMATICALLY")
	w.writeln(`
# This module contains a Python API for reading and writing data sets in
# the JSON based openLCA data exchange format. For more information see
# http://greendelta.github.io/olca-schema
`)

	// imports
	w.writeln("from enum import Enum")
	w.writeln("from dataclasses import dataclass")
	w.writeln("from typing import Dict, List, Any")
	w.writeln()
	w.writeln()

	// enums and classes
	w.model.EachEnum(w.writeEnum)
	for _, class := range topoSortClasses(w.model) {
		w.writeClass(class)
	}
}

func (w *pyWriter) writeEnum(enum *YamlEnum) {
	w.writeln("class", enum.Name+"(Enum):")
	w.writeln()
	for _, item := range enum.Items {
		w.writeln("    " + item.Name + " = '" + item.Name + "'")
	}
	w.writeln()
	w.writeln()
}

func (w *pyWriter) writeClass(class *YamlClass) {
	w.writeln("@dataclass")
	w.writeln("class", class.Name+":")
	w.writeln()
	for _, prop := range w.model.AllPropsOf(class) {
		propName := prop.Name
		switch propName {
		case "@type":
			continue
		case "@id":
			propName = "id"
		case "from":
			propName = "from_"
		}
		propType := YamlPropType(prop.Type)
		w.writeln("    " + toSnakeCase(propName) + ": " + propType.ToPython())
	}
	w.writeln("    schema_type: str = '" + class.Name + "'")

	w.writeln()
	w.writeln()
}

func (w *pyWriter) writeln(args ...string) {
	for i, arg := range args {
		if i > 0 {
			w.buff.WriteRune(' ')
		}
		w.buff.WriteString(arg)
	}
	w.buff.WriteRune('\n')
}

func topoSortClasses(model *YamlModel) []*YamlClass {

	// check if there is a link between a class A and
	// another class B where B is dependent from A
	isLinked := func(class, dependent *YamlClass) bool {
		if class == dependent {
			return false
		}
		for _, prop := range dependent.Props {
			propType := YamlPropType(prop.Type)
			if propType.IsList() {
				propType = propType.UnpackList()
			}
			if propType.ToPython() == class.Name {
				return true
			}
		}
		return false
	}

	// collect the dependencies
	dependencyCount := make(map[string]int)
	dependents := make(map[string][]string)
	model.EachClass(func(class *YamlClass) {
		if _, ok := dependencyCount[class.Name]; !ok {
			dependencyCount[class.Name] = 0
		}
		model.EachClass(func(dependent *YamlClass) {
			if isLinked(class, dependent) {
				c := class.Name
				d := dependent.Name
				dependencyCount[d] += 1
				dependents[c] = append(dependents[c], d)
			}
		})
	})

	// sort dependencies in topological order
	order := make([]string, 0)
	for len(dependencyCount) > 0 {

		// find next node with no dependencies
		node := ""
		for n, count := range dependencyCount {
			if count <= 0 {
				node = n
				break
			}
		}

		if node == "" {
			log.Println("ERROR: could not sort classes in topological order")
			break
		}
		delete(dependencyCount, node)
		order = append(order, node)

		// remove the handled dependency from its dependents
		for _, dependent := range dependents[node] {
			dependencyCount[dependent] -= 1
		}
	}

	sorted := make([]*YamlClass, 0, len(order))
	for _, name := range order {
		next := model.TypeMap[name]
		if next != nil && next.IsClass() {
			sorted = append(sorted, next.Class)
		}
	}
	return sorted
}
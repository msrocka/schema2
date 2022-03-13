package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
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
	w.writeln("import json")
	w.writeln()
	w.writeln("from enum import Enum")
	w.writeln("from dataclasses import dataclass")
	w.writeln("from typing import Any, Dict, List, Optional, Union")
	w.writeln()
	w.writeln()

	// enums and classes
	w.model.EachEnum(w.writeEnum)
	for _, class := range topoSortClasses(w.model) {
		if w.model.IsAbstract(class) {
			continue
		}
		w.writeln(w.model.ToPyClass(class))
	}

	// write RootEntity type
	w.writeln("RootEntity = Union[")
	w.model.EachClass(func(class *YamlClass) {
		if w.model.IsRoot(class) {
			w.writeln("    " + class.Name + ",")
		}
	})
	w.writeln("]")
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

func (model *YamlModel) ToPyClass(class *YamlClass) string {
	b := NewBuffer()
	b.Writeln("@dataclass")
	b.Writeln("class", class.Name+":")
	b.Writeln()

	// properties
	props := model.AllPropsOf(class)
	for _, prop := range props {
		if prop.Name == "@type" {
			continue
		}
		propType := YamlPropType(prop.Type)
		b.Writeln("    " + prop.PyName() +
			": Optional[" + propType.ToPython() + "] = None")
	}
	if class.Name == "Ref" {
		b.Writeln("    model_type: str = ''")
	}
	b.Writeln()

	// to_dict
	b.Writeln("    def to_dict(self) -> Dict[str, Any]:")
	b.Writeln("        d: Dict[str, Any] = {}")
	if model.IsRoot(class) {
		b.Writeln("        d['@type'] = '" + class.Name + "'")
	}
	if class.Name == "Ref" {
		b.Writeln("        d['@type'] = self.model_type")
	}
	for _, prop := range props {
		if prop.Name == "@type" {
			continue
		}
		selfProp := "self." + prop.PyName()
		dictProp := "            d['" + prop.Name + "']"
		propType := prop.PropType()
		b.Writeln("        if " + selfProp + ":")
		if propType.IsPrimitive() ||
			propType.IsEnumOf(model) ||
			(propType.IsList() && propType.UnpackList().IsPrimitive()) ||
			propType == "GeoJSON" {
			b.Writeln(dictProp + " = " + selfProp)
		} else if propType.IsList() {
			b.Writeln(dictProp + " = [e.to_dict() for e in " + selfProp + "]")
		} else {
			b.Writeln(dictProp + " = " + selfProp + ".to_dict()")
		}
	}
	b.Writeln("        return d")
	b.Writeln()

	// to_json
	if model.IsRoot(class) {
		b.Writeln("    def to_json(self) -> str:")
		b.Writeln("        return json.dumps(self.to_dict(), indent=2)")
		b.Writeln()
	}

	// to_ref
	if model.IsRoot(class) || class.Name == "Unit" {
		b.Writeln("    def to_ref(self) -> 'Ref':")
		b.Writeln("        ref = Ref(id=self.id, name=self.name)")
		b.Writeln("        ref.category = self.category")
		b.Writeln("        ref.model_type = '" + class.Name + "'")
		b.Writeln("        return ref")
		b.Writeln()
	}

	// from_dict
	b.Writeln("    @staticmethod")
	b.Writeln("    def from_dict(d: Dict[str, Any]) -> '" + class.Name + "':")
	instance := strings.ToLower(toSnakeCase(class.Name))
	b.Writeln("        " + instance + " = " + class.Name + "()")
	if class.Name == "Ref" {
		b.Writeln("        " + instance + ".model_type = d.get('@type', '')")
	}
	for _, prop := range props {
		b.Writeln("        if v := d.get('" + prop.Name + "'):")
		propType := prop.PropType()
		modelProp := "            " + instance + "." + prop.PyName()
		if propType.IsPrimitive() ||
			propType.IsEnumOf(model) ||
			(propType.IsList() && propType.UnpackList().IsPrimitive()) ||
			propType == "GeoJSON" {
			b.Writeln(modelProp + " = v")
		} else if propType.IsList() {
			u := propType.UnpackList()
			b.Writeln(modelProp + " = [" + string(u) + ".from_dict(e) for e in v]")
		} else {
			b.Writeln(modelProp + " = " + string(propType) + ".from_dict(v)")
		}
	}
	b.Writeln("        return " + instance)
	b.Writeln()

	// from_json
	if model.IsRoot(class) {
		b.Writeln("    @staticmethod")
		b.Writeln("    def from_json(data: Union[str, bytes]) -> '" +
			class.Name + "':")
		b.Writeln("        return " + class.Name + ".from_dict(json.loads(data))")
		b.Writeln()
	}

	return b.String()
}

func (w *pyWriter) writeln(args ...string) {
	w.write(args...)
	w.buff.WriteRune('\n')
}

func (w *pyWriter) write(args ...string) {
	for i, arg := range args {
		if i > 0 {
			w.buff.WriteRune(' ')
		}
		w.buff.WriteString(arg)
	}
}

func topoSortClasses(model *YamlModel) []*YamlClass {

	// check if there is a link between a class A and another class B where B is
	// dependent from A. B is dependent from A if it has a property of type A.
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

	// make sure that every RootEntity is dependent from 'Ref' as we generate a
	// to_ref method where the Ref type should be known
	refDeps, ok := dependents["Ref"]
	if !ok {
		refDeps = make([]string, 0)
	}
	model.EachClass(func(class *YamlClass) {
		if !model.IsRoot(class) && class.Name != "Unit" {
			return
		}
		contains := false
		for _, dep := range refDeps {
			if class.Name == dep {
				contains = true
				break
			}
		}
		if !contains {
			refDeps = append(refDeps, class.Name)
			dependencyCount[class.Name] += 1
		}
	})
	dependents["Ref"] = refDeps

	// sort dependencies in topological order
	order := make([]string, 0)
	for len(dependencyCount) > 0 {

		// find next node with no dependencies; if there are multiple options, try
		// to do this in alphabetical order so that we get a stable sort order
		node := ""
		for n, count := range dependencyCount {
			if count > 0 {
				continue
			}
			if node == "" ||
				strings.Compare(strings.ToLower(n), strings.ToLower(node)) < 0 {
				node = n
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

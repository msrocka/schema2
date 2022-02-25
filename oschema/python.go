package main

import (
	"bytes"
	"fmt"
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

	// imports
	w.writeln("from enum import Enum")
	w.writeln()
	w.writeln()

	for _, t := range w.model.Types {
		if t.IsEnum() {
			w.writeEnum(t.Enum)
		}
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

func (w *pyWriter) writeln(args ...string) {
	for i, arg := range args {
		if i > 0 {
			w.buff.WriteRune(' ')
		}
		w.buff.WriteString(arg)
	}
	w.buff.WriteRune('\n')
}

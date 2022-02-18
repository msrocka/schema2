package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type YamlType struct {
	Class *YamlClass `yaml:"class"`
	Enum  *YamlEnum  `yaml:"enum"`
}

func (yt *YamlType) String() string {
	if yt.Class != nil {
		return "ClassDef " + yt.Class.Name
	}
	if yt.Enum != nil {
		return "EnumDef " + yt.Enum.Name
	}
	return "Unknown TypeDef"
}

func (yt *YamlType) name() string {
	if yt.Class != nil {
		return yt.Class.Name
	}
	if yt.Enum != nil {
		return yt.Enum.Name
	}
	return "Unknown"
}

type YamlClass struct {
	Name       string      `yaml:"name"`
	SuperClass string      `yaml:"superClass"`
	Doc        string      `yaml:"doc"`
	Fields     []*YamlProp `yaml:"properties"`
}

type YamlProp struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Doc  string `yaml:"doc"`
}

type YamlEnum struct {
	Name  string          `yaml:"name"`
	Doc   string          `yaml:"doc"`
	Items []*YamlEnumItem `yaml:"items"`
}

type YamlEnumItem struct {
	Name string `yaml:"name"`
	Doc  string `yaml:"doc"`
}

type YamlModel struct {
	Types   []*YamlType
	TypeMap map[string]*YamlType
}

func ReadYamlModel(dir string) (*YamlModel, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	types := make([]*YamlType, 0)
	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".yaml") {
			continue
		}

		log.Println("Parse YAML file", name)
		path := filepath.Join(dir, name)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		typeDef := &YamlType{}
		if err := yaml.Unmarshal(data, typeDef); err != nil {
			return nil, err
		}

		types = append(types, typeDef)
	}
	log.Println("Collected", len(types), "YAML types")

	typeMap := make(map[string]*YamlType)
	for i := range types {
		typeDef := types[i]
		typeMap[typeDef.name()] = typeDef
	}

	model := YamlModel{Types: types, TypeMap: typeMap}

	return &model, nil
}

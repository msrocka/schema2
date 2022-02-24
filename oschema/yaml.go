package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type YamlType struct {
	Class *YamlClass `yaml:"class"`
	Enum  *YamlEnum  `yaml:"enum"`
}

func (yt *YamlType) IsClass() bool {
	return yt.Class != nil
}

func (yt *YamlType) IsEnum() bool {
	return yt.Enum != nil
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

func (yt *YamlType) Name() string {
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
	Props      []*YamlProp `yaml:"properties"`
}

type YamlProp struct {
	Name     string `yaml:"name"`
	Index    int    `yaml:"index"`
	Type     string `yaml:"type"`
	Doc      string `yaml:"doc"`
	Required bool   `yaml:"required"`
}

type YamlPropsByName []*YamlProp

func (s YamlPropsByName) Len() int { return len(s) }
func (s YamlPropsByName) Less(i, j int) bool {
	name_i := s[i].Name
	name_j := s[j].Name
	if name_i == name_j {
		return false
	}
	firstOrder := []string{"@type", "@id"}
	for _, f := range firstOrder {
		if name_i == f {
			return true
		}
		if name_j == f {
			return false
		}
	}
	return strings.Compare(name_i, name_j) < 0
}
func (s YamlPropsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type YamlEnum struct {
	Name  string          `yaml:"name"`
	Doc   string          `yaml:"doc"`
	Items []*YamlEnumItem `yaml:"items"`
}

type YamlEnumItem struct {
	Name  string `yaml:"name"`
	Doc   string `yaml:"doc"`
	Index int    `yaml:"index"`
}

type YamlModel struct {
	Types   []*YamlType
	TypeMap map[string]*YamlType
}

func (model *YamlModel) ParentOf(class *YamlClass) *YamlClass {
	parentName := class.SuperClass
	if parentName == "" {
		return nil
	}
	parent := model.TypeMap[parentName]
	if parent == nil || parent.IsEnum() {
		return nil
	}
	return parent.Class
}

func (model *YamlModel) IsEmpty() bool {
	return len(model.Types) == 0
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
		typeMap[typeDef.Name()] = typeDef
	}

	model := YamlModel{Types: types, TypeMap: typeMap}

	return &model, nil
}

func (model *YamlModel) AllPropsOf(class *YamlClass) []*YamlProp {
	props := make([]*YamlProp, 0, len(class.Props)+1)
	c := class
	for {
		if c == nil {
			break
		}
		props = append(props, c.Props...)
		c = model.ParentOf(c)
	}
	sort.Sort(YamlPropsByName(props))
	return props
}

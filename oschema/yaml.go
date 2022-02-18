package main

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

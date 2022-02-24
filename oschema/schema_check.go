package main

import (
	"fmt"
	"sort"
	"strings"
)

func checkSchema(args *args) {
	yamlModel, err := ReadYamlModel(args.yamlDir)
	if err != nil {
		fmt.Println("ERROR: Failed to parse YAML model:", err)
		return
	}

	// check that every class begins in Entity
	for _, t := range yamlModel.Types {
		if t.IsEnum() {
			continue
		}
		class := t.Class
		for {
			if class.Name == "Entity" {
				break
			}
			parent := yamlModel.ParentOf(class)
			if parent == nil {
				fmt.Println("ERROR: class hierarchy of '" +
					class.Name + "' does not starts in `Entity`")
				break
			}
			class = parent
		}
	}

	// Check if the properties in the classes are sorted by name. This is just for
	// the initial schema creation and should be removed later.
	for _, t := range yamlModel.Types {
		if t.IsEnum() {
			continue
		}
		c := t.Class
		props := c.Props
		sorted := true
		var last *YamlProp
		for i := range props {
			prop := props[i]
			if i == 0 {
				last = props[i]
				continue
			}
			if strings.Compare(prop.Name, last.Name) < 0 || prop.Index < last.Index {
				sorted = false
				break
			}
		}

		if sorted {
			continue
		}

		fmt.Println("WARNING: properties not in order:", c.Name)
		sort.Sort(YamlPropsByName(c.Props))
		for i, p := range c.Props {
			fmt.Println("  o + ", i, p.Name)
		}
	}

}

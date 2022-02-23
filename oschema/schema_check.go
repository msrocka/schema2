package main

import "fmt"

func checkSchema(args *args) {
	yamlModel, err := ReadYamlModel(args.yamlDir)
	if err != nil {
		fmt.Println("ERROR: Failed to parse YAML model:", err)
		return
	}

	// check that every class ends in Entity
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
			}
		}
	}

}

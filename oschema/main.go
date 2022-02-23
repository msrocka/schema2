package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	// parse the YAML files
	yamlDir := findYamlDir()
	yamlModel, err := ReadYamlModel(yamlDir)
	check(err)

	proto := GenProto(yamlModel)

	// print to console or write to file
	if len(os.Args) < 3 {
		fmt.Println(proto)
	} else {
		outFile := os.Args[2]
		err := ioutil.WriteFile(outFile, []byte(proto), os.ModePerm)
		check(err, "failed to write to file", outFile)
	}
}

func check(err error, msg ...interface{}) {
	if err != nil {
		fmt.Print("ERROR: ")
		fmt.Println(msg...)
		panic(err)
	}
}

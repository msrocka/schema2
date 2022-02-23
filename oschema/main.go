package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "help", "-h":
		printHelp()
	case "proto":
		proto()
	case "check":
		checkSchema()
	}

}

func check(err error, msg ...interface{}) {
	if err != nil {
		fmt.Print("ERROR: ")
		fmt.Println(msg...)
		panic(err)
	}
}

func proto() {
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

func printHelp() {
	fmt.Println(`
oschema

usage:

$ oschema [command] [options]

commands:

  help  - prints this help
  check - checks the schema
  proto - converts the schema to ProtocolBuffers

  `)
}

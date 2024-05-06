package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Contract struct {
	State  []StateVariable             `yaml:"state"`
	Defs   []Definition                `yaml:"defs"` // Changed from Functions to Defs
	Events map[string][]EventParameter `yaml:"events"`
}

type StateVariable struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Payable bool   `yaml:"payable,omitempty"`
}

type Definition struct { // Renamed from Function to Definition
	Name       string              `yaml:"name"`
	Visibility []string            `yaml:"visibility,omitempty"` // Moved inside Definition
	Type       string              `yaml:"type,omitempty"`       // Made optional
	Parameters []FunctionParameter `yaml:"parameters,omitempty"` // Made optional
	Locals     []LocalVariable     `yaml:"locals,omitempty"`     // Made optional
}

type FunctionParameter struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type LocalVariable struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	InMemory bool   `yaml:"inMemory,omitempty"`
}

type EventParameter struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func main() {
	data, err := ioutil.ReadFile("playground.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var contract Contract
	err = yaml.Unmarshal(data, &contract)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Parsed Contract: %+v\n", contract.Events)
}

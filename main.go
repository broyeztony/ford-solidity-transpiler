package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Define a generic AST node structure
type ASTNode map[string]interface{}

// Define structures for the output AST
type ContractDefinition struct {
	Type          string        `json:"type"`
	Name          string        `json:"name"`
	BaseContracts []interface{} `json:"baseContracts"`
	SubNodes      []interface{} `json:"subNodes"`
	Kind          string        `json:"kind"`
}

type StateVariableDeclaration struct {
	Type         string      `json:"type"`
	Variables    []Variable  `json:"variables"`
	InitialValue interface{} `json:"initialValue"`
}

type Variable struct {
	Type            string             `json:"type"`
	TypeName        VariableTypeName   `json:"typeName"`
	Name            string             `json:"name"`
	Identifier      VariableIdentifier `json:"identifier"`
	Expression      interface{}        `json:"expression"`
	Visibility      string             `json:"visibility"`
	IsStateVar      bool               `json:"isStateVar"`
	IsDeclaredConst bool               `json:"isDeclaredConst"`
	IsIndexed       bool               `json:"isIndexed"`
	IsImmutable     bool               `json:"isImmutable"`
	Override        interface{}        `json:"override"`
	StorageLocation interface{}        `json:"storageLocation"`
}

type VariableTypeName struct {
	Type            string      `json:"type"`
	Name            string      `json:"name"`
	StateMutability interface{} `json:"stateMutability"`
}

type VariableIdentifier struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type VariableExpression struct {
	Type            string      `json:"type"`
	Value           interface{} `json:"value,omitempty"`
	Number          string      `json:"number,omitempty"`
	Parts           []string    `json:"parts,omitempty"`
	IsUnicode       []bool      `json:"isUnicode,omitempty"`
	Subdenomination interface{} `json:"subdenomination,omitempty"`
}

type VariableInitialValue struct {
	Type            string      `json:"type"`
	Value           interface{} `json:"value,omitempty"`
	Number          string      `json:"number,omitempty"`
	Parts           []string    `json:"parts,omitempty"`
	IsUnicode       []bool      `json:"isUnicode,omitempty"`
	Subdenomination interface{} `json:"subdenomination,omitempty"`
}

func main() {

	astFilePath := "data/primitives.ast.json" // AST
	inputAST, err := loadAndUnmarshalJSON(astFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	outputAST := transpileAST(inputAST)
	outputJSON, err := json.MarshalIndent(outputAST, "", "  ")
	if err != nil {
		fmt.Println("Error generating output JSON:", err)
		return
	}

	// println(len(outputJSON))
	fmt.Println(string(outputJSON))
}

func loadAndUnmarshalJSON(filePath string) (ASTNode, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading JSON file: %w", err)
	}

	var inputAST ASTNode
	if err := json.Unmarshal(jsonData, &inputAST); err != nil {
		return nil, fmt.Errorf("parsing input JSON: %w", err)
	}

	return inputAST, nil
}

// transpileAST transpiles the input AST to the desired output format.
func transpileAST(inputAST ASTNode) ContractDefinition {
	contractName, _ := inputAST["name"].(string)
	outputAST := ContractDefinition{
		Type:          "ContractDefinition",
		Name:          contractName,
		BaseContracts: []interface{}{},
		SubNodes:      []interface{}{},
		Kind:          "contract",
	}

	// Process each body element based on its type
	for _, bodyElement := range inputAST["body"].([]interface{}) {
		bodyNode := bodyElement.(map[string]interface{})
		switch bodyNode["type"] {
		case "VariableStatement":
			processVariableStatement(&outputAST, bodyNode)
		case "FunctionDeclaration":
			processFunctionDeclaration(&outputAST, bodyNode)
		}
	}

	return outputAST
}

// processVariableStatement processes a variable statement and adds it to the contract definition.
func processVariableStatement(contractDef *ContractDefinition, node map[string]interface{}) {

	varDeclarations, ok := node["declarations"].([]interface{})
	if !ok {
		fmt.Println("Error: 'declarations' is not a valid node.")
		return
	}

	for _, declInterface := range varDeclarations {

		decl, ok := declInterface.(map[string]interface{})
		if !ok {
			fmt.Println("Error: Declaration is not a valid map.")
			continue
		}

		varIdName, _ := decl["id"].(map[string]interface{})["name"].(string)
		varInitializer := decl["initializer"].(map[string]interface{})

		var varInitializerType string
		var varInitialValue interface{} = nil
		var varExpression interface{} = nil

		// StateVariableDeclaration
		svd := StateVariableDeclaration{
			Type: "StateVariableDeclaration",
			Variables: []Variable{{
				Type: "VariableDeclaration",
				TypeName: VariableTypeName{
					Type:            "ElementaryTypeName",
					StateMutability: nil,
				},
				Name: varIdName,
				Identifier: VariableIdentifier{
					Type: "Identifier",
					Name: varIdName,
				},

				Visibility: "public",
				IsStateVar: true,
			}},
		}

		if varInitializer != nil {

			varInitializerType, _ = declInterface.(map[string]interface{})["initializer"].(map[string]interface{})["type"].(string)

			switch varInitializerType {
			case "CallExpression":
				processCallExpressionInitializer(varInitializer, &svd)
			case "StringLiteral", "BooleanLiteral":
				processSimpleInitializer(varInitializer, &svd)
			default:
				fmt.Printf("Unrecognized initializer type: %s\n", varInitializerType)
			}
		}

		if varExpression != nil {
			svd.Variables[0].Expression = varExpression.(VariableExpression)
		}

		if varInitialValue != nil {
			svd.InitialValue = varInitialValue.(VariableInitialValue)
		}

		contractDef.SubNodes = append(contractDef.SubNodes, svd)
	}
}

func processCallExpressionInitializer(varInitializer interface{}, svd *StateVariableDeclaration) {

	var varInitialValue interface{} = nil
	var varExpression interface{} = nil

	var varInitializerCallee, _ = varInitializer.(map[string]interface{})["callee"]
	var varInitializerCalleeName, _ = varInitializerCallee.(map[string]interface{})["name"]
	var varInitializerArguments, _ = varInitializer.(map[string]interface{})["arguments"].([]interface{})
	var varInitializerFirstArgument = varInitializerArguments[0].(interface{})

	// handling 'complicated' types
	// (u8, u16, ..., u256, i8, i16, ..., i256, address)
	if varInitializerCalleeName == "u8" {
		var varInitializerFirstArgumentValue float64 = varInitializerFirstArgument.(map[string]interface{})["value"].(float64)
		var number = strconv.FormatFloat(varInitializerFirstArgumentValue, 'f', -1, 64)
		svd.Variables[0].TypeName.Name = "uint8"
		varInitialValue = VariableInitialValue{
			Type:            "NumberLiteral",
			Number:          number,
			Subdenomination: nil,
		}

		varExpression = VariableExpression{
			Type:            "NumberLiteral",
			Number:          number,
			Subdenomination: nil,
		}

	} else if varInitializerCalleeName == "address" {

		var varInitializerFirstArgumentValue = varInitializerFirstArgument.(map[string]interface{})["value"]
		var numberStr = varInitializerFirstArgumentValue.(string)
		svd.Variables[0].TypeName.Name = "address"

		varInitialValue = VariableInitialValue{
			Type:            "NumberLiteral",
			Number:          numberStr,
			Subdenomination: nil,
		}

		varExpression = VariableExpression{
			Type:            "NumberLiteral",
			Number:          numberStr,
			Subdenomination: nil,
		}
	}

	if varExpression != nil {
		svd.Variables[0].Expression = varExpression.(VariableExpression)
	}

	if varInitialValue != nil {
		svd.InitialValue = varInitialValue.(VariableInitialValue)
	}
}

func processSimpleInitializer(varInitializer interface{}, svd *StateVariableDeclaration) {

	var varInitialValue interface{} = nil
	var varExpression interface{} = nil

	varInitializerType := varInitializer.(map[string]interface{})["type"].(string)
	varInitializerValue := varInitializer.(map[string]interface{})["value"]

	if varInitializerType == "StringLiteral" {

		svd.Variables[0].TypeName.Name = "string"

		varInitialValue = VariableInitialValue{
			Type:      varInitializerType,
			Value:     varInitializerValue.(string),
			Parts:     []string{varInitializerValue.(string)},
			IsUnicode: []bool{false},
		}
		varExpression = VariableExpression{
			Type:      "StringLiteral",
			Value:     varInitializerValue.(string),
			Parts:     []string{varInitializerValue.(string)},
			IsUnicode: []bool{false},
		}
	} else if varInitializerType == "BooleanLiteral" {

		svd.Variables[0].TypeName.Name = "bool"

		varInitialValue = VariableInitialValue{
			Type:  varInitializerType,
			Value: varInitializerValue.(bool),
		}
		varExpression = VariableExpression{
			Type:  varInitializerType,
			Value: varInitializerValue.(bool),
		}
	} else {
		panic(fmt.Sprintf("unrecognized varInitializerType %v", varInitializerType))
	}

	if varExpression != nil {
		svd.Variables[0].Expression = varExpression.(VariableExpression)
	}

	if varInitialValue != nil {
		svd.InitialValue = varInitialValue.(VariableInitialValue)
	}
}

// processFunctionDeclaration processes a function declaration and adds it to the contract definition.
func processFunctionDeclaration(contractDef *ContractDefinition, node map[string]interface{}) {
	// Implementation for processing function declarations
	// This function should extract information from the node and add a corresponding
	// function definition to the contractDef.SubNodes
	println("@processFunctionDeclaration", node)
}

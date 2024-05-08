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
	Value           interface{} `json:"value"`
	Number          string      `json:"number"`
	Parts           []string    `json:"parts"`
	IsUnicode       []bool      `json:"isUnicode"`
	Subdenomination interface{} `json:"subdenomination"`
}

type VariableInitialValue struct {
	Type            string      `json:"type"`
	Value           interface{} `json:"value"`
	Number          string      `json:"number"`
	Parts           []string    `json:"parts"`
	IsUnicode       []bool      `json:"isUnicode"`
	Subdenomination interface{} `json:"subdenomination"`
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
	// This function should extract information from the node and add a corresponding
	// state variable declaration to the contractDef.SubNodes

	varDeclarations, _ := node["declarations"].([]interface{})

	for _, decl := range varDeclarations {

		var varIdName, _ = decl.(map[string]interface{})["id"].(map[string]interface{})["name"].(string)
		var varInitializer, _ = decl.(map[string]interface{})["initializer"]
		// var varInitializerValue interface{}
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
		//

		if varInitializer != nil {

			var varInitializerValue interface{}
			varInitializerType, _ = decl.(map[string]interface{})["initializer"].(map[string]interface{})["type"].(string)

			if varInitializerType == "CallExpression" {
				println("@CallExpression varIdName", varIdName)
				// var varInitializerCallee = varInitializerValue.(map[string]interface{})["callee"].(map[string]interface{})["name"].(string)
				var varInitializerCallee, _ = varInitializer.(map[string]interface{})["callee"]
				var varInitializerCalleeName, _ = varInitializerCallee.(map[string]interface{})["name"]

				var varInitializerArguments, _ = varInitializer.(map[string]interface{})["arguments"].([]interface{})

				var varInitializerFirstArgument = varInitializerArguments[0].(interface{})

				// handling 'complicated' types
				// (u8, u16, ..., u256, i8, i16, ..., i256, address)
				if varInitializerCalleeName == "u8" {
					var varInitializerFirstArgumentValue float64 = varInitializerFirstArgument.(map[string]interface{})["value"].(float64)
					println("@CallExpression varInitializerFirstArgumentValue", varInitializerFirstArgumentValue)

					svd.Variables[0].TypeName.Name = "uint8"

					varInitialValue = VariableInitialValue{
						Type:            "NumberLiteral",
						Number:          strconv.FormatFloat(varInitializerFirstArgumentValue, 'f', -1, 64),
						Subdenomination: nil,
					}

					varExpression = VariableExpression{
						Type:            "NumberLiteral",
						Number:          strconv.FormatFloat(varInitializerFirstArgumentValue, 'f', -1, 64),
						Subdenomination: nil,
					}

				} else if varInitializerCalleeName == "address" {

					var varInitializerFirstArgumentValue = varInitializerFirstArgument.(map[string]interface{})["value"]
					svd.Variables[0].TypeName.Name = "address"

					varInitialValue = VariableInitialValue{
						Type:            "NumberLiteral",
						Number:          varInitializerFirstArgumentValue.(string),
						Subdenomination: nil,
					}

					varExpression = VariableExpression{
						Type:            "NumberLiteral",
						Number:          varInitializerFirstArgumentValue.(string),
						Subdenomination: nil,
					}
				}

			} else { // here variable's types are accessible directly
				varInitializerValue = decl.(map[string]interface{})["initializer"].(map[string]interface{})["value"]
				println("@varInitializerValue", varInitializerValue)

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

			}
		}

		// println("@processVariableStatement, varDeclarations type", varType, varIdName, varInitializerValue, varInitializerType)

		if varExpression != nil {
			svd.Variables[0].Expression = varExpression.(VariableExpression)
		}

		if varInitialValue != nil {
			svd.InitialValue = varInitialValue.(VariableInitialValue)
		}

		contractDef.SubNodes = append(contractDef.SubNodes, svd)
	}
}

// processFunctionDeclaration processes a function declaration and adds it to the contract definition.
func processFunctionDeclaration(contractDef *ContractDefinition, node map[string]interface{}) {
	// Implementation for processing function declarations
	// This function should extract information from the node and add a corresponding
	// function definition to the contractDef.SubNodes
	println("@processFunctionDeclaration", node)
}

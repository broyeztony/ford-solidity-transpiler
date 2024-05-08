package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	Type      string   `json:"type"`
	Value     string   `json:"value"`
	Parts     []string `json:"parts"`
	IsUnicode []bool   `json:"isUnicode"`
}

type VariableInitialValue struct {
	Type      string   `json:"type"`
	Value     string   `json:"value"`
	Parts     []string `json:"parts"`
	IsUnicode []bool   `json:"isUnicode"`
}

func main() {

	filePath := "data/helloworld.ast.json"

	inputAST, err := loadAndUnmarshalJSON(filePath)
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
		var varInitializerValue interface{}
		var varInitializerType string
		var varInitialValue interface{} = nil
		var varExpression interface{} = nil

		if varInitializer != nil {
			varInitializerValue, _ = decl.(map[string]interface{})["initializer"].(map[string]interface{})["value"].(string)
			varInitializerType, _ = decl.(map[string]interface{})["initializer"].(map[string]interface{})["type"].(string)

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
		}

		// println("@processVariableStatement, varDeclarations type", varType, varIdName, varInitializerValue, varInitializerType)

		svd := StateVariableDeclaration{
			Type: "StateVariableDeclaration",
			Variables: []Variable{{
				Type: "VariableDeclaration",
				TypeName: VariableTypeName{
					Type:            "ElementaryTypeName",
					Name:            "string",
					StateMutability: nil,
				},
				Name: varIdName,
				Identifier: VariableIdentifier{
					Type: "Identifier",
					Name: varIdName,
				},

				Visibility: "public",
				IsStateVar: true,
				// IsDeclaredConst: false,
				// IsIndexed:       false,
				// IsImmutable:     false,
			}},
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

// processFunctionDeclaration processes a function declaration and adds it to the contract definition.
func processFunctionDeclaration(contractDef *ContractDefinition, node map[string]interface{}) {
	// Implementation for processing function declarations
	// This function should extract information from the node and add a corresponding
	// function definition to the contractDef.SubNodes
	println("@processFunctionDeclaration", node)
}

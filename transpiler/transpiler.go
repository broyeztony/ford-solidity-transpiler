package transpiler

import (
	"fmt"
	"strconv"
)

type Transpiler struct{}

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
	Type            string           `json:"type"`
	TypeName        VariableTypeName `json:"typeName"`
	Name            interface{}      `json:"name"`
	Identifier      interface{}      `json:"identifier"`
	Expression      interface{}      `json:"expression,omitempty"`
	Visibility      interface{}      `json:"visibility,omitempty"`
	IsStateVar      bool             `json:"isStateVar"`
	IsDeclaredConst bool             `json:"isDeclaredConst"`
	IsIndexed       bool             `json:"isIndexed"`
	IsImmutable     bool             `json:"isImmutable"`
	Override        interface{}      `json:"override"`
	StorageLocation interface{}      `json:"storageLocation"`
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

func NewTranspiler() *Transpiler {

	transpiler := &Transpiler{}
	return transpiler
}

func (t *Transpiler) Transpile(inputAST ASTNode) ContractDefinition {
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
			t.processVariableStatement(&outputAST, bodyNode)
		case "FunctionDeclaration":
			t.processFunctionDeclaration(&outputAST, bodyNode)
		}
	}

	return outputAST
}

// processVariableStatement processes a variable statement and adds it to the contract definition.
func (t *Transpiler) processVariableStatement(contractDef *ContractDefinition, node map[string]interface{}) {

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
				t.processCallExpressionInitializer(varInitializer, &svd)
			case "StringLiteral", "BooleanLiteral":
				t.processSimpleInitializer(varInitializer, &svd)
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

func (t *Transpiler) processCallExpressionInitializer(varInitializer interface{}, svd *StateVariableDeclaration) {

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

func (t *Transpiler) processSimpleInitializer(varInitializer interface{}, svd *StateVariableDeclaration) {

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

// Define structures for the output AST
type FunctionDefinition struct {
	Type             string        `json:"type"`
	Name             string        `json:"name"`
	Parameters       []Variable    `json:"parameters"`
	ReturnParameters []Variable    `json:"returnParameters"`
	Body             Block         `json:"body"`
	Visibility       string        `json:"visibility"`
	Modifiers        []interface{} `json:"modifiers"` // Assuming modifiers can vary; adjust as needed
	Override         interface{}   `json:"override"`  // Assuming override can vary; adjust as needed
	IsConstructor    bool          `json:"isConstructor"`
	IsReceiveEther   bool          `json:"isReceiveEther"`
	IsFallback       bool          `json:"isFallback"`
	IsVirtual        bool          `json:"isVirtual"`
	StateMutability  string        `json:"stateMutability"`
}

type Block struct {
	Type       string      `json:"type"`
	Statements []Statement `json:"statements"`
}

type Statement struct {
	Type       string     `json:"type"`
	Expression Expression `json:"expression"`
}

type Expression struct {
	Type     string             `json:"type"`
	Operator string             `json:"operator"`
	Left     VariableIdentifier `json:"left"`
	Right    VariableIdentifier `json:"right"`
}

// processFunctionDeclaration processes a function declaration and adds it to the contract definition.
func (t *Transpiler) processFunctionDeclaration(contractDef *ContractDefinition, node map[string]interface{}) {
	println("@processFunctionDeclaration", node)

	funcNameNode, ok := node["name"].(interface{})
	if !ok {
		fmt.Println("Error: 'FunctionDeclaration' is not a valid node.")
		return
	}

	funcNameNodeDecl := funcNameNode.(map[string]interface{})
	funcName, ok := funcNameNodeDecl["name"].(string)

	if !ok {
		fmt.Println("Error: invalid function's name.")
		return
	}

	fdef := FunctionDefinition{
		Type:             "FunctionDefinition",
		Name:             funcName,
		Parameters:       nil,
		ReturnParameters: nil,
		Body:             Block{},
		Visibility:       "public", // TODO: handle all visibility settings
		Modifiers:        nil,
		Override:         nil,
		IsConstructor:    false,
		IsReceiveEther:   false,
		IsFallback:       false,
		IsVirtual:        false,
		StateMutability:  "",
	}

	funcParams, ok := node["params"].(interface{})
	funcParamsDecl := funcParams.([]interface{})

	for _, declInterface := range funcParamsDecl { // processes all function's parameters

		decl, ok := declInterface.(map[string]interface{})
		if !ok {
			fmt.Println("Error: `params` declaration is not a valid map.")
			continue
		}

		// function's parameter's name
		paramName, _ := decl["name"].(string)

		variable := Variable{
			Type: "VariableDeclaration",
			TypeName: VariableTypeName{
				Type:            "ElementaryTypeName",
				Name:            "uint8", // use yaml spec to get the correct type
				StateMutability: nil,
			},
			Name: paramName,
			Identifier: VariableIdentifier{
				Type: "Identifier",
				Name: paramName,
			},
			Expression:      nil,
			IsStateVar:      false,
			IsDeclaredConst: false,
			IsIndexed:       false,
			IsImmutable:     false,
			Override:        nil,
			StorageLocation: nil,
		}

		fdef.Parameters = append(fdef.Parameters, variable)
	}

	// processes function's return statement
	fdef.ReturnParameters = []Variable{{
		Type: "VariableDeclaration",
		TypeName: VariableTypeName{
			Type:            "ElementaryTypeName",
			Name:            "uint8", // use yaml spec to get the correct return type
			StateMutability: nil,
		},
		Name:            nil,
		Identifier:      nil,
		Expression:      nil,
		IsStateVar:      false,
		IsIndexed:       false,
		StorageLocation: nil,
	}}

	contractDef.SubNodes = append(contractDef.SubNodes, fdef)
}

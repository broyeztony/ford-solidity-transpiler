package main

import (
	"encoding/json"
	"fmt"
	transpiler "ford-solidity-transpiler/transpiler"
	"os"
)

func main() {

	inputASTFilePath := "data/primitives.ast.json" // AST
	inputAST, err := loadAndUnmarshalJSON(inputASTFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t := transpiler.NewTranspiler()
	outputAST := t.Transpile(inputAST)
	outputJSON, err := json.MarshalIndent(outputAST, "", "  ")
	if err != nil {
		fmt.Println("Error generating output JSON:", err)
		return
	}

	// println(len(outputJSON))
	fmt.Println(string(outputJSON))
}

func loadAndUnmarshalJSON(filePath string) (transpiler.ASTNode, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading JSON file: %w", err)
	}

	var inputAST transpiler.ASTNode
	if err := json.Unmarshal(jsonData, &inputAST); err != nil {
		return nil, fmt.Errorf("parsing input JSON: %w", err)
	}

	return inputAST, nil
}

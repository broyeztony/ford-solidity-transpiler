package transpiler

import "fmt"

func (t *Transpiler) processReturnStatement(body interface{}) Statement {

	stmtTypeArgument := body.(map[string]interface{})["argument"]
	stmtTypeArgumentType := stmtTypeArgument.(map[string]interface{})["type"].(string)
	fmt.Println("@stmtType", stmtTypeArgumentType)

	var stmt Statement
	switch stmtTypeArgumentType {
	case "BinaryExpression":
		stmtTypeArgumentOperator := stmtTypeArgument.(map[string]interface{})["operator"].(string)
		stmtTypeArgumentLeft := stmtTypeArgument.(map[string]interface{})["left"].(map[string]interface{})
		stmtTypeArgumentRight := stmtTypeArgument.(map[string]interface{})["right"].(map[string]interface{})

		stmt = Statement{
			Type: "ReturnStatement",
			Expression: Expression{
				Type:     stmtTypeArgumentType,
				Operator: stmtTypeArgumentOperator,
				Left: VariableIdentifier{
					Type: stmtTypeArgumentLeft["type"].(string),
					Name: stmtTypeArgumentLeft["name"].(string),
				},
				Right: VariableIdentifier{
					Type: stmtTypeArgumentRight["type"].(string),
					Name: stmtTypeArgumentRight["name"].(string),
				},
			},
		}
	}

	return stmt
}

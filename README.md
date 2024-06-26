**Simply put, this is a [Ford](https://github.com/broyeztony/ford-lang-parser) to Solidity AST-to-AST transpiler.**

It takes a json-formatted AST generated by https://github.com/broyeztony/ford-solidity-parser and transpile it to a Solidity-compliant AST.

Once we have a Solidity-compliant AST, we can either:
 - retro-engineer the original Solidity source code (.sol file) using https://github.com/Consensys/solc-typed-ast/tree/master?tab=readme-ov-file#converting-an-ast-back-to-source and then compile it the usual way.
 - or we could go down the rabbit hole and compile the AST to EVM bytecode 

## Usage

1. Create a Ford Smart Contract and parse it using the [Ford parser](https://github.com/broyeztony/ford-solidity-parser). Example 
```ford
contract PrimitiveTypes;

let aString = ""; // string
let aBool = true; // boolean
let aUint8 = u8(255); // uint8
let anAddress = address("0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c"); // address
```
2. `/ford-solidity-parser main*
   ❯ go run main.go > primitives.ast.json`  

```json
{
  "body": [
    {
      "declarations": [
        {
          "id": {
            "name": "aString",
            "type": "Identifier"
          },
          "initializer": {
            "type": "StringLiteral",
            "value": ""
          },
          "type": "VariableDeclaration"
        }
      ],
      "type": "VariableStatement"
    },
    {
      "declarations": [
        {
          "id": {
            "name": "aBool",
            "type": "Identifier"
          },
          "initializer": {
            "type": "BooleanLiteral",
            "value": true
          },
          "type": "VariableDeclaration"
        }
      ],
      "type": "VariableStatement"
    },
    {
      "declarations": [
        {
          "id": {
            "name": "aUint8",
            "type": "Identifier"
          },
          "initializer": {
            "arguments": [
              {
                "type": "NumericLiteral",
                "value": 255
              }
            ],
            "callee": {
              "name": "u8",
              "type": "Identifier"
            },
            "type": "CallExpression"
          },
          "type": "VariableDeclaration"
        }
      ],
      "type": "VariableStatement"
    },
    {
      "declarations": [
        {
          "id": {
            "name": "anAddress",
            "type": "Identifier"
          },
          "initializer": {
            "arguments": [
              {
                "type": "StringLiteral",
                "value": "0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c"
              }
            ],
            "callee": {
              "name": "address",
              "type": "Identifier"
            },
            "type": "CallExpression"
          },
          "type": "VariableDeclaration"
        }
      ],
      "type": "VariableStatement"
    }
  ],
  "name": "PrimitiveTypes",
  "type": "Contract"
}
```

3. Then, transpile `/ford-solidity-transpiler main ≡
   ❯ go run main.go`

and you get a Solidity-compliant AST.

```json
{
  "type": "ContractDefinition",
  "name": "PrimitiveTypes",
  "baseContracts": [],
  "subNodes": [
    {
      "type": "StateVariableDeclaration",
      "variables": [
        {
          "type": "VariableDeclaration",
          "typeName": {
            "type": "ElementaryTypeName",
            "name": "string",
            "stateMutability": null
          },
          "name": "aString",
          "identifier": {
            "type": "Identifier",
            "name": "aString"
          },
          "expression": {
            "type": "StringLiteral",
            "value": "",
            "parts": [
              ""
            ],
            "isUnicode": [
              false
            ]
          },
          "visibility": "public",
          "isStateVar": true,
          "isDeclaredConst": false,
          "isIndexed": false,
          "isImmutable": false,
          "override": null,
          "storageLocation": null
        }
      ],
      "initialValue": {
        "type": "StringLiteral",
        "value": "",
        "parts": [
          ""
        ],
        "isUnicode": [
          false
        ]
      }
    },
    {
      "type": "StateVariableDeclaration",
      "variables": [
        {
          "type": "VariableDeclaration",
          "typeName": {
            "type": "ElementaryTypeName",
            "name": "bool",
            "stateMutability": null
          },
          "name": "aBool",
          "identifier": {
            "type": "Identifier",
            "name": "aBool"
          },
          "expression": {
            "type": "BooleanLiteral",
            "value": true
          },
          "visibility": "public",
          "isStateVar": true,
          "isDeclaredConst": false,
          "isIndexed": false,
          "isImmutable": false,
          "override": null,
          "storageLocation": null
        }
      ],
      "initialValue": {
        "type": "BooleanLiteral",
        "value": true
      }
    },
    {
      "type": "StateVariableDeclaration",
      "variables": [
        {
          "type": "VariableDeclaration",
          "typeName": {
            "type": "ElementaryTypeName",
            "name": "uint8",
            "stateMutability": null
          },
          "name": "aUint8",
          "identifier": {
            "type": "Identifier",
            "name": "aUint8"
          },
          "expression": {
            "type": "NumberLiteral",
            "number": "255"
          },
          "visibility": "public",
          "isStateVar": true,
          "isDeclaredConst": false,
          "isIndexed": false,
          "isImmutable": false,
          "override": null,
          "storageLocation": null
        }
      ],
      "initialValue": {
        "type": "NumberLiteral",
        "number": "255"
      }
    },
    {
      "type": "StateVariableDeclaration",
      "variables": [
        {
          "type": "VariableDeclaration",
          "typeName": {
            "type": "ElementaryTypeName",
            "name": "address",
            "stateMutability": null
          },
          "name": "anAddress",
          "identifier": {
            "type": "Identifier",
            "name": "anAddress"
          },
          "expression": {
            "type": "NumberLiteral",
            "number": "0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c"
          },
          "visibility": "public",
          "isStateVar": true,
          "isDeclaredConst": false,
          "isIndexed": false,
          "isImmutable": false,
          "override": null,
          "storageLocation": null
        }
      ],
      "initialValue": {
        "type": "NumberLiteral",
        "number": "0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c"
      }
    }
  ],
  "kind": "contract"
}
```

The above AST would then convert back to Solidity source code as
```solidity
pragma solidity ^0.8.4;

contract PrimitiveTypes {
    string public aString = "";
    bool public aBool = true; // boolean
    uint8 public aUint8 = 255; // uint8
    address public anAddress = 0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c; // address
}
```

## Another example (transpiling functions)

1. Given the below Ford Smart Contract
```ford
contract View;

let x = u8(1);

// by default, all defs are `public`
// `y` is of type `u8` (uint8)
// `addToX` return type is `u8` (uint8)
def addToX(y) {
    return x + y;
}
```
and the companion metadata file
```yaml
state:
  - name: x
    type: u8

defs:
  - name: addToX
    parameters:
      y: u8
    returnType: u8
    visibility:
      - public
    stateMutability:
      - view

events: {}
```

2. `/ford-solidity-parser main*
      ❯ go run main.go > functions.ast.json`

```json
{
  "body": [
    {
      "declarations": [
        {
          "id": {
            "name": "x",
            "type": "Identifier"
          },
          "initializer": {
            "arguments": [
              {
                "type": "NumericLiteral",
                "value": 1
              }
            ],
            "callee": {
              "name": "u8",
              "type": "Identifier"
            },
            "type": "CallExpression"
          },
          "type": "VariableDeclaration"
        }
      ],
      "type": "VariableStatement"
    },
    {
      "body": {
        "body": [
          {
            "argument": {
              "left": {
                "name": "x",
                "type": "Identifier"
              },
              "operator": "+",
              "right": {
                "name": "y",
                "type": "Identifier"
              },
              "type": "BinaryExpression"
            },
            "type": "ReturnStatement"
          }
        ],
        "type": "BlockStatement"
      },
      "name": {
        "name": "addToX",
        "type": "Identifier"
      },
      "params": [
        {
          "name": "y",
          "type": "Identifier"
        }
      ],
      "type": "FunctionDeclaration"
    }
  ],
  "name": "ViewAndPure",
  "type": "Contract"
}
```

3. Then, transpile `/ford-solidity-transpiler main ≡
   ❯ go run main.go`

and you get a Solidity-compliant AST.

```json
{
  "type": "ContractDefinition",
  "name": "ViewAndPure",
  "baseContracts": [],
  "subNodes": [
    {
      "type": "StateVariableDeclaration",
      "variables": [
        {
          "type": "VariableDeclaration",
          "typeName": {
            "type": "ElementaryTypeName",
            "name": "uint8",
            "stateMutability": null
          },
          "name": "x",
          "identifier": {
            "type": "Identifier",
            "name": "x"
          },
          "expression": {
            "type": "NumberLiteral",
            "number": "1"
          },
          "visibility": "public",
          "isStateVar": true,
          "isDeclaredConst": false,
          "isIndexed": false,
          "isImmutable": false,
          "override": null,
          "storageLocation": null
        }
      ],
      "initialValue": {
        "type": "NumberLiteral",
        "number": "1"
      }
    },
    {
      "type": "FunctionDefinition",
      "name": "addToX",
      "parameters": [
        {
          "type": "VariableDeclaration",
          "typeName": {
            "type": "ElementaryTypeName",
            "name": "uint8",
            "stateMutability": null
          },
          "name": "y",
          "identifier": {
            "type": "Identifier",
            "name": "y"
          },
          "isStateVar": false,
          "isDeclaredConst": false,
          "isIndexed": false,
          "isImmutable": false,
          "override": null,
          "storageLocation": null
        }
      ],
      "returnParameters": [
        {
          "type": "VariableDeclaration",
          "typeName": {
            "type": "ElementaryTypeName",
            "name": "uint8",
            "stateMutability": null
          },
          "name": null,
          "identifier": null,
          "isStateVar": false,
          "isDeclaredConst": false,
          "isIndexed": false,
          "isImmutable": false,
          "override": null,
          "storageLocation": null
        }
      ],
      "body": {
        "type": "Block",
        "statements": [
          {
            "type": "ReturnStatement",
            "expression": {
              "type": "BinaryExpression",
              "operator": "+",
              "left": {
                "type": "Identifier",
                "name": "x"
              },
              "right": {
                "type": "Identifier",
                "name": "y"
              }
            }
          }
        ]
      },
      "visibility": "public",
      "modifiers": [],
      "override": null,
      "isConstructor": false,
      "isReceiveEther": false,
      "isFallback": false,
      "isVirtual": false,
      "stateMutability": "view"
    }
  ],
  "kind": "contract"
}
```

The above AST would then convert back to Solidity source code as
```solidity
contract ViewAndPure {
  uint8 public x = 1;

  // Promise not to modify the state.
  function addToX(uint8 y) public view returns (uint8) {
    return x + y;
  }
}
```

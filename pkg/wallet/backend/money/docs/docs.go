// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/balance": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get balance",
                "operationId": "2",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Public key prefixed with 0x",
                        "name": "pubkey",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/money.BalanceResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/money.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/block-height": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Money partition's latest block number",
                "operationId": "4",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/money.BlockHeightResponse"
                        }
                    }
                }
            }
        },
        "/list-bills": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "List bills",
                "operationId": "1",
                "parameters": [
                    {
                        "type": "string",
                        "example": "0x000000000000000000000000000000000000000000000000000000000000000123",
                        "description": "Public key prefixed with 0x",
                        "name": "pubkey",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "limits how many bills are returned in response",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "response will include bills starting after offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/money.ListBillsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/money.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/proof": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get proof",
                "operationId": "3",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the bill (hex)",
                        "name": "bill_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_alphabill-org_alphabill_pkg_wallet_backend_money.Bills"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/money.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/money.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "anypb.Any": {
            "type": "object",
            "properties": {
                "type_url": {
                    "description": "A URL/resource name that uniquely identifies the type of the serialized\nprotocol buffer message. This string must contain at least\none \"/\" character. The last segment of the URL's path must represent\nthe fully qualified name of the type (as in\n` + "`" + `path/google.protobuf.Duration` + "`" + `). The name should be in a canonical form\n(e.g., leading \".\" is not accepted).\n\nIn practice, teams usually precompile into the binary all types that they\nexpect it to use in the context of Any. However, for URLs which use the\nscheme ` + "`" + `http` + "`" + `, ` + "`" + `https` + "`" + `, or no scheme, one can optionally set up a type\nserver that maps type URLs to message definitions as follows:\n\n  - If no scheme is provided, ` + "`" + `https` + "`" + ` is assumed.\n  - An HTTP GET on the URL must yield a [google.protobuf.Type][]\n    value in binary format, or produce an error.\n  - Applications are allowed to cache lookup results based on the\n    URL, or have them precompiled into a binary to avoid any\n    lookup. Therefore, binary compatibility needs to be preserved\n    on changes to types. (Use versioned type names to manage\n    breaking changes.)\n\nNote: this functionality is not currently available in the official\nprotobuf release, and it is not used for type URLs beginning with\ntype.googleapis.com.\n\nSchemes other than ` + "`" + `http` + "`" + `, ` + "`" + `https` + "`" + ` (or the empty scheme) might be\nused with implementation specific semantics.",
                    "type": "string"
                },
                "value": {
                    "description": "Must be a valid serialized protocol buffer of the above specified type.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "block.BlockProof": {
            "type": "object",
            "properties": {
                "block_header_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "block_tree_hash_chain": {
                    "$ref": "#/definitions/block.BlockTreeHashChain"
                },
                "hash_value": {
                    "description": "hash value of either primary tx or secondary txs or zero hash, depending on proof type",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "proof_type": {
                    "$ref": "#/definitions/block.ProofType"
                },
                "sec_tree_hash_chain": {
                    "$ref": "#/definitions/block.SecTreeHashChain"
                },
                "transactions_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "unicity_certificate": {
                    "$ref": "#/definitions/certificates.UnicityCertificate"
                }
            }
        },
        "block.BlockTreeHashChain": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/block.ChainItem"
                    }
                }
            }
        },
        "block.ChainItem": {
            "type": "object",
            "properties": {
                "hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "val": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "block.MerklePathItem": {
            "type": "object",
            "properties": {
                "direction_left": {
                    "description": "DirectionLeft direction from parent node; left=true right=false",
                    "type": "boolean"
                },
                "path_item": {
                    "description": "PathItem Hash of Merkle Tree node",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "block.ProofType": {
            "type": "integer",
            "enum": [
                0,
                1,
                2,
                3,
                4
            ],
            "x-enum-varnames": [
                "ProofType_PRIM",
                "ProofType_SEC",
                "ProofType_ONLYSEC",
                "ProofType_NOTRANS",
                "ProofType_EMPTYBLOCK"
            ]
        },
        "block.SecTreeHashChain": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/block.MerklePathItem"
                    }
                }
            }
        },
        "certificates.InputRecord": {
            "type": "object",
            "properties": {
                "block_hash": {
                    "description": "hash of the block",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "hash": {
                    "description": "state hash to be certified",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "previous_hash": {
                    "description": "previously certified state hash",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "round_number": {
                    "description": "transaction system's round number",
                    "type": "integer"
                },
                "summary_value": {
                    "description": "summary value to certified",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "certificates.UnicityCertificate": {
            "type": "object",
            "properties": {
                "input_record": {
                    "$ref": "#/definitions/certificates.InputRecord"
                },
                "unicity_seal": {
                    "$ref": "#/definitions/certificates.UnicitySeal"
                },
                "unicity_tree_certificate": {
                    "$ref": "#/definitions/certificates.UnicityTreeCertificate"
                }
            }
        },
        "certificates.UnicitySeal": {
            "type": "object",
            "properties": {
                "hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "previous_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "root_chain_round_number": {
                    "type": "integer"
                },
                "signatures": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        }
                    }
                }
            }
        },
        "certificates.UnicityTreeCertificate": {
            "type": "object",
            "properties": {
                "sibling_hashes": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        }
                    }
                },
                "system_description_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "system_identifier": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "github_com_alphabill-org_alphabill_pkg_wallet_backend_money.Bill": {
            "type": "object",
            "properties": {
                "OwnerPredicate": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "id": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "isDcBill": {
                    "type": "boolean"
                },
                "orderNumber": {
                    "description": "OrderNumber insertion order of given bill in pubkey =\u003e list of bills bucket, needed for determistic paging",
                    "type": "integer"
                },
                "txHash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "txProof": {
                    "$ref": "#/definitions/money.TxProof"
                },
                "value": {
                    "type": "integer"
                }
            }
        },
        "github_com_alphabill-org_alphabill_pkg_wallet_backend_money.Bills": {
            "type": "object",
            "properties": {
                "bills": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_alphabill-org_alphabill_pkg_wallet_backend_money.Bill"
                    }
                }
            }
        },
        "money.BalanceResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "string",
                    "example": "0"
                }
            }
        },
        "money.BlockHeightResponse": {
            "type": "object",
            "properties": {
                "blockHeight": {
                    "type": "string",
                    "example": "0"
                },
                "lastRoundNumber": {
                    "type": "integer"
                }
            }
        },
        "money.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "money.ListBillVM": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "format": "base64",
                    "example": "AAAAAAgwv3UA1HfGO4qc1T3I3EOvqxfcrhMjJpr9Tn4="
                },
                "isDcBill": {
                    "type": "boolean",
                    "example": false
                },
                "txHash": {
                    "type": "string",
                    "format": "base64",
                    "example": "Q4ShCITC0ODXPR+j1Zl/teYcoU3/mAPy0x8uSsvQFM8="
                },
                "value": {
                    "type": "string",
                    "example": "1000"
                }
            }
        },
        "money.ListBillsResponse": {
            "type": "object",
            "properties": {
                "bills": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/money.ListBillVM"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "money.TxProof": {
            "type": "object",
            "properties": {
                "blockNumber": {
                    "type": "integer"
                },
                "proof": {
                    "$ref": "#/definitions/block.BlockProof"
                },
                "tx": {
                    "$ref": "#/definitions/txsystem.Transaction"
                }
            }
        },
        "txsystem.Transaction": {
            "type": "object",
            "properties": {
                "owner_proof": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "system_id": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "timeout": {
                    "type": "integer"
                },
                "transaction_attributes": {
                    "$ref": "#/definitions/anypb.Any"
                },
                "unit_id": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Money Partition Indexing Backend API",
	Description:      "This service processes blocks from the Money partition and indexes ownership of bills.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
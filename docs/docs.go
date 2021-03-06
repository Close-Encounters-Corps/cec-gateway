// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/0/login/discord": {
            "get": {
                "tags": [
                    "private"
                ],
                "summary": "Authenticate using Discord",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Second phase: State to fetch from CEC Auth",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "URL to redirect after second phase",
                        "name": "redirect_url",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AuthPhaseResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/0/users/current": {
            "get": {
                "tags": [
                    "private"
                ],
                "summary": "Get current user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization token",
                        "name": "X-Auth-Token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AuthPhaseResult": {
            "type": "object",
            "properties": {
                "next_url": {
                    "description": "next url",
                    "type": "string"
                },
                "phase": {
                    "description": "phase",
                    "type": "integer"
                },
                "token": {
                    "description": "token",
                    "type": "string"
                },
                "user": {
                    "description": "user",
                    "$ref": "#/definitions/models.User"
                }
            }
        },
        "models.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "description": "message",
                    "type": "string"
                },
                "request_id": {
                    "description": "request id",
                    "type": "string"
                }
            }
        },
        "models.Principal": {
            "type": "object",
            "properties": {
                "admin": {
                    "description": "admin",
                    "type": "boolean"
                },
                "created_on": {
                    "description": "created on",
                    "type": "string"
                },
                "id": {
                    "description": "id",
                    "type": "integer"
                },
                "last_login": {
                    "description": "last login",
                    "type": "string"
                },
                "state": {
                    "description": "state",
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "id",
                    "type": "integer"
                },
                "principal": {
                    "description": "principal",
                    "$ref": "#/definitions/models.Principal"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "CEC Gateway",
	Description:      "Gateway endpoint of a CEC Platform v2, serves as a proxy and one swagger to rule them all. Find more at Close Encounters Corps Discord server!",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

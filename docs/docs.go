// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Berke Emrecan Arslan",
            "url": "https://beremaran.com",
            "email": "berke.emrecan.arslan@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/_actuator/bulk-generate": {
            "get": {
                "description": "Generate users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actuator"
                ],
                "summary": "Generate users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "how many users to generate",
                        "name": "n",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "500": {}
                }
            }
        },
        "/_actuator/flush-all": {
            "delete": {
                "description": "Remove all data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actuator"
                ],
                "summary": "Flush Redis Cache",
                "responses": {
                    "200": {},
                    "500": {}
                }
            }
        },
        "/_actuator/user-count": {
            "get": {
                "description": "Get total number of users",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "actuator"
                ],
                "summary": "Get total number of users",
                "responses": {
                    "200": {},
                    "500": {}
                }
            }
        },
        "/leaderboard": {
            "get": {
                "description": "Get leaderboard",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "leaderboard"
                ],
                "summary": "Get leaderboard",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of records in a page",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of records in a page",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.LeaderboardRow"
                            }
                        }
                    },
                    "500": {}
                }
            }
        },
        "/leaderboard/{country_iso_code}": {
            "get": {
                "description": "Get leaderboard",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "leaderboard"
                ],
                "summary": "Get leaderboard",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of records in a page",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of records in a page",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "ISO standard country code",
                        "name": "country_iso_code",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.LeaderboardRow"
                            }
                        }
                    },
                    "500": {}
                }
            }
        },
        "/score/submit": {
            "post": {
                "description": "submit a new score",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "leaderboard",
                    "score"
                ],
                "summary": "submit a new score",
                "parameters": [
                    {
                        "description": "score submission",
                        "name": "score",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.ScoreSubmission"
                        }
                    }
                ],
                "responses": {
                    "200": {},
                    "500": {}
                }
            }
        },
        "/user/create": {
            "post": {
                "description": "Create a new user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "user info",
                        "name": "profile",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UserProfile"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.UserProfile"
                            }
                        }
                    },
                    "500": {}
                }
            }
        },
        "/user/profile/{id}": {
            "get": {
                "description": "Get user details by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user details by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user GUID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.UserProfile"
                            }
                        }
                    },
                    "500": {}
                }
            }
        }
    },
    "definitions": {
        "api.LeaderboardRow": {
            "type": "object",
            "properties": {
                "country": {
                    "type": "string"
                },
                "display_name": {
                    "type": "string"
                },
                "points": {
                    "type": "integer"
                },
                "rank": {
                    "type": "integer"
                }
            }
        },
        "api.ScoreSubmission": {
            "type": "object",
            "required": [
                "score",
                "timestamp",
                "user_id"
            ],
            "properties": {
                "score": {
                    "type": "number"
                },
                "timestamp": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "api.UserProfile": {
            "type": "object",
            "required": [
                "country",
                "display_name"
            ],
            "properties": {
                "country": {
                    "type": "string"
                },
                "display_name": {
                    "type": "string"
                },
                "points": {
                    "type": "number"
                },
                "rank": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.0.4",
	Host:        "leaderboard-v2-lb-ecs-tg-584908050.eu-central-1.elb.amazonaws.com",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Leaderboard Service",
	Description: "Simple & fast leaderboard service",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}

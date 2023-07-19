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
        "/ping": {
            "get": {
                "tags": [
                    "首頁"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/post/createPost": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "貼文資料"
                ],
                "summary": "創建貼文",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用戶名",
                        "name": "author",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "標題",
                        "name": "title",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "內容",
                        "name": "content",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "照片",
                        "name": "image",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code\",\"message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/post/getPostList": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "貼文資料"
                ],
                "summary": "貼文列表",
                "responses": {
                    "200": {
                        "description": "code\",\"message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/createUser": {
            "post": {
                "tags": [
                    "用戶資料"
                ],
                "summary": "新增用戶",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用戶名",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "密碼",
                        "name": "password",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "確認密碼",
                        "name": "repassword",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code\",\"message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/deleteUser": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "用戶資料"
                ],
                "summary": "刪除用戶",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code\",\"message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/findUserByNameAndPwd": {
            "post": {
                "tags": [
                    "用戶資料"
                ],
                "summary": "用戶登入",
                "parameters": [
                    {
                        "description": "用戶名和密碼",
                        "name": "UserInput",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UserInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code\",\"message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/getUserList": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "用戶資料"
                ],
                "responses": {
                    "200": {
                        "description": "code\",\"message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/updateUser": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "用戶資料"
                ],
                "summary": "修改用戶",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "用戶名",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "密碼",
                        "name": "password",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "phone",
                        "name": "phone",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "email",
                        "name": "email",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code\",\"message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "service.UserInput": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Description for what is this security definition being used",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

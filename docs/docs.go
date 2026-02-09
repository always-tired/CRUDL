package docs

import "github.com/swaggo/swag"

const docTemplate = `{
"swagger": "2.0",
"info": {
  "description": "REST API for subscription aggregation service",
  "title": "Subscription Aggregator API",
  "version": "1.0.0"
},
"basePath": "/",
"paths": {
  "/health": {
    "get": {
      "summary": "Health check",
      "responses": {
        "200": {"description": "OK"}
      }
    }
  },
  "/subscriptions": {
    "post": {
      "summary": "Create subscription",
      "parameters": [
        {"in": "body", "name": "subscription", "required": true, "schema": {"$ref": "#/definitions/SubscriptionRequest"}}
      ],
      "responses": {
        "201": {"description": "Created", "schema": {"$ref": "#/definitions/Subscription"}},
        "400": {"description": "Bad request", "schema": {"$ref": "#/definitions/Error"}}
      }
    },
    "get": {
      "summary": "List subscriptions",
      "parameters": [
        {"in": "query", "name": "user_id", "type": "string", "format": "uuid"},
        {"in": "query", "name": "service_name", "type": "string"},
        {"in": "query", "name": "limit", "type": "integer"},
        {"in": "query", "name": "offset", "type": "integer"}
      ],
      "responses": {
        "200": {"description": "OK", "schema": {"type": "array", "items": {"$ref": "#/definitions/Subscription"}}}
      }
    }
  },
  "/subscriptions/{id}": {
    "get": {
      "summary": "Get subscription by id",
      "parameters": [
        {"in": "path", "name": "id", "required": true, "type": "string", "format": "uuid"}
      ],
      "responses": {
        "200": {"description": "OK", "schema": {"$ref": "#/definitions/Subscription"}},
        "404": {"description": "Not found", "schema": {"$ref": "#/definitions/Error"}}
      }
    },
    "put": {
      "summary": "Update subscription",
      "parameters": [
        {"in": "path", "name": "id", "required": true, "type": "string", "format": "uuid"},
        {"in": "body", "name": "subscription", "required": true, "schema": {"$ref": "#/definitions/SubscriptionRequest"}}
      ],
      "responses": {
        "200": {"description": "OK", "schema": {"$ref": "#/definitions/Subscription"}},
        "400": {"description": "Bad request", "schema": {"$ref": "#/definitions/Error"}},
        "404": {"description": "Not found", "schema": {"$ref": "#/definitions/Error"}}
      }
    },
    "delete": {
      "summary": "Delete subscription",
      "parameters": [
        {"in": "path", "name": "id", "required": true, "type": "string", "format": "uuid"}
      ],
      "responses": {
        "204": {"description": "No Content"},
        "404": {"description": "Not found", "schema": {"$ref": "#/definitions/Error"}}
      }
    }
  },
  "/subscriptions/summary": {
    "get": {
      "summary": "Get total cost for period",
      "parameters": [
        {"in": "query", "name": "start", "required": true, "type": "string", "example": "07-2025"},
        {"in": "query", "name": "end", "required": true, "type": "string", "example": "12-2025"},
        {"in": "query", "name": "user_id", "type": "string", "format": "uuid"},
        {"in": "query", "name": "service_name", "type": "string"}
      ],
      "responses": {
        "200": {
          "description": "OK",
          "schema": {
            "type": "object",
            "properties": {"total": {"type": "integer"}}
          }
        },
        "400": {"description": "Bad request", "schema": {"$ref": "#/definitions/Error"}}
      }
    }
  }
},
"definitions": {
  "SubscriptionRequest": {
    "type": "object",
    "required": ["service_name", "price", "user_id", "start_date"],
    "properties": {
      "service_name": {"type": "string"},
      "price": {"type": "integer"},
      "user_id": {"type": "string", "format": "uuid"},
      "start_date": {"type": "string", "example": "07-2025"},
      "end_date": {"type": "string", "example": "12-2025"}
    }
  },
  "Subscription": {
    "type": "object",
    "properties": {
      "id": {"type": "string", "format": "uuid"},
      "service_name": {"type": "string"},
      "price": {"type": "integer"},
      "user_id": {"type": "string", "format": "uuid"},
      "start_date": {"type": "string"},
      "end_date": {"type": "string"},
      "created_at": {"type": "string", "format": "date-time"},
      "updated_at": {"type": "string", "format": "date-time"}
    }
  },
  "Error": {
    "type": "object",
    "properties": {"error": {"type": "string"}}
  }
}
}`

func init() {
	swag.Register(swag.Name, &s{})
}

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}

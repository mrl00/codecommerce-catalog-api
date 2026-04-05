# Create a category
POST http://localhost:8080/api/categories
Content-Type: application/json
{ "name": "Electronics" }

HTTP 201
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" exists

[Captures]
category_id: jsonpath "$.id"

# Create a second category
POST http://localhost:8080/api/categories
Content-Type: application/json
{ "name": "Books" }

HTTP 201

# List categories (paginated)
GET http://localhost:8080/api/categories?page=1&per_page=10

HTTP 200
[Asserts]
jsonpath "$.items" exists
jsonpath "$.total" exists
jsonpath "$.page" == 1
jsonpath "$.per_page" == 10

# Get category by ID
GET http://localhost:8080/api/categories/{{category_id}}

HTTP 200
[Asserts]
jsonpath "$.name" == "Electronics"

# Update category
PUT http://localhost:8080/api/categories/{{category_id}}
Content-Type: application/json
{ "name": "Electronics Updated" }

HTTP 200
[Asserts]
jsonpath "$.name" == "Electronics Updated"

# Get updated category
GET http://localhost:8080/api/categories/{{category_id}}

HTTP 200
[Asserts]
jsonpath "$.name" == "Electronics Updated"

# Delete category
DELETE http://localhost:8080/api/categories/{{category_id}}

HTTP 204

# Get deleted category should return 404
GET http://localhost:8080/api/categories/{{category_id}}

HTTP 404
[Asserts]
jsonpath "$.error" == "category not found"

# Create category with empty name should fail
POST http://localhost:8080/api/categories
Content-Type: application/json
{ "name": "" }

HTTP 400
[Asserts]
jsonpath "$.error" == "name is required"

# Invalid UUID for category ID should fail
GET http://localhost:8080/api/categories/not-a-uuid

HTTP 400
[Asserts]
jsonpath "$.error" == "invalid category id"
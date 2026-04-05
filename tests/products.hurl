# Create a category first (needed for products)
POST http://localhost:8080/api/categories
Content-Type: application/json
{ "name": "Tech Gadgets" }

HTTP 201
[Asserts]
jsonpath "$.id" exists

[Captures]
category_id: jsonpath "$.id"

# Create a product
POST http://localhost:8080/api/products
Content-Type: application/json
{
    "name": "Wireless Mouse",
    "description": "Ergonomic wireless mouse",
    "price": 2999,
    "image_url": "http://example.com/mouse.jpg",
    "category": "{{category_id}}"
}

HTTP 201
[Asserts]
jsonpath "$.name" == "Wireless Mouse"
jsonpath "$.price" == 2999
jsonpath "$.id" exists

[Captures]
product_id: jsonpath "$.id"

# Create a second product in same category
POST http://localhost:8080/api/products
Content-Type: application/json
{
    "name": "USB Keyboard",
    "description": "Mechanical keyboard",
    "price": 7999,
    "category": "{{category_id}}"
}

HTTP 201
[Asserts]
jsonpath "$.name" == "USB Keyboard"

# List products (paginated)
GET http://localhost:8080/api/products?page=1&per_page=10

HTTP 200
[Asserts]
jsonpath "$.items" exists
jsonpath "$.page" == 1
jsonpath "$.per_page" == 10

# Get product by ID
GET http://localhost:8080/api/products/{{product_id}}

HTTP 200
[Asserts]
jsonpath "$.name" == "Wireless Mouse"
jsonpath "$.price" == 2999

# Update product
PUT http://localhost:8080/api/products/{{product_id}}
Content-Type: application/json
{
    "name": "Wireless Mouse V2",
    "description": "Updated ergonomic mouse",
    "price": 3499,
    "image_url": "http://example.com/mouse-v2.jpg",
    "category": "{{category_id}}"
}

HTTP 200
[Asserts]
jsonpath "$.name" == "Wireless Mouse V2"
jsonpath "$.price" == 3499

# List products by category ID
GET http://localhost:8080/api/categories/{{category_id}}/products?page=1&per_page=10

HTTP 200
[Asserts]
jsonpath "$.items" exists

# Delete product
DELETE http://localhost:8080/api/products/{{product_id}}

HTTP 204

# Get deleted product should return 404
GET http://localhost:8080/api/products/{{product_id}}

HTTP 404
[Asserts]
jsonpath "$.error" == "product not found"

# Create product with empty name should fail
POST http://localhost:8080/api/products
Content-Type: application/json
{ "name": "", "price": 1000, "category": "{{category_id}}" }

HTTP 400
[Asserts]
jsonpath "$.error" == "name is required"

# Create product with zero price should fail
POST http://localhost:8080/api/products
Content-Type: application/json
{ "name": "Free Item", "price": 0, "category": "{{category_id}}" }

HTTP 400
[Asserts]
jsonpath "$.error" == "price must be greater than zero"

# Create product without category should fail
POST http://localhost:8080/api/products
Content-Type: application/json
{ "name": "NoCategory", "price": 1000 }

HTTP 400
[Asserts]
jsonpath "$.error" == "category is required"

# Invalid UUID for product ID should fail
GET http://localhost:8080/api/products/not-a-uuid

HTTP 400
[Asserts]
jsonpath "$.error" == "invalid product id"
Financial Management API

This API provides functionalities for user registration, login, password management, and tracking of financial transactions, including creating, updating, deleting, and summarizing them. It also allows users to manage categories for better transaction organization.

Endpoints

1. User Registration

	•	Endpoint: /v1/registration
	•	Method: POST
	•	Description: Register a new user.
	•	Request Body:

{
    "name": "string",
    "username": "string",
    "email": "string",
    "phone_number": "string"
    "password": "string"
}

	•	Response:
	•	200 OK: User is successfully registered.
	•	400 Bad Request: Invalid input data.

2. User Login

	•	Endpoint: /v1/login
	•	Method: POST
	•	Description: Authenticate the user and return a JWT token.
	•	Request Body:

{
    "username": "string",
    "password": "string"
}


	•	Response:
	•	200 OK: Returns a JWT token.
	•	401 Unauthorized: Invalid credentials.

3. Change Password

	•	Endpoint: /v1/change-password
	•	Method: POST
	•	Description: Change the user’s password.
	•	Request Body:

{
    "email": "string",
    "old_password": "string",
    "new_password": "string"
}


	•	Response:
	•	200 OK: Password updated successfully.
	•	400 Bad Request: Invalid request or password.

4. Create Transaction (JWT Protected)

	•	Endpoint: /v1/transaction
	•	Method: POST
	•	Description: Create a new financial transaction.
	•	Request Body:

{
    "amount": "number",
    "category_id": "int",
    "description": "string",
    "category_type": "IN/OUT"
}


	•	Response:
	•	201 Created: Transaction created.
	•	400 Bad Request: Invalid input.

5. Get Transaction Summary (JWT Protected)

	•	Endpoint: /v1/summary-transaction
	•	Method: GET
	•	Description: Get a summary of transactions including total income and expenses.
	•	Response:
	•	200 OK: Returns transaction summary.
	•	401 Unauthorized: Missing or invalid JWT token.

6. Get Transaction List (JWT Protected)

	•	Endpoint: /v1/transactions
	•	Method: GET
	•	Description: Fetch a list of transactions based on filters like category type, date range, offset, and limit.
	•	Request Body:

{
    "category_type": "",
    "start_date": "2024-10-01",
    "end_date": "2024-10-31",
    "offset": 0,
    "limit": 0
}

	•	category_type: Filter transactions by category (e.g., IN, OUT, or leave empty for all).
	•	start_date: The start date for filtering transactions (format: YYYY-MM-DD).
	•	end_date: The end date for filtering transactions (format: YYYY-MM-DD).
	•	offset: The number of items to skip (for pagination).
	•	limit: The maximum number of items to return (for pagination).

	•	Response:
	•	200 OK: Returns the filtered list of transactions.
	•	401 Unauthorized: Missing or invalid JWT token.

7. Create Category (JWT Protected)

	•	Endpoint: /v1/category
	•	Method: POST
	•	Description: Create a new category for organizing transactions.
	•	Request Body:

{
    "category_name": "string",
    "category_type": "IN/OUT"
    "category_description": "string"
}


	•	Response:
	•	201 Created: Category created successfully.
	•	400 Bad Request: Invalid input.

8. Get Categories (JWT Protected)

	•	Endpoint: /v1/categories?type=OUT
	•	Method: GET
	•	Description: Retrieve a list of categories filtered by type.
	•	Query Parameter:
	•	type: Filter categories by type (IN, OUT).
Example:

/v1/categories?type=OUT

	•	Response:
	•	200 OK: Returns the list of categories based on the specified type.
	•	401 Unauthorized: Missing or invalid JWT token.

9. Update Category (JWT Protected)

	•	Endpoint: /v1/update-category
	•	Method: POST
	•	Description: Update an existing category.
	•	Request Body:

{
    "id": "int",
    "category_name": "string",
    "category_description": "string"
}


	•	Response:
	•	200 OK: Category updated.
	•	400 Bad Request: Invalid request.

10. Update Transaction (JWT Protected)

	•	Endpoint: /v1/update-transaction
	•	Method: POST
	•	Description: Update an existing transaction.
	•	Request Body:

{
    "id": "int",
    "amount": "number",
    "description": "string"
}


	•	Response:
	•	200 OK: Transaction updated successfully.
	•	400 Bad Request: Invalid input.

11. Delete Transaction (JWT Protected)

	•	Endpoint: /v1/delete-transaction
	•	Method: POST
	•	Description: Delete a transaction.
	•	Request Body:

{
    "id": "int"
}


	•	Response:
	•	200 OK: Transaction deleted.
	•	400 Bad Request: Invalid input.

12. Delete Category (JWT Protected)

	•	Endpoint: /v1/delete-category
	•	Method: POST
	•	Description: Delete a transaction.
	•	Request Body:

{
    "id": "int"
}


	•	Response:
	•	200 OK: Transaction deleted.
	•	400 Bad Request: Invalid input.

Authentication

This API uses JWT (JSON Web Token) for authentication on protected routes. To access protected routes, you need to include a valid JWT token in the Authorization header:

Authorization: Bearer <token>

Error Responses

	•	400 Bad Request: Invalid or missing parameters.
	•	401 Unauthorized: Missing or invalid JWT token.
	•	500 Internal Server Error: An unexpected error occurred.

How to Run

	1.	Install Go and set up the environment.
	2.	Run the server using:

go run main.go

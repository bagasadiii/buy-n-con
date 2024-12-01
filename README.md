# Buy N Con

**Buy N Con** is a platform that allows users to socialize and sell items they own. The project is currently under development, and more features are planned to be added in the future.
This is my project to showcase my skill

## Features (Current State)
- REST API with middleware for secure and efficient processing.
- Basic CRUD (Create, Read, Update, Delete) operations for user items.

## Technologies Used
- **PostgreSQL Driver**: [PGX](https://github.com/jackc/pgx) for database interaction.
- **Router**: [httprouter](https://github.com/julienschmidt/httprouter), a lightweight and high-performance HTTP router.
- Docker: For Virtualization

## Getting Started
### Prerequisites
- Go (Golang) installed on your system.
- PostgreSQL database set up and running.
  
### Setup Instructions
1. Clone the repository:
   ```bash
   git clone https://github.com/bagasadiii/buy-n-con.git
   cd buy-n-con
2. Run with docker:
   ```bash
   docker compose up --build
3. Check if its working by accessing http://localhost:8080
4. If status is 404, its already working

# API Documentation for Buy-n-Con

This documentation provides an overview of the API endpoints for the **Buy-n-Con** application. The API allows users to register, login, manage items, and posts.

## Base URL
The API base URL is:  
`http://localhost:8080/`

## Authentication
Some routes require authentication. Use the `Authorization` header with the format `Bearer <token>`.

---

## User Endpoints

### 1. **Register User**
- **POST** `/api/register`
- Registers a new user.
- **Request Body**:
    ```json
    {
      "username": "string",
      "email": "string",
      "password": "string"
    }
    ```
- **Response**:
    ```json
    {
      "status": 201,
      "message": "user created",
      "data": {
        "id": "user_id",
        "username": "username",
        "email": "email"
      }
    }
    ```

### 2. **Login User**
- **POST** `/api/login`
- Logs in a user and returns a JWT token.
- **Request Body**:
    ```json
    {
      "username": "string",
      "password": "string"
    }
    ```
- **Response**:
    ```json
    {
      "status": 200,
      "message": "Login successful",
      "data": {
        "token": "JWT_TOKEN"
      }
    }
    ```

### 3. **Get User by Username**
- **GET** `/api/u/:username`
- Retrieves a user's profile by their username.
- **Response**:
    ```json
    {
      "status": 200,
      "message": "OK",
      "data": {
        "id": "user_id",
        "username": "username",
        "email": "email"
      }
    }
    ```

---

## Item Endpoints (Requires Authentication)

### 1. **Create Item**
- **POST** `/api/u/:username/items`
- Creates a new item for the user.
- **Request Body**:
    ```json
    {
      "name": "string",
      "description": "string",
      "price": "number",
      "category": "string"
    }
    ```
- **Response**:
    ```json
    {
      "status": 201,
      "message": "item created",
      "data": {
        "item_id": "item_id",
        "name": "name",
        "description": "description",
        "price": "price",
        "category": "category"
      }
    }
    ```

### 2. **Get Item by ID**
- **GET** `/api/u/:username/items/:item_id`
- Retrieves a specific item by its ID.
- **Response**:
    ```json
    {
      "status": 200,
      "message": "OK",
      "data": {
        "item_id": "item_id",
        "name": "name",
        "description": "description",
        "price": "price",
        "category": "category"
      }
    }
    ```

### 3. **Get All Items**
- **GET** `/api/u/:username/items`
- Retrieves all items for a specific user.
- **Query Parameters**:
    - `limit`: (Optional) Number of items to retrieve (default is 10).
    - `offset`: (Optional) Page offset (default is 0).
- **Response**:
    ```json
    {
      "status": 200,
      "message": "Items fetched",
      "data": [
        {
          "item_id": "item_id",
          "name": "name",
          "description": "description",
          "price": "price",
          "category": "category"
        }
      ]
    }
    ```

### 4. **Update Item**
- **PATCH** `/api/u/:username/items/:item_id`
- Updates an existing item.
- **Request Body**:
    ```json
    {
      "name": "updated name",
      "description": "updated description",
      "price": "updated price",
      "category": "updated category"
    }
    ```
- **Response**:
    ```json
    {
      "status": 200,
      "message": "Item updated",
      "data": {
        "item_id": "item_id",
        "name": "name",
        "description": "description",
        "price": "price",
        "category": "category"
      }
    }
    ```

### 5. **Delete Item**
- **DELETE** `/api/u/:username/items/:item_id`
- Deletes an item by its ID.
- **Response**:
    ```json
    {
      "status": 200,
      "message": "Item deleted successfully",
      "data": null
    }
    ```

---

## Post Endpoints (Requires Authentication)

### 1. **Create Post**
- **POST** `/api/u/:username/post`
- Creates a new post.
- **Request Body**:
    ```json
    {
      "title": "string",
      "content": "string"
    }
    ```
- **Response**:
    ```json
    {
      "status": 201,
      "message": "post created",
      "data": {
        "post_id": "post_id",
        "title": "title",
        "content": "content"
      }
    }
    ```

### 2. **Get Post by ID**
- **GET** `/api/u/:username/post/:post_id`
- Retrieves a post by its ID.
- **Response**:
    ```json
    {
      "status": 200,
      "message": "OK",
      "data": {
        "post_id": "post_id",
        "title": "title",
        "content": "content"
      }
    }
    ```

### 3. **Get All Posts**
- **GET** `/api/u/:username/post`
- Retrieves all posts for a specific user.
- **Query Parameters**:
    - `limit`: (Optional) Number of posts to retrieve (default is 10).
    - `offset`: (Optional) Page offset (default is 0).
- **Response**:
    ```json
    {
      "status": 200,
      "message": "posts fetched",
      "data": [
        {
          "post_id": "post_id",
          "title": "title",
          "content": "content"
        }
      ]
    }
    ```

### 4. **Update Post**
- **PATCH** `/api/u/:username/post/:post_id`
- Updates an existing post.
- **Request Body**:
    ```json
    {
      "title": "updated title",
      "content": "updated content"
    }
    ```
- **Response**:
    ```json
    {
      "status": 200,
      "message": "post updated",
      "data": {
        "post_id": "post_id",
        "title": "title",
        "content": "content"
      }
    }
    ```

### 5. **Delete Post**
- **DELETE** `/api/u/:username/post/:post_id`
- Deletes a post by its ID.
- **Response**:
    ```json
    {
      "status": 200,
      "message": "post deleted successfully",
      "data": null
    }
    ```

---

## Middleware
- **Authentication**: Some routes are protected and require a JWT token for access.
- **Authorization**: Only the owner of the items or posts can update or delete them.

---

## Error Responses
- **Bad Request (400)**: Invalid or missing data.
- **Unauthorized (401)**: Missing or invalid token.
- **Forbidden (403)**: User does not have permission to access the resource.
- **Internal Server Error (500)**: A server-side error occurred.

---

## Notes
- All routes that involve creating, updating, or deleting resources (items, posts) require the user to be authenticated using a valid JWT token in the `Authorization` header.
- Ensure to handle edge cases like invalid IDs, missing parameters, and unauthorized access accordingly.

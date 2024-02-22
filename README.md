# Go User Management App

This is a simple RESTful API for managing user data written in Go.

## Features

- Create, Read, Update, and Delete users (CRUD operations)
- Uses MySQL database for storing user data
- Implements a RESTful API using Gorilla Mux

## Prerequisites

- Go 1.13 or higher installed on your machine
- MySQL database server

## Installation

1. Clone the repository:
git clone [https://github.com/abhishek-ics/go-crud.git](https://github.com/abhishek-ics/go-crud.git)
cd go-crud

2. Set up your MySQL database and configure the environment variables in a `.env` file. You can use the `.env.example` file as a template.

3. Install dependencies:
go mod tidy


4. Run the application:
go run main.go



The server should be running at `http://localhost:8080`.

## Endpoints

- `GET /users`: Get all users
- `GET /users/{id}`: Get a user by ID
- `POST /users`: Create a new user
- `PUT /users/{id}`: Update an existing user
- `DELETE /users/{id}`: Delete a user by ID

## Environment Variables

- `DB_USERNAME`: MySQL database username
- `DB_PASSWORD`: MySQL database password
- `DB_HOST`: MySQL database host
- `DB_PORT`: MySQL database port
- `DB_NAME`: MySQL database name

## Creating the database schema

1. To create the database and schema, you can follow the instructions below:
Log into your database:
```bash 
sudo mysql -u root -p
```

2. Create the database using SQL commands or a MySQL management tool:

```bash 
CREATE DATABASE DB_NAME;
```
3. Select the newly created database:
```bash
USE gocrud_app;
```

4. Create the users table with the desired columns:
```bash
  CREATE TABLE users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  email VARCHAR(100),
);
```
# Go Key-Value Database (go-kv)

A simple in-memory key-value database implemented in Go. This application provides a RESTful API to store, retrieve, update, and delete key-value pairs.

## Features

- **GET /[key]**: Retrieve the value for a given key. Returns a 404 if the key does not exist.
- **PUT /[key]**: Set the value for a given key. Updates the value if the key already exists.
- **DELETE /[key]**: Delete the value for a given key. Returns a 404 if the key does not exist.
- **GET /**: Retrieve a list of all keys in the database.

## Security

The application uses the `github.com/unrolled/secure` package to apply basic security headers, ensuring a more secure web application.

## Requirements

- Go 1.20 or later

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mrdaleyoung/go-kv.git
   cd go-kv

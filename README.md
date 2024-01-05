# Snippetbox

## Overview

Snippetbox is a CRUD (Create, Read, Update, Delete) web application built with Go, MySQL, Bootstrap, and JavaScript. It serves as a training project to learn how to build production-ready web applications with Go, inspired by the concepts introduced in Alex Edwards "Let’s Go" book.

## Features

- **CRUD Functionality**: Create, read, update, and delete snippets with ease.
- **Web Technologies**: Utilizes Go for server-side logic, MySQL for data storage, Bootstrap for frontend styling, and JavaScript for dynamic interactions.
- **Third-party Libraries**:
  - **zerolog**: Fast and structured logging.
  - **lumberjack**: Log rolling files for efficient log management.
  - **goland-migrate**: SQL database migration for easy versioning and updates.
  - **justinas/alice**: Painless middleware chaining for Go.
  - **justinas/nosurf**: CSRF protection for Go web applications.
  - **julienschmidt/httprouter**: Efficient router for handling HTTP requests.

### Running

```
  cd snippetbox
  go run cmd/web/*
  go test -v cmd/web/   ( to run the tests )
  
```

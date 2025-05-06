# portal-cli

CLI helper for Portalnesia projects: a lightweight, modular tool to automate boilerplate code generation and initialization tasks in Go applications.

## Installation

```bash
go install go.portalnesia.com/portal-cli@latest
```

## Init Golang Project

Golang helper for init structure directory

Usage:

```bash
portal-cli golang init [flags]
```

Flags:

| Shorthand | Flag         | Description                |
|-----------|--------------|----------------------------|
| -h        | --help       | help for init              |
| -a        | --all        | Add all library            |
| -f        | --flag       | Use flag instead of prompt |
| -o        | --override   | Override existing files    |
| -         | --firebase   | Add firebase               |
| -         | --handlebars | Add handlebars             |
| -         | --redis      | Add redis                  |


## Add Service

Add new service and CRUD routes, handler, and usecase

Usage:

```bash
portal-cli golang add-service [flags]
```

Flags:

| Shorthand | Flag       | Description                                             |
|-----------|------------|---------------------------------------------------------|
| -h        | --help     | help for add-service                                    |
| -f        | --flag     | Use flag instead of prompt                              |
| -n        | --name     | Service name                                            |
| -v        | --version  | Endpoint version, example: v1. Default without version  |
| -p        | --path     | Endpoint path,  example: user. Default use service name |
| -o        | --override | Override existing files                                 |


## Add Endpoint

Add new endpoint to existing routes, handler, and usecase

Usage:

```bash
portal-cli golang add-endpoint [flags]
```

Flags:

| Shorthand | Flag      | Description                                          |
|-----------|-----------|------------------------------------------------------|
| -h        | --help    | help for add-service                                 |
| -f        | --flag    | Use flag instead of prompt                           |
| -s        | --service | Service name                                         |
| -n        | --name    | Method name, example: FollowUsers                    |
| -p        | --path    | Endpoint path,  example: /v1/user/:id/follow         |
| -m        | --method  | HTTP method, example: GET, POST, PUT, PATCH, DELETE) |



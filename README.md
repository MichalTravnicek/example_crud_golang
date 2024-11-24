# example_crud_golang
Crud operations in Go with Gorm and Gin

## Docker

Service is configured for running in Docker 
- Running in Docker is very simple, the repo includes a `docker-compose.yaml` file.
- Service is setup with Postgres database (`.env` file) and includes database web admin interface at http://localhost:8090 
- At the repository root, execute `docker compose up` command to deploy the database instance and the API service.
- To force rebuild use `docker compose up --build --force-recreate` command.


## Run without Docker
- Service can be launched standalone using `go run main.go` but database config in .env file must be loaded to environment (not implemented)
- For running tests from command line use `go test -v`


## TODO

- [X] Basic Go project with Github Actions CI/CD
- [X] Docker configuration with database
- [ ] Crud application with tests
- [ ] Refactoring - repository, dto etc.

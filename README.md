# example_crud_golang
Crud operations in Go with Gorm and Gin

## REST

1. **Add User to database** (POST)
   - Endpoint: `localhost:8080/users` or `localhost:8080/save`
   - Description: Add user to the database.
   - Responses: 201 (created), 400 (invalid uuid or date)
   - Example curl: 
   ```
   curl -X POST http://localhost:8080/save/ -d '{"id":"b05cc5a3-62d7-49ce-a094-f65a82caac5c","name":"Michal","email":"michal@michal.cz","date_of_birth":"2020-01-01T12:12:34+00:00"}
   ```
   - Response:
   ```JSON
   {
    "id": "b05cc5a3-62d7-49ce-a094-f65a82caac5c"
    "name": "Michal",
    "email": "michal@michal.cz",
    "date_of_birth": "2020-01-01T12:12:34+00:00"
   }
   ```

2. **Get User by ID** (GET)
   - Endpoint: `localhost:8080/users/{UUID}` or `localhost:8080/{UUID}`
   - Description: Get user by uuid.
   - Responses: 200 (ok), 400 (invalid uuid), 404 (user does not exist)
   - Example curl: 
   ```
   curl localhost:8080/b05cc5a3-62d7-49ce-a094-f65a82caac5c
   ``` 

   - Response:
   ```JSON
   {
    "id": "b05cc5a3-62d7-49ce-a094-f65a82caac5c"
    "name": "Michal",
    "email": "michal@michal.cz",
    "date_of_birth": "2020-01-01T12:12:34+00:00"
   }
   ```

3. **Update User by ID** (PUT)
   - Endpoint: `localhost:8080/users/{UUID}`
   - Description: Update data of user with this ID.

4. **Delete User by ID** (DELETE)
   - Endpoint: `localhost:8080/users/{UUID}`
   - Description: Delete user with this ID.


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
- [X] Crud application with tests
- [ ] Refactoring - repository, dto etc.

# Todo REST API

A basic REST API for todo app with PostgreSQL

## How to run

### Running with Docker

* Clone the repo
* Copy .env.sample to .env
* Run `docker-compose up --build`
* Go to [http://127.0.0.1:8080](http://127.0.0.1:8080)

### Running manually

* Clone the repo
* Copy .env.sample to .env
* Set postgres database URL with postgres variable in the env file
* Run `go mod download`
* Run `go build -o main`
* Run `./main`
* Go to [http://127.0.0.1:8080](http://127.0.0.1:8080)


## API Endpoints

| URL                                               | METHOD | Description                                     |
|---------------------------------------------------|--------|-------------------------------------------------|
| [/swaggerui/](http://127.0.0.1:8080/swaggerui/#/) | GET    | API documentation with Open API 3.0 and Swagger |
| [/todo/](http://127.0.0.1:8080/todo)              | GET    | Fetch all the todos                             |
| [/todo/list](http://127.0.0.1:8080/todo/list)     | GET    | Fetch all the todos                             |
| /todo/:id                                         | GET    | Fetch single todo                               |
| /todo/:id                                         | DELETE | Delete a todo                                   |
| /todo/create                                      | POST   | Creates a todo item with the given title prop   |
| /todo/update                                      | POST   | Updates a todo either with title or done props  |
| /user/register                                    | POST   | Register                                        |
| /user/login                                       | POST   | Login                                           |


## Roadmap

1. Todo REST API without DB ✅
2. Todo REST API with PostgreSQL ✅
3. Todo REST API with Auth ✅
4. Logging ✅
5. Reset password / EMAIL integration ✅
6. Tests ✅ (only for two files)
7. Open API Integration ✅
8. Docker integration ✅
9. Refactor


## Questions

- Does storing the user in the context of a request make sense?

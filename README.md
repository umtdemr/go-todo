# Todo REST API without DB

A basic REST API for todo app without using database

## Run

* Clone the repo
* Create .env file
* Set postgres database URL with postgres variable in env file 
* Run go build
* Run ./go-todo (linux & mac
* Go to [http://127.0.0.1](http://127.0.0.1)


## API Endpoints

| URL     | METHOD | Description                                    |
|---------|--------|------------------------------------------------|
| /list   | GET    | Fetch all the todos                            |
| /:id    | GET    | Fetch single todo                              |
| /:id    | DELETE | Delete a todo                                  |
| /create | POST   | Creates a todo item with the given title prop  |
| /update | POST   | Updates a todo either with title or done props |


## Roadmap

1. Todo REST API without DB ✅
2. Todo REST API with PostgreSQL ✅
3. Todo REST API with Auth

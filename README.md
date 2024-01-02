# Todo REST API without DB

A basic REST API for todo app with PostgreSQL

## Run

* Clone the repo
* Create .env file
* Set postgres database URL with postgres variable in env file 
* Run go build
* Run ./go-todo (linux & mac
* Go to [http://127.0.0.1](http://127.0.0.1)


## API Endpoints

| URL            | METHOD | Description                                    |
|----------------|--------|------------------------------------------------|
| /todo/         | GET    | Fetch all the todos                            |
| /todo/list     | GET    | Fetch all the todos                            |
| /todo/:id      | GET    | Fetch single todo                              |
| /todo/:id      | DELETE | Delete a todo                                  |
| /todo/create   | POST   | Creates a todo item with the given title prop  |
| /todo/update   | POST   | Updates a todo either with title or done props |
| /user/register | POST   | Register                                       |
| /user/login    | POST   | Login                                          |


## Roadmap

1. Todo REST API without DB ✅
2. Todo REST API with PostgreSQL ✅
3. Todo REST API with Auth ✅


## Questions

- Does storing the user in the context of a request make sense?
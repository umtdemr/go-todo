package main

import (
	"encoding/json"
	"fmt"
	"github.com/umtdemr/go-todo/todo"
	"net/http"
	"os"
	"sync"
)

type APIServer struct {
	listenAddr string
	repository Repository
}

func NewAPIServer(listenAddr string, repository Repository) *APIServer {
	return &APIServer{listenAddr: listenAddr, repository: repository}
}

func (s *APIServer) Run() {
	http.HandleFunc("/list", s.handleList)
	http.HandleFunc("/create", s.handleAdd)
	//http.HandleFunc("/update", handleUpdate)
	//http.HandleFunc("/", handleFetchAndDelete)
	http.ListenAndServe(s.listenAddr, nil)
}

func Respond(w http.ResponseWriter, data interface{}) {
	isThereError := false
	jsonData, err := json.Marshal(data)

	if err != nil {
		isThereError = true
	}

	w.Header().Set("Content-Type", "application/json")
	_, writeError := w.Write(jsonData)

	if writeError != nil {
		isThereError = true
	}

	if isThereError {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func RespondWithError(w http.ResponseWriter, msg string, errCode int) {
	if msg == "" {
		msg = "An error has occurred while processing"
	}
	if errCode == 0 {
		errCode = http.StatusBadRequest
	}
	w.WriteHeader(errCode)

	resp := make(map[string]string)
	resp["message"] = msg

	encoder := json.NewEncoder(w)
	err := encoder.Encode(resp)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": %v}`, err)
	}
}

type Database struct {
	Todos *[]*todo.Todo
	mutex sync.Mutex
}

var db = Database{
	Todos: &[]*todo.Todo{todo.NewTodo("initial todo")},
}

func (s *APIServer) handleList(w http.ResponseWriter, r *http.Request) {
	todos, err := s.repository.GetAllTodos()

	if err != nil {
		RespondWithError(w, fmt.Sprintf("error while getting list: %s", err), http.StatusBadRequest)
		return
	}
	Respond(w, todos)
}

func (s *APIServer) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, "not valid", http.StatusBadRequest)
		return
	}

	var createTodoType todo.CreateTodoData

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&createTodoType)

	if err != nil {
		RespondWithError(w, fmt.Sprintf("parsing error: %v", err), http.StatusBadRequest)
		return
	}

	if createTodoType.Title == "" {
		RespondWithError(w, "Title need to be sent", http.StatusBadRequest)
		return
	}

	createdTodo := todo.NewTodo(createTodoType.Title)
	createErr := s.repository.CreateTodo(createdTodo)
	if createErr != nil {
		fmt.Fprintf(os.Stderr, "error while generating the todo: %s\n", createErr)
	}

	Respond(w, createdTodo)
}

//func handleUpdate(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		RespondWithError(w, "Only POST methods are allowed", http.StatusBadRequest)
//		return
//	}
//
//	var updateData todo.UpdateTodoData
//
//	decoder := json.NewDecoder(r.Body)
//
//	err := decoder.Decode(&updateData)
//
//	if err != nil {
//		RespondWithError(w, fmt.Sprintf("error while parsing: %s", err), http.StatusBadRequest)
//		return
//	}
//
//	id, uuidParseErr := uuid.Parse(updateData.Id)
//
//	if uuidParseErr != nil {
//		RespondWithError(w, "the id you sent is invalid", http.StatusBadRequest)
//		return
//	}
//
//	db.mutex.Lock()
//	todoIndex := slices.IndexFunc(*db.Todos, func(c *todo.Todo) bool { return c.Id == id })
//	if todoIndex == -1 {
//		RespondWithError(w, "We couldn't find the todo you are searching for", http.StatusBadRequest)
//		return
//	}
//	todoItem := (*db.Todos)[todoIndex]
//
//	if updateData.Title != nil {
//		todoItem.Title = *updateData.Title
//	}
//
//	if updateData.Done != nil {
//		todoItem.Done = *updateData.Done
//	}
//
//	todoItem.UpdatedAt = time.Now()
//	db.mutex.Unlock()
//
//	Respond(w, todoItem)
//}

/*func handleFetchAndDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		RespondWithError(w, "Only GET and DELETE requests are allowed", http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 2 {
		RespondWithError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	todoId := parts[1]

	parsedTodoId, parseUuidErr := uuid.Parse(todoId)

	if parseUuidErr != nil {
		RespondWithError(w, "the id you sent is invalid", http.StatusBadRequest)
		return
	}

	db.mutex.Lock()
	todoIndex := slices.IndexFunc(*db.Todos, func(c *todo.Todo) bool { return c.Id == parsedTodoId })
	if todoIndex == -1 {
		RespondWithError(w, "We couldn't find the todo you are searching for", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		nextId := todoIndex + 1
		if len(*db.Todos)-1 > nextId {
			nextId = len(*db.Todos) - 1
		}
		*(db.Todos) = slices.Delete(*db.Todos, todoIndex, nextId)
		db.mutex.Unlock()
		Respond(w, (*db.Todos)[todoIndex])
	} else {
		db.mutex.Unlock()
		Respond(w, (*db.Todos)[todoIndex])
	}
}
*/

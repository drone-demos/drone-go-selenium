// Copyright 2014 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/drone-demos/drone-go-selenium/task"
	"github.com/gorilla/mux"
)

var tasks, _ = task.NewTaskManager(
	"sqlite3",
	"todo.sqlite",
)

func RegisterHandlers() {
	r := mux.NewRouter()
	r.HandleFunc("/task/", errorHandler(ListTasks)).Methods("GET")
	r.HandleFunc("/task/", errorHandler(NewTask)).Methods("POST")
	r.HandleFunc("/task/{id}", errorHandler(GetTask)).Methods("GET")
	r.HandleFunc("/task/{id}", errorHandler(UpdateTask)).Methods("PUT")
	r.HandleFunc("/task/{id}", errorHandler(DeleteTask)).Methods("DELETE")
	http.Handle("/task/", r)
}

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

// errorHandler wraps a function returning an error by handling the error and returning a http.Handler.
// If the error is of the one of the types defined above, it is handled as described for every type.
// If the error is of another type, it is considered as an internal error and its message is logged.
func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}

// ListTask handles GET requests on /task.
//
// Example:
//
//   req: GET /task/
//   res: 200 [
//          {"ID": 1, "Title": "Learn Go"},
//          {"ID": 2, "Title": "Buy bread"}
//        ]
func ListTasks(w http.ResponseWriter, r *http.Request) error {
	all, _ := tasks.List()
	return json.NewEncoder(w).Encode(all)
}

// NewTask handles POST requests on /task.
// The request body must contain a JSON object with a Title field.
// The status code of the response is used to indicate any error.
//
// Examples:
//
//   req: POST /task/ {"Title": ""}
//   res: 400 empty title
func NewTask(w http.ResponseWriter, r *http.Request) error {
	task := task.Task{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return badRequest{err}
	}
	return tasks.Save(&task)
}

// GetTask handles GET requsts to /task/{taskID}.
// There's no parameters and it returns a JSON encoded task.
//
// Examples:
//
//   req: GET /task/1
//   res: 200 {"ID": 1, "Title": "Buy bread", "Done": true}
//
//   req: GET /task/42
//   res: 404 task not found
func GetTask(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	log.Println("Task is ", id)
	if err != nil {
		return badRequest{err}
	}
	t, err := tasks.Find(id)
	if err != nil {
		return notFound{}
	}
	return json.NewEncoder(w).Encode(t)
}

// DeleteTask handles DELETE requests to /task/{taskID}.
//
// Example:
//
//   req: DELETE /task/1
//   res: 200
func DeleteTask(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}
	return tasks.Delete(id)
}

// UpdateTask handles PUT requests to /task/{taskID}.
// The request body must contain a JSON encoded task.
//
// Example:
//
//   req: PUT /task/1 {"ID": 1, "Title": "Learn Go", "Done": true}
//   res: 200
//
//   req: PUT /task/2 {"ID": 2, "Title": "Learn Go", "Done": true}
//   res: 400 inconsistent task IDs
func UpdateTask(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}
	var t task.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return badRequest{err}
	}
	if t.ID != id {
		return badRequest{fmt.Errorf("inconsistent task IDs")}
	}
	if _, err := tasks.Find(id); err != nil {
		return notFound{}
	}
	return tasks.Update(&t)
}

// parseID obtains the id variable from the given request url,
// parses the obtained text and returns the result.
func parseID(r *http.Request) (int64, error) {
	txt, ok := mux.Vars(r)["id"]
	if !ok {
		return 0, fmt.Errorf("task id not found")
	}
	return strconv.ParseInt(txt, 10, 0)
}

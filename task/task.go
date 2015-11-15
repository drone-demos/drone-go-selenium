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
package task

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID    int64  // Unique identifier
	Title string // Description
	Done  bool   // Is this task done?
}

// TaskManager manages a list of tasks in a sql database.
type TaskManager struct {
	db *sql.DB // Database connection
}

// Find returns the Task with the given id in the database.
func (t *TaskManager) Find(id int64) (*Task, error) {
	task := &Task{}
	row := t.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Done,
	)
	return task, err
}

// All returns the list of all the Tasks in the database.
func (t *TaskManager) List() ([]*Task, error) {
	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		task := &Task{}
		err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Done,
		)
		if err != nil {
			break
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

// Save saves the given Task in the database.
func (t *TaskManager) Save(task *Task) error {
	res, err := t.db.Exec("INSERT INTO tasks VALUES (null, ?, ?)", task.Title, task.Done)
	if err != nil {
		return err
	}
	task.ID, err = res.LastInsertId()
	return err
}

// Update updates the given Task in the database.
func (t *TaskManager) Update(task *Task) error {
	_, err := t.db.Exec("UPDATE tasks SET title=?, done=? WHERE id=?", task.Title, task.Done, task.ID)
	return err
}

// Delete deltes the Task with the given id in the database.
func (t *TaskManager) Delete(id int64) error {
	_, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

// NewTaskManager returns a TaskManager with a sql database
// setup and configured.
func NewTaskManager(driver, datasource string) (*TaskManager, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}
	return &TaskManager{db}, nil
}

const schema = `
CREATE TABLE IF NOT EXISTS tasks (
	id    INTEGER PRIMARY KEY AUTOINCREMENT, 
	title TEXT,
	done  BOOLEAN
);
`

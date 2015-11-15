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
	"testing"
)

var testTask = Task{Title: "download drone"}

var testTasks = []Task{
	Task{Title: "download drone"},
	Task{Title: "setup continuous integration"},
	Task{Title: "profit"},
}

var tasks *TaskManager

func TestSave(t *testing.T) {
	setup()
	defer teardown()

	err := tasks.Save(&testTask)
	if err != nil {
		t.Errorf("Wanted to save task, got error. %s", err)
	}
	if testTask.ID == 0 {
		t.Errorf("Wanted tasks id assignment, got 0")
	}

	after, _ := tasks.List()
	if len(after) != 1 {
		t.Errorf("Wanted 1 item in the task list, got %d tasks", len(after))
	}
}

func TestList(t *testing.T) {
	setup()
	defer teardown()

	for _, task := range testTasks {
		tasks.Save(&task)
	}

	list, err := tasks.List()
	if err != nil {
		t.Errorf("Error listing task items. %s", err)
	}
	if len(list) != len(testTasks) {
		t.Errorf("Wanted %d items in list, got %d", len(testTasks), len(list))
	}
}

func TestDelete(t *testing.T) {
	setup()
	defer teardown()

	err := tasks.Save(&testTask)
	if err != nil {
		t.Errorf("Wanted to save tasl, got error. %s", err)
	}

	err = tasks.Delete(testTask.ID)
	if err != nil {
		t.Errorf("Wanted to delete task, got error. %s", err)
	}

	after, _ := tasks.List()
	if len(after) != 0 {
		t.Errorf("Wanted empty task list, got %d tasks", len(after))
	}
}

func setup() {
	tasks, _ = NewTaskManager("sqlite3", ":memory:")
	tasks.db.Exec("DELETE FROM tasks")
}

func teardown() {
	tasks.db.Close()
}

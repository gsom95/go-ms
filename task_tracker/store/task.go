package taskstore

import (
	"errors"
	"sync"
	"time"
)

// Task represents a task with a due date.
type Task struct {
	ID   int       `json:"id"`
	Text string    `json:"text"`
	Due  time.Time `json:"due"`
}

// TaskStore is a in-memory database of tasks.
// Methods are "safe" to use concurrently.
type TaskStore struct {
	sync.Mutex

	tasks  map[int]Task
	nextID int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.nextID = 0
	ts.tasks = make(map[int]Task)

	return ts
}

// CreateTask save a new task to the storage.
func (ts *TaskStore) CreateTask(text string, due time.Time) (int, error) {
	ts.Lock()
	defer ts.Unlock()

	t := Task{
		ID:   ts.nextID,
		Text: text,
		Due:  due,
	}
	ts.tasks[ts.nextID] = t
	ts.nextID++
	return t.ID, nil
}

// ErrTaskNotFound occurs when a task is not in the storage.
var ErrTaskNotFound = errors.New("task is not found")

// GetTask searches for the task and returns it.
// If the task is not found, returns the ErrTaskNotFound error.
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	}

	return Task{}, ErrTaskNotFound
}

// DeleteTask removes the task from the storage.
// If the task is not found, returns the ErrTaskNotFound error.
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return ErrTaskNotFound
	}
	delete(ts.tasks, id)
	return nil
}

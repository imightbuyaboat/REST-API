package basic_types

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TaskStore interface {
	AddTask(task *Task) error
	GetTask(id int) (*Task, error)
	GetAllTasks() ([]Task, error)
	UpdateTask(task *Task) (*Task, error)
	DeleteTask(id int) error
}

type TaskCache interface {
	Get(taskID int) (*Task, error)
	Set(task *Task) error
	Delete(taskID int) error
}

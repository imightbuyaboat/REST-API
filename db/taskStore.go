package db

import (
	bt "restapi/basic_types"
)

type TaskStore interface {
	AddTask(task *bt.Task) error
	GetTask(id int) (*bt.Task, error)
	GetAllTasks() ([]bt.Task, error)
	UpdateTask(task *bt.Task) (*bt.Task, error)
	DeleteTask(id int) error
	CheckUser(data *UserData) (int, error)
}

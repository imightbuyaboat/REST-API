package cache

import (
	bt "restapi/basic_types"
)

type TaskCache interface {
	Get(taskID int) (*bt.Task, error)
	Set(task *bt.Task) error
	Delete(taskID int) error
}

package mocks

import (
	bt "restapi/basic_types"
	"restapi/db"

	"github.com/stretchr/testify/mock"
)

type MockTaskStore struct {
	mock.Mock
}

func (m *MockTaskStore) AddTask(task *bt.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskStore) GetTask(id int) (*bt.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*bt.Task), args.Error(1)
}

func (m *MockTaskStore) GetAllTasks() ([]bt.Task, error) {
	args := m.Called()
	return args.Get(0).([]bt.Task), args.Error(1)
}

func (m *MockTaskStore) DeleteTask(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskStore) UpdateTask(task *bt.Task) (*bt.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*bt.Task), args.Error(1)
}

func (m *MockTaskStore) CheckUser(data *db.UserData) (int, error) {
	args := m.Called(data)
	return args.Get(0).(int), args.Error(1)
}

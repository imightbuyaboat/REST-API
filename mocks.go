package main

import (
	bt "restapi/basic_types"

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

type MockTaskCache struct {
	mock.Mock
}

func (m *MockTaskCache) Get(taskID int) (*bt.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(*bt.Task), args.Error(1)
}

func (m *MockTaskCache) Set(task *bt.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskCache) Delete(taskID int) error {
	args := m.Called(taskID)
	return args.Error(0)
}

package mocks

import (
	bt "restapi/basic_types"

	"github.com/stretchr/testify/mock"
)

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

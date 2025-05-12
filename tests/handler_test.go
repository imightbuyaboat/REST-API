package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	bt "restapi/basic_types"
	"restapi/cache"
	db "restapi/db"
	"restapi/handler"
	"restapi/tests/mocks"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/mock"
)

func TestCreateTaskHandler(t *testing.T) {
	mockDB := &mocks.MockTaskStore{}
	mockCache := &mocks.MockTaskCache{}
	h := &handler.Handler{DB: mockDB, Cache: mockCache}

	type taskInfo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	tests := []struct {
		name             string
		inputID          int
		inputInfo        taskInfo
		mockAddTaskError error
		expectedStatus   int
		expectedResponse bt.Task
	}{
		{
			name:    "Succesfully create task",
			inputID: 1,
			inputInfo: taskInfo{
				Name:        "Test Task",
				Description: "Test Description",
			},
			mockAddTaskError: nil,
			expectedStatus:   http.StatusCreated,
			expectedResponse: bt.Task{
				ID:          1,
				Name:        "Test Task",
				Description: "Test Description",
			},
		},
		{
			name:    "Task already exists",
			inputID: 1,
			inputInfo: taskInfo{
				Name:        "Test Task",
				Description: "Test Description",
			},
			mockAddTaskError: db.ErrTaskAlreadyExists,
			expectedStatus:   http.StatusConflict,
			expectedResponse: bt.Task{},
		},
		{
			name:    "Invalid request body",
			inputID: 0,
			inputInfo: taskInfo{
				Name:        "",
				Description: "",
			},
			mockAddTaskError: nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: bt.Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := strconv.Itoa(tt.inputID)

			mockDB.ExpectedCalls = nil
			mockDB.On("AddTask", mock.AnythingOfType("*basic_types.Task")).Return(tt.mockAddTaskError)

			body, _ := json.Marshal(tt.inputInfo)

			req, err := http.NewRequest("POST", "/tasks/"+id, bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{
				"id": id,
			})

			rr := httptest.NewRecorder()
			h.CreateTaskHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				var responseTask bt.Task
				if err := json.NewDecoder(rr.Body).Decode(&responseTask); err != nil {
					t.Fatal(err)
				}
				if responseTask != tt.expectedResponse {
					t.Errorf("Expected response %v, got %v", tt.expectedResponse, responseTask)
				}
			}
		})
	}
}

func TestGetTaskHandler(t *testing.T) {
	mockDB := &mocks.MockTaskStore{}
	mockCache := &mocks.MockTaskCache{}
	h := &handler.Handler{DB: mockDB, Cache: mockCache}

	tests := []struct {
		name             string
		taskID           string
		cachedTask       *bt.Task
		cachedGetError   error
		dbTask           *bt.Task
		dbGetError       error
		expectedStatus   int
		expectedResponse bt.Task
	}{
		{
			name:   "Succesfully get task from cache",
			taskID: "1",
			cachedTask: &bt.Task{
				ID:          1,
				Name:        "Test Task",
				Description: "Test Description",
			},
			cachedGetError: nil,
			dbTask:         nil,
			dbGetError:     nil,
			expectedStatus: http.StatusOK,
			expectedResponse: bt.Task{
				ID:          1,
				Name:        "Test Task",
				Description: "Test Description",
			},
		},
		{
			name:           "Succesfully get task from DB",
			taskID:         "2",
			cachedTask:     nil,
			cachedGetError: nil,
			dbTask: &bt.Task{
				ID:          2,
				Name:        "Test Task",
				Description: "Test Description",
			},
			dbGetError:     nil,
			expectedStatus: http.StatusOK,
			expectedResponse: bt.Task{
				ID:          2,
				Name:        "Test Task",
				Description: "Test Description",
			},
		},
		{
			name:             "Task not found",
			taskID:           "100",
			cachedTask:       nil,
			cachedGetError:   cache.ErrTaskNotFound,
			dbTask:           nil,
			dbGetError:       db.ErrTaskNotFound,
			expectedStatus:   http.StatusNotFound,
			expectedResponse: bt.Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := strconv.Atoi(tt.taskID)

			mockCache.ExpectedCalls = nil
			mockCache.On("Get", id).Return(tt.cachedTask, tt.cachedGetError)
			mockCache.On("Set", tt.dbTask).Return(nil)

			mockDB.ExpectedCalls = nil
			mockDB.On("GetTask", id).Return(tt.dbTask, tt.dbGetError)

			req, err := http.NewRequest("GET", "/tasks/"+tt.taskID, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{
				"id": tt.taskID,
			})

			rr := httptest.NewRecorder()
			h.GetTaskHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var responseTask bt.Task
				if err := json.NewDecoder(rr.Body).Decode(&responseTask); err != nil {
					t.Fatal(err)
				}

				if responseTask != tt.expectedResponse {
					t.Errorf("Expected response %v, got %v", tt.expectedResponse, responseTask)
				}
			}
		})
	}
}

func TestGetAllTasksHandler(t *testing.T) {
	mockDB := &mocks.MockTaskStore{}
	mockCache := &mocks.MockTaskCache{}
	h := &handler.Handler{DB: mockDB, Cache: mockCache}

	tests := []struct {
		name             string
		dbTasks          []bt.Task
		dbGetError       error
		expectedStatus   int
		expectedResponse []bt.Task
	}{
		{
			name: "Succesfully get all tasks",
			dbTasks: []bt.Task{
				{
					ID:          1,
					Name:        "Test Task 1",
					Description: "Test Description 1",
				},
				{
					ID:          2,
					Name:        "Test Task 2",
					Description: "Test Description 2",
				},
			},
			dbGetError:     nil,
			expectedStatus: http.StatusOK,
			expectedResponse: []bt.Task{
				{
					ID:          1,
					Name:        "Test Task 1",
					Description: "Test Description 1",
				},
				{
					ID:          2,
					Name:        "Test Task 2",
					Description: "Test Description 2",
				},
			},
		},
		{
			name:             "Failed to get all tasks",
			dbTasks:          nil,
			dbGetError:       fmt.Errorf("failed to select tasks from DB"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB.ExpectedCalls = nil
			mockDB.On("GetAllTasks").Return(tt.dbTasks, tt.dbGetError)

			req, err := http.NewRequest("GET", "/tasks", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			h.GetAllTasksHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var responseTask []bt.Task
				if err := json.NewDecoder(rr.Body).Decode(&responseTask); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(responseTask, tt.expectedResponse) {
					t.Errorf("Expected response %v, got %v", tt.expectedResponse, responseTask)
				}
			}
		})
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	mockDB := &mocks.MockTaskStore{}
	mockCache := &mocks.MockTaskCache{}
	h := &handler.Handler{DB: mockDB, Cache: mockCache}

	type taskInfo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	tests := []struct {
		name              string
		inputID           int
		inputInfo         taskInfo
		cachedDeleteError error
		dbTask            *bt.Task
		dbUpdateError     error
		expectedStatus    int
		expectedResponse  bt.Task
	}{
		{
			name:    "Succesfully update task",
			inputID: 1,
			inputInfo: taskInfo{
				Name:        "Test Task",
				Description: "Test Description",
			},
			cachedDeleteError: nil,
			dbTask: &bt.Task{
				ID:          1,
				Name:        "Test Task",
				Description: "Test Description",
			},
			dbUpdateError:  nil,
			expectedStatus: http.StatusOK,
			expectedResponse: bt.Task{
				ID:          1,
				Name:        "Test Task",
				Description: "Test Description",
			},
		},
		{
			name:    "Invalid request body",
			inputID: 0,
			inputInfo: taskInfo{
				Name:        "",
				Description: "",
			},
			cachedDeleteError: nil,
			dbTask:            nil,
			dbUpdateError:     nil,
			expectedStatus:    http.StatusBadRequest,
			expectedResponse:  bt.Task{},
		},
		{
			name:    "Task not found",
			inputID: 100,
			inputInfo: taskInfo{
				Name:        "Test Task",
				Description: "Test Description",
			},
			cachedDeleteError: nil,
			dbTask:            nil,
			dbUpdateError:     db.ErrTaskNotFound,
			expectedStatus:    http.StatusNotFound,
			expectedResponse:  bt.Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := strconv.Itoa(tt.inputID)

			mockCache.ExpectedCalls = nil
			mockCache.On("Delete", tt.inputID).Return(tt.cachedDeleteError)

			mockDB.ExpectedCalls = nil
			mockDB.On("UpdateTask", mock.AnythingOfType("*basic_types.Task")).Return(tt.dbTask, tt.dbUpdateError)

			body, _ := json.Marshal(tt.inputInfo)

			req, err := http.NewRequest("PUT", "/tasks/"+id, bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{
				"id": id,
			})

			rr := httptest.NewRecorder()
			h.UpdateTaskHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var responseTask bt.Task
				if err := json.NewDecoder(rr.Body).Decode(&responseTask); err != nil {
					t.Fatal(err)
				}
				if responseTask != tt.expectedResponse {
					t.Errorf("Expected response %v, got %v", tt.expectedResponse, responseTask)
				}
			}
		})
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	mockDB := &mocks.MockTaskStore{}
	mockCache := &mocks.MockTaskCache{}
	h := &handler.Handler{DB: mockDB, Cache: mockCache}

	tests := []struct {
		name              string
		taskID            string
		cachedDeleteError error
		dbDeleteError     error
		expectedStatus    int
	}{
		{
			name:              "Succesfully delete task",
			taskID:            "1",
			cachedDeleteError: nil,
			dbDeleteError:     nil,
			expectedStatus:    http.StatusNoContent,
		},
		{
			name:              "Task not found",
			taskID:            "100",
			cachedDeleteError: cache.ErrTaskNotFound,
			dbDeleteError:     db.ErrTaskNotFound,
			expectedStatus:    http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := strconv.Atoi(tt.taskID)

			mockCache.ExpectedCalls = nil
			mockCache.On("Delete", id).Return(tt.cachedDeleteError)

			mockDB.ExpectedCalls = nil
			mockDB.On("DeleteTask", id).Return(tt.dbDeleteError)

			req, err := http.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{
				"id": tt.taskID,
			})

			rr := httptest.NewRecorder()
			h.DeleteTaskHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

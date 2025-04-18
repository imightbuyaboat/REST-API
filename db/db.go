package db

import (
	"database/sql"
	"fmt"
	"os"
	bt "restapi/basic_types"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("SQL_HOST"), os.Getenv("SQL_PORT"),
		os.Getenv("SQL_USER"), os.Getenv("SQL_PASSWORD"),
		os.Getenv("SQL_DB"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresStore{DB: db}, nil
}

func (ps *PostgresStore) AddTask(task *bt.Task) error {
	var exists bool
	query := "select EXISTS (select 1 from tasks where id = $1)"
	err := ps.DB.QueryRow(query, task.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if task %d exists: %v", task.ID, err)
	}
	if exists {
		return ErrTaskAlreadyExists
	}

	query = "insert into tasks (id, name, description) values ($1, $2, $3)"
	_, err = ps.DB.Exec(query, task.ID, task.Name, task.Description)
	if err != nil {
		return fmt.Errorf("failed to insert task %d: %v", task.ID, err)
	}
	return nil
}

func (ps *PostgresStore) GetTask(taskID int) (*bt.Task, error) {
	var task bt.Task
	query := "select id, name, description from tasks where id = $1"

	err := ps.DB.QueryRow(query, taskID).Scan(&task.ID, &task.Name, &task.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to select task %d from DB: %v", taskID, err)
	}
	return &task, nil
}

func (ps *PostgresStore) GetAllTasks() ([]bt.Task, error) {
	query := "select id, name, description from tasks"

	rows, err := ps.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to select tasks from DB: %v", err)
	}
	defer rows.Close()

	var tasks []bt.Task
	for rows.Next() {
		var task bt.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Description); err != nil {
			return nil, fmt.Errorf("failed to scan task %d from DB: %v", len(tasks)+1, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (ps *PostgresStore) UpdateTask(task *bt.Task) (*bt.Task, error) {
	var exists bool
	query := "select EXISTS (select 1 from tasks where id = $1)"
	err := ps.DB.QueryRow(query, task.ID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if task %d exists: %v", task.ID, err)
	}
	if !exists {
		return nil, ErrTaskNotFound
	}

	query = "update tasks set name = $1, description = $2 where id = $3 returning id, name, description"
	var updatedTask bt.Task

	err = ps.DB.QueryRow(query, task.Name, task.Description, task.ID).Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to update task %d: %v", task.ID, err)
	}
	return &updatedTask, nil
}

func (ps *PostgresStore) DeleteTask(taskID int) error {
	query := "delete from tasks where id = $1"

	res, err := ps.DB.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task %d from DB: %v", taskID, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

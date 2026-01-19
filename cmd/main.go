package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v2"
)

var pool *pgxpool.Pool
var ctx = context.Background()

type Task struct {
  ID int `json:"id"`
  Text string `json:"text"`
  Completed bool `json:"completed"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

func init() {
  pool, err := pgxpool.New(ctx, "postgresql://postgres:postgres@localhost:5432/gotodoapp")
  if err != nil {
    log.Fatal("Unable to connect to database", err)
  }

  if err := pool.Ping(ctx); err != nil {
    log.Fatal("Unable to ping the database", err)
  }

  fmt.Println("Connected to the database successfully")
}

func main() {
  app := &cli.App{
    Name: "GoTodo",
    Usage: "A simple cli program to manage your tasks",
    Commands: []*cli.Command{},
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}

func createTask(text string) error {
  sql := `
    INSERT INTO tasks (text, completed)
    VALUES ($1, $2)
    RETURNING id
  `

  var id int 
  err := pool.QueryRow(ctx, sql, text, false).Scan(&id)
  if err != nil {
    return fmt.Errorf("error creating task: %w", err)
  }

  fmt.Printf("Created task successfully with ID: %d\n", id)
  return nil
}

func getAllTasks() ([]Task, error) {
  sql := `
    SELECT id, text, completed, created_at, updated_at
    FROM tasks
    ORDER BY created_at DESC
  `

  rows, err := pool.Query(ctx, sql)
  if err != nil {
    return nil, fmt.Errorf("error querying  tasks: %w", err)
  }
  defer rows.Close()

  var tasks []Task
  for rows.Next() {
    var task Task
    err := rows.Scan(
      &task.ID,
      &task.Text,
      &task.Completed,
      &task.CreatedAt,
      &task.UpdatedAt,
    )

    if err !=nil {
      return nil, fmt.Errorf("Error scan task row: %w", err)
    }

    tasks = append(tasks, task)
  }

  return tasks, nil
}

func completeTask(id int) error {
  sql := `
    UPDATE tasks
    SET completed = true, updated_at = NOW()
    WHERE id = $1
  `

  commandTag, err := pool.Exec(ctx, sql, id)
  if err != nil {
    return fmt.Errorf("error completing task: %w", err)
  }

  if commandTag.RowsAffected() == 0 {
    return fmt.Errorf("no task found with id %d", id)
  }

  return nil
}

func deleteTask(id int) error {
  sql := `DELETE FROM tasks WHERE id = $1`

  commandTag, err := pool.Exec(ctx, sql, id)
  if err != nil {
    return fmt.Errorf("error while deleting task with id: %d", id)
  }

  if commandTag.RowsAffected() == 0 {
    return fmt.Errorf("no task found with id %d", id)
  }

  return nil
}

func getPendingTasks() ([]Task, error) {
  sql := `
    SELECT id, text, completed, created_at, updated_at 
    FROM tasks
    WHERE completed = false
    ORDER BY created_at DESC
  `

  rows, err := pool.Query(ctx, sql)
  if err != nil {
    return nil, fmt.Errorf("error querying pending tasks: %w", err)
  }
  defer rows.Close()
  
  var tasks []Task 
  for rows.Next() {
    var task Task 
    err := rows.Scan(
      &task.ID,
      &task.Text,
      &task.Completed,
      &task.CreatedAt,
      &task.UpdatedAt,
    )
    if err != nil {
      return nil, fmt.Errorf("error scanning task row: %w", err)
    }

    tasks = append(tasks, task)
  }

  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("error iterating task row: %w", err)
  }

  return tasks, nil
}

func getCompletedTasks() ([]Task, error) {
  sql := `
    SELECT id, text, completed, created_at, updated_at 
    FROM tasks
    WHERE completed = true
    ORDER BY created_at DESC
  `

  rows, err := pool.Query(ctx, sql)
  if err != nil {
    return nil, fmt.Errorf("error querying pending tasks: %w", err)
  }
  defer rows.Close()
  
  var tasks []Task 
  for rows.Next() {
    var task Task 
    err := rows.Scan(
      &task.ID,
      &task.Text,
      &task.Completed,
      &task.CreatedAt,
      &task.UpdatedAt,
    )
    if err != nil {
      return nil, fmt.Errorf("error scanning task row: %w", err)
    }

    tasks = append(tasks, task)
  }

  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("error iterating task row: %w", err)
  }

  return tasks, nil
}

func printTask(tasks []Task) {
  if len(tasks) == 0 {
    fmt.Println("No tasks found")
    return
  }

  for _, task := range tasks {
    status := "[ ]"
    if task.Completed {
      status = "[✓]"
    }

    fmt.Printf("%d. %s %s (Created: %s)\n", task.ID, status, task.Text, task.CreatedAt.Format("2006-01-02 15:04:05"))
  }
}
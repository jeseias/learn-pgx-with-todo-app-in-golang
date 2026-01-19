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
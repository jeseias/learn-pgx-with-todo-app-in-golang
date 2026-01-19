package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

}
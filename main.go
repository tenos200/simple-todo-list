package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"todoList/dbhandler"
	"todoList/taskhandler"
)

const filePath = "./test.db"

func main() {

	// if we can't open the file then we should create the schema for the db.
	file, openErr := os.Open(filePath)
	if openErr != nil {
		fmt.Println(openErr)
		createDbFile()
		dbhandler.CreateSchema(filePath)
	}
	// If it already exist, ensure to close.
	defer file.Close()
	// Start the task
	taskhandler.TodoListRunner(filePath)
	os.Exit(0)
}

// createDbFile creates a database on file should it not exist.
func createDbFile() {
	_, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Something went wrong when creating file", "err", err)
	}
}

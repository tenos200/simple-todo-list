package main

import (
	"fmt"
	"log"
	"os"
	"todoList/dbhandler"
	"todoList/taskhandler"

	_ "github.com/mattn/go-sqlite3"
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
	// if it already exist, ensure to close then we move to "game" loop
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

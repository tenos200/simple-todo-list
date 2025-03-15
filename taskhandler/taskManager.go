package taskhandler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"todoList/dbhandler"
)

// TodoListRunner the entry point for the todo list, handles the event loop
// for the todo list. Takes a file path to the required database file.
func TodoListRunner(filePath string) {

	scanner := bufio.NewScanner(os.Stdin)
	// Fetch relevant entries from db, and show menu
	dbOutputCache := dbhandler.GetDbRows(filePath)
	for {
		showMenu()
		scanner.Scan()
		text := scanner.Text()

		// Switch statement for initialising different menu options.
		switch text {
		case "1":
			addTaskToList(&dbOutputCache)
		case "2":
			showTasks(&dbOutputCache)
		case "3":
			markAsDone()
		case "4":
			deleteTask()
		case "5":
			os.Exit(0)
		default:
			fmt.Println("Error: Invalid input")
		}

	}
}

// addTaskToList prompts the user to enter a new task and appends it to dbOutputCache.
// It continuously asks for task details until the user confirms the addition.
// The function takes a reference to a slice of dbhandler.DbRow and modifies it directly.
//
// Steps:
//
//  1. Prompts the user for a task name.
//
//  2. Asks the user to select a task date (Tomorrow, Next Week, or custom).
//
//  3. Sets the task status to "Incomplete" by default.
//
//  4. Displays the task details and asks for confirmation to add the task.
//
//  5. If confirmed, the task is appended to dbOutputCache with an ID of -1
//
// (indicating a new, unsaved task in database).
//
//  6. Returns the updated reference to dbOutputCache.
func addTaskToList(dbOutputCache *[]dbhandler.DbRow) *[]dbhandler.DbRow {
	var rowToAdd dbhandler.DbRow
	// For loop for adding task
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Add task name:")
		scanner.Scan()
		taskName := scanner.Text()
		rowToAdd.Name = taskName
		fmt.Println(
			"Add task date:\n1.Tomorrow\n2.Next Week\n3.Enter own date.")
		scanner.Scan()
		dateInput := scanner.Text()
		rowToAdd.Date = getInputDate(dateInput)
		// set default status to imcomplete
		rowToAdd.Status = "Incomplete"
		fmt.Printf("Task: %s\nDate :%s\n", rowToAdd.Name, rowToAdd.Date)
		fmt.Printf("Would you like to add this task to list.\n[y/n]\n")
		scanner.Scan()
		userInput := scanner.Text()
		if userInput == "y" {
			rowToAdd.Id = -1
			*dbOutputCache = append(*dbOutputCache, rowToAdd)
			fmt.Println("Task added!")
			break
		}
	}
	return dbOutputCache
}

// getInputDate takes a userInput from 1 -> returns tomorrows date, 2 -> returns
// next weeks date. 3 -> allows you to input a date.
//
// Example:
//
//	date := getInputDate("1")
//	fmt.Println(date)
//
// Output:
//
//	"2006-01-02"
func getInputDate(userDateChoice string) string {
	var pickedDate string
	currentDate := time.Now()
	switch userDateChoice {
	case "1":
		tomorrow := currentDate.Add(24 * time.Hour)
		pickedDate = tomorrow.Format("2006-01-02")
	case "2":
		nextWeek := currentDate.Add(time.Hour * 24 * 7)
		pickedDate = nextWeek.Format("2006-01-02")
	case "3":
		// Initialise a new scanner and take the date.
		fmt.Printf("Enter date formatted as '2006-01-02': ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		dateOutput := scanner.Text()
		parsedDate, dateErr := time.Parse("2006-01-02", dateOutput)
		if dateErr != nil {
			fmt.Println("Incorrectly formatted date", "date", dateErr)
		}
		pickedDate = parsedDate.Format("2006-01-02")
	}
	return pickedDate
}

// TODO: fix so that the output is correct.
func showTasks(tasks *[]dbhandler.DbRow) {

	if len(*tasks) == 0 {
		fmt.Println("No tasks have been added.")
	} else {
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("| %-25s | %-10s | %-10s |\n", "Task Name", "Status", "Date")
		fmt.Println(strings.Repeat("=", 50))

		for _, task := range *tasks {
			fmt.Printf("| %-25s | %-10s | %-10s |\n", task.Name, task.Status, task.Date)
		}

		fmt.Println(strings.Repeat("=", 50))
	}
}

func markAsDone() {
	fmt.Println("mark task as done")
}

func deleteTask() {
	fmt.Println("delete tasks")
}

// showMenu shows the output menu options for the todo list.
func showMenu() {
	menu := `====================================
          📌 TO-DO LIST 📌         
====================================
[1] ➜ Add a new task
[2] ➜ View tasks
[3] ➜ Mark task as done ✅
[4] ➜ Delete a task ❌
[5] ➜ Exit 🚪
------------------------------------
Enter your choice: `
	fmt.Println(menu)
}

package taskhandler

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"todoList/dbhandler"

	"github.com/charmbracelet/huh"
)

const timeFormat = "2006-01-02"

// TodoListRunner the takes a file path to required database file,
// handles the event loop. Provides the user with 5 choices for input:
//
// Picks:
//  1. Add a task
//  2. Show all tasks
//  3. Show all tasks
//  4. Mark a task as done
//  5. Exit program
func TodoListRunner(filePath string) {

	// Fetch relevant entries from db, and show menu
	dbOutputCache := dbhandler.GetDbRows(filePath)
	for {
		choice := showMenu()

		// Switch statement for initialising different menu options.
		switch choice {
		case 1:
			showTasks(&dbOutputCache)
		case 2:
			addTaskToList(&dbOutputCache)
		case 3:
			markAsDone(&dbOutputCache)
		case 4:
			deleteTask()
		case 5:
			os.Exit(0)
		default:
			fmt.Println("Error: Invalid input")
		}
	}
}

func addTaskToList(dbOutputCache *[]dbhandler.DbRow) *[]dbhandler.DbRow {
	var rowToAdd dbhandler.DbRow

	// For loop for adding task to list
	for {
		var confirmation bool

		// Add name to task.
		huh.NewInput().
			Title("Enter task name.").
			Prompt("> ").
			Value(&rowToAdd.Name).
			Run()

		// Get the date for the task
		rowToAdd.Date = getInputDate()

		// Set default status to imcomplete
		rowToAdd.Status = "Incomplete"
		formattedDate := fmt.Sprintf("Task: %s\nDate: %s\n",
			rowToAdd.Name,
			rowToAdd.Date)

		huh.NewConfirm().
			Title(formattedDate).
			Affirmative("Yes.").
			Negative("No.").
			Value(&confirmation).
			Run()

		if confirmation {
			*dbOutputCache = append(*dbOutputCache, rowToAdd)
			break
		}

	}
	return dbOutputCache
}

// getInputDate allows user to pick an input date from menu, returns a date string
// formatted as 2006-01-02.
func getInputDate() string {

	// Instansiate variables for huh input values
	var pickedDate string
	var switchSelect int
	currentDate := time.Now()
	huh.NewSelect[int]().
		Title("Input a date for task completion.").
		Options(
			huh.NewOption("Tomorrow.", 1),
			huh.NewOption("In 7 days.", 2),
			huh.NewOption("Enter own date.", 3),
		).
		Value(&switchSelect).
		Run()

	// Switch statement to pick the correct date
	switch switchSelect {
	case 1:
		tomorrow := currentDate.Add(24 * time.Hour)
		pickedDate = tomorrow.Format(timeFormat)
	case 2:
		nextWeek := currentDate.Add(time.Hour * 24 * 7)
		pickedDate = nextWeek.Format(timeFormat)
	case 3:

		// Initialise a new scanner and take the date.
		var dateInput string
		huh.NewInput().
			Title("Enter date formatted as '2006-01-02'").
			Prompt("> ").
			Value(&dateInput).
			Run()

		parsedDate, dateErr := time.Parse(timeFormat, dateInput)
		if dateErr != nil {
			fmt.Println("Incorrectly formatted date", "date", dateErr)
		}
		pickedDate = parsedDate.Format(timeFormat)
	}
	return pickedDate
}

// showTasks displays all tasks for tasklist.
func showTasks(tasks *[]dbhandler.DbRow) {

	if len(*tasks) == 0 {
		fmt.Println("No tasks have been added.")
	} else {
		for _, v := range *tasks {
			fmt.Printf("ID: %d\nName: %s\nStatus: %s\nDate: %s\n",
				v.Id, v.Name, v.Status, v.Date)
		}
	}
}

// TODO: Implement this function
func markAsDone(tasks *[]dbhandler.DbRow) {
	var selectedTasks []string
	var choices []huh.Option[string]

	// convert the tasks into huhNewOptions
	for _, v := range *tasks {
		convertedId := strconv.Itoa(v.Id)
		choices = append(choices, huh.NewOption(v.Name, convertedId))
	}
	huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Options").
				Options(choices...).
				Value(&selectedTasks),
		),
	).Run()
}

func deleteTask() {
	fmt.Println("delete tasks")
}

// showMenu displays a menu with options from 1-5, returns int from user input.
//
// Choices:
//  1. View Tasks
//  2. Add a new Task
//  3. Mark Complete
//  4. Delete Task
//  5. Exit
func showMenu() int {
	var userInput int
	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(
					huh.NewOption("View Tasks", 1),
					huh.NewOption("Add New Task", 2),
					huh.NewOption("Mark Complete", 3),
					huh.NewOption("Delete Task", 4),
					huh.NewOption("Exit", 5),
				).
				Value(&userInput),
		),
	).WithHeight(10).Run()
	return userInput
}

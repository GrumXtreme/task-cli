package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type TaskList []Task

var TasksLists TaskList

func main() {
	args := os.Args
	if len(args) < 2 {
		printError()
		return
	}

	file, err := OpenFile()
	if err != nil {
		fmt.Println("Error reading json file")
		return
	}
	defer file.Close()

	err = ReadFile(file, &TasksLists)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Println(err)
		return
	}

	process := args[1]
	switch process {
	case "add":
		if len(args) < 3 {
			fmt.Println("Please add a description to the task.")
			return
		}
		description := args[2]
		err := add(description)
		if err != nil {
			fmt.Println("Error adding a new task")
			return
		}
	case "update":
		if len(args) < 3 {
			fmt.Println("Please add an ID number")
			return
		}
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) < 4 {
			fmt.Println("Please add a description to the task.")
			return
		}
		description := args[3]
		err = update(id, description)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		if len(args) < 3 {
			fmt.Println("Please add an ID number")
			return
		}
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = delete(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "mark-in-progress":
		if len(args) < 3 {
			fmt.Println("Please add an ID number")
			return
		}
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = markProgress(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "mark-done":
		if len(args) < 3 {
			fmt.Println("Please add an ID number")
			return
		}
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = markDone(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "list":
		option := "all"
		if len(args) > 3 {
			fmt.Println("To many arguments")
			return
		}
		if len(args) == 3 {
			option = args[2]
		}
		err := list(option)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		printError()
	}
}

func printError() {
	fmt.Print("Please add a valid argument:\n 1.-add\n 2.-update\n 3.-delete\n 4.-mark-in-progress\n 5.-mark-done\n 6.-list\n")
}

func OpenFile() (*os.File, error) {
	file, err := os.OpenFile("task.json", os.O_RDWR, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create("task.json")
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return file, nil
}

func ReadFile(file *os.File, data *TaskList) error {
	decoder := json.NewDecoder(file)
	err := decoder.Decode(data)
	return err
}

func WriteFile(data *TaskList) error {
	file, err := os.Create("task.json")
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	return err
}

func add(description string) error {
	var newTask Task

	newTask.Description = description
	newTask.CreatedAt = time.Now().Format("2006/02/01 15:04:05")
	newTask.UpdatedAt = time.Now().Format("2006/02/01 15:04:05")
	newTask.Status = "todo"

	if isEmpty(&TasksLists) {
		newTask.ID = 1
	} else {
		newTask.ID = TasksLists[len(TasksLists)-1].ID + 1
	}

	TasksLists = append(TasksLists, newTask)

	err := WriteFile(&TasksLists)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)

	return nil
}

func delete(id int) error {
	if isEmpty(&TasksLists) {
		return fmt.Errorf("task list is empty")
	}

	if !isIDValid(TasksLists, id) {
		return fmt.Errorf("input a valid id")
	}

	TasksLists = append(TasksLists[:id-1], TasksLists[id:]...)

	err := WriteFile(&TasksLists)
	if err != nil {
		return err
	}
	return nil
}

func update(id int, description string) error {
	if isEmpty(&TasksLists) {
		return fmt.Errorf("task list is empty")
	}

	if !isIDValid(TasksLists, id) {
		return fmt.Errorf("input a valid id")
	}

	TasksLists[id-1].Description = description
	TasksLists[id-1].UpdatedAt = time.Now().Format("2006/02/01 15:04:05")

	err := WriteFile(&TasksLists)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func markProgress(id int) error {
	if isEmpty(&TasksLists) {
		return fmt.Errorf("task list is empty")
	}

	if !isIDValid(TasksLists, id) {
		return fmt.Errorf("input a valid id")
	}

	TasksLists[id-1].Status = "in-progress"
	TasksLists[id-1].UpdatedAt = time.Now().Format("2006/02/01 15:04:05")

	err := WriteFile(&TasksLists)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func markDone(id int) error {
	if isEmpty(&TasksLists) {
		return fmt.Errorf("task list is empty")
	}

	if !isIDValid(TasksLists, id) {
		return fmt.Errorf("input a valid id")
	}

	TasksLists[id-1].Status = "done"
	TasksLists[id-1].UpdatedAt = time.Now().Format("2006/02/01 15:04:05")

	err := WriteFile(&TasksLists)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func list(option string) error {
	switch option {
	case "all":
		for _, task := range TasksLists {
			fmt.Println(task)
		}
	case "done":
		for _, task := range TasksLists {
			if task.Status == "done" {
				fmt.Println(task)
			}
		}
	case "todo":
		for _, task := range TasksLists {
			if task.Status == "todo" {
				fmt.Println(task)
			}
		}
	case "in-progress":
		for _, task := range TasksLists {
			if task.Status == "in-progress" {
				fmt.Println(task)
			}
		}
	default:
		fmt.Println("Add a valid option (todo | done | in-progress)")
	}
	return nil
}

func isEmpty(list *TaskList) bool {
	return len(*list) == 0
}

func isIDValid(list []Task, id int) bool {
	for _, task := range list {
		if task.ID == id {
			return true
		}
	}
	return false
}

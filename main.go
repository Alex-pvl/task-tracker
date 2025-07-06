package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const taskFile = "/path/to/file.json"

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func loadTasks() ([]Task, error) {
	var tasks []Task
	file, err := os.Open(taskFile)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	file, err := os.Create(taskFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(tasks)
}

func printTask(task Task) {
	fmt.Printf("ID: %d | %s | %s\n", task.ID, task.Status, task.Description)
}

func printHelp() {
	fmt.Println(`Task CLI - Manage your to-do list from the command line

Usage:
  task-cli <command> [arguments]

Commands:
  add "task description"         Add a new task
  update <id> "new description"  Update a task's description
  delete <id>                    Delete a task
  mark-in-progress <id>          Mark a task as in progress
  mark-done <id>                 Mark a task as done
  list                           List all tasks
  list todo                      List tasks with status "todo"
  list in-progress               List tasks with status "in-progress"
  list done                      List tasks with status "done"

Options:
  -h, --help                     Show this help message
`)
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp()
		return
	}

	command := os.Args[1]
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add \"task description\"")
			return
		}
		id := len(tasks) + 1
		desc := strings.Join(os.Args[2:], " ")
		tasks = append(tasks, Task{
			ID:          id,
			Description: desc,
			Status:      "todo",
		})
		if err := saveTasks(tasks); err != nil {
			fmt.Println("Error saving task:", err)
			return
		}
		fmt.Println("Task", id, "added successfully")

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task-cli update <id> \"new description\"")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		desc := strings.Join(os.Args[3:], " ")

		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].ID < tasks[j].ID
		})
		index := sort.Search(len(tasks), func(i int) bool {
			return tasks[i].ID >= id
		})

		if index < len(tasks) && tasks[index].ID == id {
			tasks[index].Description = desc
			saveTasks(tasks)
			fmt.Println("Task", id, "updated.")
		} else {
			fmt.Println("Task", id, "not found.")
		}

	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: task-cli delete <id>")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		found := false
		newTasks := tasks[:0]
		for _, t := range tasks {
			if t.ID == id {
				found = true
				continue
			}
			newTasks = append(newTasks, t)
		}
		if !found {
			fmt.Println("Task", id, "not found.")
			return
		}
		saveTasks(newTasks)
		fmt.Println("Task", id, "deleted.")

	case "mark-in-progress", "mark-done":
		if len(os.Args) != 3 {
			fmt.Printf("Usage: task-cli %s <id>\n", command)
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].ID < tasks[j].ID
		})

		index := sort.Search(len(tasks), func(i int) bool {
			return tasks[i].ID >= id
		})

		if index < len(tasks) && tasks[index].ID == id {
			tasks[index].Status = strings.TrimPrefix(command, "mark-")
		} else {
			fmt.Println("Task", id, "not found.")
			return
		}
		saveTasks(tasks)
		fmt.Println("Task", id, "status updated.")

	case "list":
		if len(os.Args) == 2 {
			for _, t := range tasks {
				printTask(t)
			}
			return
		}
		status := os.Args[2]
		for _, t := range tasks {
			if t.Status == status {
				printTask(t)
			}
		}

	default:
		fmt.Println("Unknown command:", command)
	}
}

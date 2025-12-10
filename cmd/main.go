package main

import (
	"fmt"
	"os"
)

const TasksPath = "./tasks.csv"

func main() {

	tasks, err := GetTasks(TasksPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading tasks", err)
	}

	Run(tasks)
}

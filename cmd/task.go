package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	ListFlag   = flag.Int("list", -1, "List tasks optionally limit number")
	UpdateFlag = flag.Int("update", -1, "Update the value of a task by id")
	CheckFlag  = flag.Int("check", -1, "Check or Uncheck a task by id")
	DeleteFlag = flag.Int("delete", -1, "Delete task by id")
)

func Run(tasks []Task) {
	flag.Parse()
	pos := flag.Args()

	// List tasks
	if *ListFlag != -1 {
		if *ListFlag == 0 {
			DisplayTasks(tasks)
			return
		} else {
			start := len(tasks) - *ListFlag
			if start < 0 {
				end := len(tasks)
				DisplayTasks(tasks[:end])
				return
			}
			DisplayTasks(tasks[start:])
			return
		}
	}

	// Update task
	if *UpdateFlag != -1 {
		if len(pos) < 1 {
			println("Please provide a value to update.")
			return
		}
		err := UpdateTask(*UpdateFlag, pos[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error updating task:", err)
		}
		fmt.Printf("Updated task #%v", *UpdateFlag)
		return
	}

	// Check task
	if *CheckFlag != -1 {
		err := CheckTask(*CheckFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error checking task:", err)
		}
		fmt.Printf("Task #%v marked", *CheckFlag)
		return
	}

	// Delete task
	if *DeleteFlag != -1 {
		err := DeleteTask(*DeleteFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error deleting task:", err)
		}
		fmt.Printf("Deleted task #%v", *DeleteFlag)
		return
	}

	// if no flag add new task
	if len(pos) >= 1 {
		ids := make([]int, 0)

		for _, t := range tasks {
			ids = append(ids, t.id)
		}

		newId := slices.Max(ids) + 1

		newTask := Task{id: newId, value: pos[0], done: false}
		err := WriteTask(newTask, TasksPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating task:", err)
			return
		}

		fmt.Printf("Added task #%v", newTask.id)
		return
	}

	// NO COMMAND
	fmt.Println("Usage:")
	fmt.Println("'Task name' \n Create a new task")
	flag.PrintDefaults()
}

type Task struct {
	id    int
	value string
	done  bool
}

func (t *Task) IsDone() string {
	if t.done {
		return "Finished"
	} else {
		return "Unfinished"
	}
}

func GetTasks(path string) ([]Task, error) {
	// if file doesn't exist -> create with header and return empty slice
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("create file: %v", err)
		}

		defer f.Close()

		// Write file
		_, err = f.WriteString("id,value,done\n")
		if err != nil {
			return nil, fmt.Errorf("write header: %v", err)
		}
		return []Task{}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	rows, err := r.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("read csv: %v", err)
	}

	var tasks []Task

	for i, row := range rows {
		// skip header
		if i == 0 {
			continue
		}
		// skip fucked up row
		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		// skip fucked up id
		if err != nil {
			continue
		}

		done, _ := strconv.ParseBool(row[2]) // default to false on error
		tasks = append(tasks, Task{
			id:    id,
			value: row[1],
			done:  done,
		})
	}

	return tasks, nil
}

func WriteTask(t Task, path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("create file: %v", err)
		}

		_, err = f.WriteString("id,value,done\n")
		if err != nil {
			f.Close()
			return fmt.Errorf("write file: %v", err)
		}
		f.Close()
	} else if err != nil {
		return fmt.Errorf("status file: %v", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	err = writer.Write([]string{
		strconv.Itoa(t.id),
		t.value,
		strconv.FormatBool(t.done),
	})
	if err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}

func UpdateTask(id int, value string) error {
	tasks, err := GetTasks(TasksPath)
	if err != nil {
		return fmt.Errorf("get tasks: %v", err)
	}

	// find and update task
	updated := false

	for i := range tasks {
		if tasks[i].id == id {
			tasks[i].value = value
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("task with ID: %v not found", id)
	}

	f, err := os.Create(TasksPath)
	if err != nil {
		return fmt.Errorf("overwrite file: %v", err)
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{"id", "value", "done"})

	for _, t := range tasks {
		writer.Write([]string{
			strconv.Itoa(t.id),
			t.value,
			strconv.FormatBool(t.done),
		})
	}

	return nil
}

func CheckTask(id int) error {
	tasks, err := GetTasks(TasksPath)
	if err != nil {
		return fmt.Errorf("get tasks: %v", err)
	}

	updated := false

	for i := range tasks {
		if tasks[i].id == id {
			tasks[i].done = !tasks[i].done
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("task with ID: %v not found", id)
	}

	f, err := os.Create(TasksPath)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{"id", "value", "done"})

	for _, t := range tasks {
		writer.Write([]string{
			strconv.Itoa(t.id),
			t.value,
			strconv.FormatBool(t.done),
		})
	}

	return nil
}

func DeleteTask(id int) error {
	tasks, err := GetTasks(TasksPath)
	if err != nil {
		return fmt.Errorf("get tasks: %v", err)
	}

	taskExists := false

	for i := range tasks {
		if tasks[i].id == id {
			taskExists = true
			break
		}
	}

	if !taskExists {
		return fmt.Errorf("task with ID: %v not found", id)
	}

	tasks = slices.DeleteFunc(tasks, func(t Task) bool {
		return t.id == id
	})

	f, err := os.Create(TasksPath)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{"id", "value", "done"})

	for _, t := range tasks {
		writer.Write([]string{
			strconv.Itoa(t.id),
			t.value,
			strconv.FormatBool(t.done),
		})
	}

	return nil
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			MarginBottom(1)

	taskBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			MarginBottom(1)

	idStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	doneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	pendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
)

func DisplayTasks(tasks []Task) {
	if len(tasks) < 1 {
		fmt.Println("\n  No tasks to display")
		return
	}

	fmt.Println(titleStyle.Render("Tasks"))

	for _, t := range tasks {
		var status string
		if t.done {
			status = doneStyle.Render("✓ Done")
		} else {
			status = pendingStyle.Render("○ Pending")
		}

		content := fmt.Sprintf("%s  %s\n%s",
			idStyle.Render(fmt.Sprintf("#%d", t.id)),
			t.value,
			status,
		)

		fmt.Println(taskBoxStyle.Render(content))
	}
}

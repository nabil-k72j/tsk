# `tsk` — A Simple CLI Task Manager

`tsk` is a lightweight command-line task manager written in Go.
Tasks are stored in a CSV file (`tasks.csv`) in the same directory.
You can add tasks, list them, delete them, and toggle them as done/undone.

---

## ⭐ Features

* Add tasks with a simple command:
  `tsk "Buy groceries"`
* List all tasks or limit the number shown
* Update the value of a task by ID
* Delete tasks by ID
* Toggle completion status
* Persistent CSV storage

---

## 📦 Installation

```sh
go build -o tsk .
```

Move it somewhere in your `$PATH` if you want global access:

```sh
sudo mv tsk /usr/local/bin/
```

---

## 📄 Usage

### ➕ Add a task

Adds a new task with the next available ID.

```sh
tsk "Finish homework"
```

---

### 📋 List tasks

List all:

```sh
tsk --list
```

List the latest N tasks:

```sh
tsk --list 10
```

---

### 📝 Update a task

```sh
tsk --update 3 "new value"
```

Updates the value of task with ID 3

---

### ❌ Delete a task

```sh
tsk --delete 3
```

Deletes the task with ID 3.

---

### ✔ Toggle completion (check/uncheck)

```sh
tsk --check 5
```

If task 5 is done → it becomes undone
If task 5 is undone → it becomes done

---

## 📁 Data Storage

Tasks are stored in:

```
tasks.csv
```

Format:

```
id,value,done
1,Buy groceries,false
2,Walk the dog,true
```

Your code automatically:

* Creates the file if missing
* Writes a header
* Loads tasks into memory
* Rewrites the file on updates/deletes

---

## 🛠 Flags Summary

| Flag       | Description                    | Example          |
| ---------- | ------------------------------ | ---------------- |
| `--list`   | List all tasks or last N tasks | `tsk --list 5`   |
| `--update` | Update the value of task by ID | `tsk --update 5` |
| `--delete` | Delete task by ID              | `tsk --delete 3` |
| `--check`  | Toggle done/undone for task ID | `tsk --check 10` |

If no flags are provided, the positional argument is treated as a new task.

---

## 🧱 Code Structure (optional section)

```
.
├── main.go
├── tasks.go (load/save/add/delete/check)
└── tasks.csv (generated automatically)
```

---

## 📝 Example Session

```
$ tsk "Write README"
Added task #1

$ tsk "Push code to GitHub"
Added task #2

$ tsk --list
1  Write README        Unfinished
2  Push code to GitHub Unfinished

$ tsk --check 1
Task #1 marked

$ tsk --list
1  Write README        Finished
2  Push code to GitHub Unfinished

$ tsk --update 1 "Write README.md"
Updated task #1

$ tsk --Delete 2
Deleted task #2
```

---

## 📜 License

MIT


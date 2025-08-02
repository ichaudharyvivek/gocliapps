# 🛠️ GoTask CLI
Task Manager CLI in Go

## Objective

Create a command-line task manager in Go named `gotask` that allows users to add, list, delete and complete the tasks via subcommands. Tasks should persist between sessions using a local `tasks.json` file.

## Functional Requirements

### Subcommands

#### 1. `add`

```bash
gotask add -task="Write Go CLI project"
```

- Add a new task to the list.
- Accept a `-task` flag (string, required).
- Automatically assign a unique incremental ID.
- Set default status as pending.

#### 2. `list`

```bash
gotask list
```

- List all tasks in the system.
- Display all tasks in tabular format.
- Show `ID`, `Description`, and `Status`.
- Handle case when no tasks exist.

#### 2. `done`

```bash
gotask done -id=2
```

- Mark a task as completed.
- Accept an `-id` flag (int, required).
- Update the task status to done.
- Validate if task with given ID exists.

#### 4. `delete`

```bash
gotask delete -id=3
```

- Delete a task from the list.
- Accept an `-id` flag (int, required).
- Remove the task with the given ID.
- Validate if task with given ID exists.

## Notes

#### Folder structure

When we make a complex cli app, with the use of a CLI package like `cobra` the folder structure we usually follow is this:

```
gotask/
├── cmd/
│   ├── root.go          # main entry point for CLI
│   ├── add.go           # `add` command
│   ├── list.go          # `list` command
│   ├── done.go          # `done` command
│   └── delete.go        # `delete` command
│
├── internal/
│   └── task/
│       ├── store.go     # read/write from JSON
│       ├── model.go     # task struct, types
│       └── manager.go   # business logic: add, delete, etc.
│
├── tasks.json           # data file created at runtime
├── main.go              # just calls cmd.Execute()
├── go.mod
└── README.md
```

However, we are building a simple todo app cli. Hence, we will use a simpler approach:

```
gotask/
├── main.go                  # CLI entry point, handles flag parsing
├── go.mod
├── task/
│   ├── task.go              # Task struct and status constants
│   └── store.go             # Read/write to tasks.json
├── tasks.json               # Created at runtime
└── README.md                # Project overview
```

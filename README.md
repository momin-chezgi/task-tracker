# task-tracker
A CLI command that manages your tasks: what you need to do, what you have done, what you are currently working on, and what you have cancelled.

#### Where does the idea come from?
It is an educational project in Go, based on the projects on roadmap.sh in the [Task Tracker CLI](https://roadmap.sh/projects/task-tracker).

## Installation

### Prerequisites

This project requires **Go 1.27.1 or later**.

Check whether Go is already installed:

```bash
go version
```

If Go is not installed, install it from the official [Go website](https://go.dev/dl/).

### Linux and macOS

Install `task-tracker` with:

```bash
go install github.com/momin-chezgi/task-tracker/cmd/task-tracker@latest
```

Go installs the executable in the directory reported by `go env GOBIN`. If `GOBIN` is not set, it uses `GOPATH/bin` (usually `~/go/bin`). Ensure that directory is in your `PATH`. For Bash or Zsh, if you use the default location, add:

```bash
export PATH="$HOME/go/bin:$PATH"
```

To make this change permanent, add the line above to your `~/.bashrc` or `~/.zshrc`.

You can then run the program from any directory, for example:

```bash
task-tracker list
```

### Windows

Open **PowerShell** and run:

```powershell
go install github.com/momin-chezgi/task-tracker/cmd/task-tracker@latest
```

Go installs the executable in the directory reported by `go env GOBIN`. If `GOBIN` is not set, it uses `GOPATH\bin` (usually `%USERPROFILE%\go\bin`).

To add the default installation directory to your user's `PATH`:

1. Open the Start menu and search for **Edit environment variables for your account**.
2. Select `Path` under **User variables**, then click **Edit**.
3. Click **New** and add:

   ```text
   %USERPROFILE%\go\bin
   ```

4. Save the changes and restart PowerShell.

If you have configured a custom `GOBIN`, add that directory instead.

Run the following command to verify the installation:

```powershell
task-tracker list
```

### Verify the installation

You can check whether the executable is available with:

**Linux and macOS:**

```bash
command -v task-tracker
```

**Windows PowerShell:**

```powershell
Get-Command task-tracker
```

If the command returns the location of the executable, it is available through your `PATH`.

### Update

To update an existing installation to the latest version:

```bash
go install github.com/momin-chezgi/task-tracker/cmd/task-tracker@latest
```

The new executable will replace the previous version.

### Install a specific version

Replace `@latest` with an available version tag:

```bash
go install github.com/momin-chezgi/task-tracker/cmd/task-tracker@v1.0.0
```

Make sure the version tag exists in the repository before using it.

## Data file

Tasks are stored in `data.json` under the operating system's user config directory, inside a `task-tracker` folder. On Linux, this is typically `$XDG_CONFIG_HOME/task-tracker/data.json`, or `~/.config/task-tracker/data.json` when `XDG_CONFIG_HOME` is not set. On Windows, the corresponding location is typically `%AppData%\task-tracker\data.json` when using Go's standard user configuration directory.

The directory is created automatically. This keeps task data available when the command is run from different working directories.

## Options:

* `add`: makes a new task and adds it to the list of tasks. You should pass a string argument for the description of the task.

Example:
```
$task-tracker add "Come up with a jogging routine"
Task added successfully (ID: 2)
```

> Note: Every task has one of these four stages: `done`, `todo`, `in-progress` or `cancelled`.
>
> Every newly added task is a `todo` task if you don't change its status with the `mark` command.

* `list`: you can write this option without any other arguments, or you can filter the tasks by their status.

Example:
```
$task-tracker list
---------------------
Buy groceries     (Done)
ID: 1
Created at: 2026-10-04 08:20:53.261562075 +0330 +0330
Updated at: 2026-10-04 08:20:53.261562141 +0330 +0330
---------------------

---------------------
Come up with a jogging routine  (To-Do)
ID: 2
Created at: 2026-10-04 08:21:09.564909106 +0330 +0330
Updated at: 2026-10-04 08:21:09.564909198 +0330 +0330
---------------------

---------------------
Go to the bank and open an account      (In Progress)
ID: 4
Created at: 2026-10-05 09:23:01.072078 +0330 +0330
Updated at: 2026-10-05 09:23:01.072078 +0330 +0330
---------------------

$task-tracker list done
---------------------
Buy groceries     (Done)
ID: 1
Created at: 2026-10-04 08:20:53.261562075 +0330 +0330
Updated at: 2026-10-04 08:20:53.261562141 +0330 +0330
---------------------

$task-tracker list todo
---------------------
Come up with a jogging routine  (To-Do)
ID: 2
Created at: 2026-10-04 08:21:09.564909106 +0330 +0330
Updated at: 2026-10-04 08:21:09.564909198 +0330 +0330
---------------------
```

* `show`: Shows a single task distinguished by its ID.

Example:
```
$task-tracker show 2
---------------------
Come up with a jogging routine  (To-Do)
ID: 2
Created at: 2026-10-04 08:21:09.564909106 +0330 +0330
Updated at: 2026-10-04 08:21:09.564909198 +0330 +0330
---------------------
```

* `update`: Updates the description of a task distinguished by its ID. At the end, provide the new description string.

Example:
```
$task-tracker show 1
---------------------
Buy groceries     (Done)
ID: 1
Created at: 2026-10-04 08:20:53.261562075 +0330 +0330
Updated at: 2026-10-04 08:20:53.261562141 +0330 +0330
---------------------

$task-tracker update 1 "Buy bread"
$task-tracker show 1
---------------------
Buy bread       (Done)
ID: 1
Created at: 2026-10-04 08:20:53.261562075 +0330 +0330
Updated at: 2026-10-05 10:31:01.9459913 +0330 +0330
---------------------
```

* `cancel`: Cancels a task. Cancelled tasks are not shown in the output of the `list` command.

Example:
```
$task-tracker cancel 2
$task-tracker list
---------------------
Buy bread       (Done)
ID: 1
Created at: 2026-10-04 08:20:53.261562075 +0330 +0330
Updated at: 2026-10-05 10:31:01.9459913 +0330 +0330
---------------------

---------------------
Go to the bank and open an account      (In progress)
ID: 4
Created at: 2026-10-05 09:23:01.072078 +0330 +0330
Updated at: 2026-10-05 09:23:01.072078 +0330 +0330
---------------------

$task-tracker show 2
---------------------
Come up with a jogging routine  (Cancelled)
ID: 2
Created at: 2026-10-04 08:21:09.564909106 +0330 +0330
Updated at: 2026-10-04 08:21:09.564909198 +0330 +0330
---------------------
```

* `mark`: Changes the status of a task. Use one of these commands: `mark-todo`, `mark-in-progress`, or `mark-done`. Running `mark` without a valid subcommand returns an error.

Example:
```
$task-tracker mark-todo 4
$task-tracker list
---------------------
Buy bread       (Done)
ID: 1
Created at: 2026-10-04 08:20:53.261562075 +0330 +0330
Updated at: 2026-10-05 10:31:01.9459913 +0330 +0330
---------------------

---------------------
Go to the bank and open an account      (To-Do)
ID: 4
Created at: 2026-10-05 09:23:01.072078 +0330 +0330
Updated at: 2026-10-05 09:23:01.072078 +0330 +0330
---------------------

$task-tracker mark-done 2
something went wrong: This task has been cancelled!
```

> Note: You can't change the status of a cancelled task.

* `help`
* `version`

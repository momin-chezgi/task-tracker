# task-tracker
A CLI command that manages your tasks: what you need to do, what you have done, what you are currently working on, and what you have cancelled.

#### Where does the idea come from?
It is already an educational project in Go, based on the projects on roadmap.sh in the [Task Tracker CLI](https://roadmap.sh/projects/task-tracker).


## Options:
1. `add`: makes a new task and adds it to the list of tasks. You should pass a string argument for the description of the task.
Example:
```
$task-tracker add "Come up with a jogging routine"
Task added successfully (ID: 2)
```
> Note: Every task has one of these four stages: `done`, `todo`, `in-progress` or cancelled.
> (Every newly added task is a `todo` task if you don't change its status with the `mark` command)
2. `list`: you can just write this option without any other arguments, or you can filter the tasks by their status.
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
3. `show`: Shows a single task distinguished by its ID.
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
4. `update`: Update the description of a task distinguished by its ID. At the end you should give the new description string.
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
5. `cancel`: Cancels a command. This command won't be shown in the output of the `list` command
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
7. `mark`: You can change the status of a task. There are three options(if you enter `mark`, the output includes an error): `mark-todo`, `mark-in-progress`, `mark-done`
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
$task-tracker  mark-done 2 
something went wrong: This task has been cancelled!
```
>Note: you can't change the status of a cancelled task.

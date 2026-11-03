# API Test Storage Refactor Design

## Goal

Update `internal/api/api_test.go` to use the API's `tskStorage` storage pointer
instead of the obsolete `tasks` test state.

## Scope

- Initialize `tskStorage` as a pointer to an empty `store.Storage` in
  `resetTaskStore`.
- Reset `idCounter` along with the storage in test setup and cleanup.
- Replace test setup, mutations, snapshots, and assertions that access `tasks`
  with equivalent access to `tskStorage.Tasks`.
- Preserve existing test cases, direct handler calls, and request URLs.
- Do not modify the API implementation or `cmd/todo-api/main.go`.

## Validation

Run `go test ./internal/api` and `go test ./...` after the test migration.

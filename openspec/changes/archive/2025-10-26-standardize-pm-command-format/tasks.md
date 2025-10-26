## 1. Specification
- [x] 1.1 Finalize cli-command-format requirements and scenarios
- [x] 1.2 Validate the change with `openspec validate standardize-pm-command-format --strict`

## 2. Implementation
- [x] 2.1 Refactor `internal/adapters/handlers/cli/delete` to the shared layout (exported `Command()`, grouped flag constants, helper separation)
- [x] 2.2 Refactor `internal/adapters/handlers/cli/init` to the shared layout and replace the global `InitCmd` with `Command()`
- [x] 2.3 Refactor `internal/adapters/handlers/cli/list` to the shared layout and align flag registration
- [x] 2.4 Refactor `internal/adapters/handlers/cli/new` to close remaining gaps (constructor style, helper structure) with `pm edit`
- [x] 2.5 Extract reusable helpers for shared error handling or flag parsing when duplication appears
- [x] 2.6 Align unit and integration tests with the new command format expectations

## 3. Verification
- [x] 3.1 Add or update lint/build checks that guard the formatting contract
- [x] 3.2 Run `make lint`, `make unit`, and targeted CLI integration tests
- [x] 3.3 Update documentation or contributor guides referencing command implementation patterns

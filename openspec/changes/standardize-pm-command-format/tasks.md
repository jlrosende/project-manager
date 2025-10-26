## 1. Specification
- [ ] 1.1 Finalize cli-command-format requirements and scenarios
- [ ] 1.2 Validate the change with `openspec validate standardize-pm-command-format --strict`

## 2. Implementation
- [ ] 2.1 Refactor `internal/adapters/handlers/cli/delete` to the shared layout (exported `Command()`, grouped flag constants, helper separation)
- [ ] 2.2 Refactor `internal/adapters/handlers/cli/init` to the shared layout and replace the global `InitCmd` with `Command()`
- [ ] 2.3 Refactor `internal/adapters/handlers/cli/list` to the shared layout and align flag registration
- [ ] 2.4 Refactor `internal/adapters/handlers/cli/new` to close remaining gaps (constructor style, helper structure) with `pm edit`
- [ ] 2.5 Extract reusable helpers for shared error handling or flag parsing when duplication appears
- [ ] 2.6 Align unit and integration tests with the new command format expectations

## 3. Verification
- [ ] 3.1 Add or update lint/build checks that guard the formatting contract
- [ ] 3.2 Run `make lint`, `make unit`, and targeted CLI integration tests
- [ ] 3.3 Update documentation or contributor guides referencing command implementation patterns

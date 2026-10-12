# Agent Guidelines

## TDD Workflow

All development **must** follow a strict Red → Green → Refactor cycle. Never write production code before a failing test exists.

### Steps (in order, no skipping)

1. **Red** — Write a failing test that captures the intended behaviour. Run the tests and confirm the new test fails (and only the new test fails). **Stop and ask the user to review the failing test before proceeding.**

2. **Green** — Write the *minimum* code required to make the failing test pass. No gold-plating. Run the tests and confirm all tests are green. **Stop and ask the user to review the implementation before proceeding.**

3. **Refactor** — Clean up both the production code and the test code (remove duplication, improve naming, simplify logic). Tests must remain green throughout. **Stop and ask the user to review the refactored code before considering the task done.**

> At every stage, wait for explicit user approval before moving to the next step.

---

## Testing Strategy by Package

The project is organized by application contexts under `app/` (one app per context, e.g. `inventory`) and shared reusable packages under `pkg/` (e.g. `database`, `domain`, `message-relay`). Each layer has a dedicated testing approach:

### Coverage Enforcement for Existing Code
- When adding or updating tests for existing production code, run the relevant test command with coverage enabled (e.g. `go test ./... -cover` or a package-scoped equivalent with `-cover`).
- 100% coverage remains mandatory for the impacted `command/` and `projection/` packages. Do not consider the task done until coverage confirms this threshold.

### `pkg/`
- **Unit tests are mandatory for shared libraries under `pkg/`.**
- Test files live alongside the production code (e.g. `pkg/message-relay/source/binlog/redis/state_test.go`).
- Mock external dependencies in unit tests so they do not require running databases, caches, or other services.
- Integration tests are optional for libraries that integrate with databases, caches, or other persistence systems. When added, they exercise the real (or containerised) dependency and complement, rather than replace, unit tests.
- These rules take precedence over the `domain/` and `database/` guidance below for packages under `pkg/`.

### `command/`
- **Coverage target: 100% unit tests.**
- These packages orchestrate domain logic. Test them with unit tests only.
- **Do not mock the `domain/` layer.** Use real domain objects so the unit tests implicitly validate domain behaviour too.
- Test files live alongside the production code (e.g. `command/create_vehicle_test.go`).

### `projection/`
- **Coverage target: 100% unit tests.**
- Projections handle the read side of the CQRS model and do not orchestrate domain logic.
- Test them with unit tests only.
- Test files live alongside the production code (e.g. `projection/get_vehicle_test.go`).

### `domain/`
- No dedicated test suite required.
- Domain behaviour is verified implicitly through the `command/` and `projection/` unit tests.
- Add domain-level tests only when a piece of logic cannot be reached from the layers above.

### `database/`
- Holds persistence implementations (repositories, etc.).
- **Covered by integration tests** that exercise a real (or containerised) database.
- Unit-test only when a piece of logic is clearly independent of the database (e.g. pure query-building helpers).

### `handler/`
- Holds HTTP controllers/handlers.
- **Covered by unit tests** that mock the underlying use-case/command/projection dependencies.
- E2E tests may be added in the future but are not required today.
- Test files live alongside the production code (e.g. `handler/create_vehicle_handler_test.go`).

## Integration Tests Policy

- Integration tests are mandatory for integration behavior under `app/`, including but not limited to database persistence and Redis-backed state.
- Under `pkg/`, integration tests are optional for libraries that integrate with databases, caches, or other persistence systems; unit tests remain mandatory.
- Integration tests live in:
  - `app/{context}/test/integration` (for example, `app/inventory/test/integration`)
  - `pkg/{name}/test/integration` (for example, `pkg/database/test/integration`, `pkg/domain/test/integration`, `pkg/message-relay/test/integration`)
- Coverage target for integration tests is **100%** for the impacted integration package/scope.
- When adding or updating integration tests, run them with coverage enabled and do not consider the task complete until the required 100% coverage is confirmed.

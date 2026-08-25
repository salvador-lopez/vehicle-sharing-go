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

The project follows a layered architecture under `internal/vehicle/`. Each layer has a dedicated testing approach:

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
- Integration tests live in `test/integration/`.
- Unit-test only when a piece of logic is clearly independent of the database (e.g. pure query-building helpers).

### `handler/`
- Holds HTTP controllers/handlers.
- **Covered by unit tests** that mock the underlying use-case/command/projection dependencies.
- E2E tests may be added in the future but are not required today.
- Test files live alongside the production code (e.g. `handler/create_vehicle_handler_test.go`).

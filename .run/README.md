# IDE Run Configurations

## Integration test configs

- `integration-test-inventory` (folder: `inventory`) — runs all integration
  tests in the `app/inventory` Go module (`-tags=integration`).
- `integration-test-pkg` (folder: `pkg`) — runs all integration tests in the
  `pkg` Go module (`-tags=integration`). Separate config because `pkg` is a
  distinct Go module and can't be covered by a single `go test` invocation
  rooted in `app/inventory`.
- `integration-test` (Compound, project root) — runs both of the above
  together, each in its own tab (JetBrains Compound configs run children in
  parallel, not merged into one process). Supports plain Run, Debug, and
  Run with Coverage — each child tab shows its own results/coverage.

### Notes

- Coverage requires `-coverpkg=./...` in each child's Go tool arguments,
  otherwise coverage is (incorrectly) measured only against the test's own
  package instead of the production code it exercises, showing as 0%.
- Debugging via the compound config debugs each child process
  independently; set breakpoints in the relevant module's code as usual.

Both configs require the local MySQL (`inventory_test` DB) and Redis
containers to be up (`docker compose up -d` from the repo root).

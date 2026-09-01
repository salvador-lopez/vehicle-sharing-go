# IDE Run Configurations

## Integration test configs

- `integration-test-inventory` (folder: `inventory`) — runs all integration
  tests in the `app/inventory` Go module (`-tags=integration`).
- `integration-test-pkg` (folder: `pkg`) — runs all integration tests in the
  `pkg` Go module (`-tags=integration`). Separate config because `pkg` is a
  distinct Go module and can't be covered by a single `go test` invocation
  rooted in `app/inventory`.
- `integration-test` (Compound, project root) — runs both of the above
  together for a quick full pass/fail across both modules.

### Known limitation: Compound + Debug/Coverage

JetBrains Compound run configurations start all child configurations, but
**do not properly support "Debug" or "Run with Coverage"** — only one child
process gets instrumented/debugged, so you'll see 0% coverage or only one
child appearing to execute.

**Use `integration-test` (compound) only for plain Run.**

For **Debug** or **Run with Coverage**, run `integration-test-inventory` and
`integration-test-pkg` **individually** — each is a standalone Go Test
configuration and fully supports breakpoints and coverage reporting.

Both configs require the local MySQL (`inventory_test` DB) and Redis
containers to be up (`docker compose up -d` from the repo root).

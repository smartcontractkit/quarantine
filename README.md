# Quarantine Flaky Tests

Quarantine and gate execution of flaky tests with a simple package.

```go
import "github.com/smartcontractkit/quarantine"

func TestFlaky(t *testing.T) {
    quarantine.Flaky(t, "TICKET-Number")
    // Rest of test
}
```

Use the env var `RUN_QUARANTINED_TESTS = "true"` to run these tests. All flaky tests get some info logged to `t.Log`, and also log `"flaky_test":"<TICKET-Number>"` with [t.Attr](https://pkg.go.dev/testing#T.Attr) for easy output parsing.
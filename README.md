# Go Once Benchmark

Benchmark various ways to ensure that a task is done only once:

- Using `sync.Once`
- Using `sync.Mutex`
- Using `sync.RWMutex`
- Using `chan`

All methods have a `Precheck` variant that checks whether the task is done before
actually locking.

We benchmark those methods in two scenarios:

- When the task is not previously done: in this scenario, we do not run the `Precheck` variant, as prechecking always
  return false unless you are running this benchmark on a potato
- When the task is already completed: in this case, we expect the `Precheck` variant to be faster.

## Benchmark

Configuration can be changed in the `scenario_test.go` file:

```go
package scenario_test

import (
	"time"
)

const (
	// DoStuffTook controls how long doStuff takes.
	// You should subtract this duration from the benchmark results when running BenchmarkStandardScenario.
	DoStuffTook = 2 * time.Second
	// GoRoutinesCount controls how many goroutines to create.
	GoRoutinesCount = 10_000
)
```

Run the benchmark:

```shell
go test -bench . -count=3
```

The results on my machine (with chrome and some stuff running in the background):

```shell
goos: windows
goarch: amd64
pkg: scenario
cpu: AMD Ryzen 9 7900 12-Core Processor
BenchmarkStandardScenario/OncePrecheck-24                      1        2004853800 ns/op
BenchmarkStandardScenario/OncePrecheck-24                      1        2005557700 ns/op
BenchmarkStandardScenario/OncePrecheck-24                      1        2005833900 ns/op
BenchmarkStandardScenario/Lock-24                              1        2005363400 ns/op
BenchmarkStandardScenario/Lock-24                              1        2005569500 ns/op
BenchmarkStandardScenario/Lock-24                              1        2005285100 ns/op
BenchmarkStandardScenario/RWLock-24                            1        2001862000 ns/op
BenchmarkStandardScenario/RWLock-24                            1        2001504600 ns/op
BenchmarkStandardScenario/RWLock-24                            1        2001793900 ns/op
BenchmarkStandardScenario/Channel-24                           1        2004951000 ns/op
BenchmarkStandardScenario/Channel-24                           1        2004595800 ns/op
BenchmarkStandardScenario/Channel-24                           1        2004771200 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                   962           1198190 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                  1008           1204863 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                  1027           1169837 ns/op
BenchmarkDoneScenario/Done_Lock-24                           882           1248253 ns/op
BenchmarkDoneScenario/Done_Lock-24                           998           1258346 ns/op
BenchmarkDoneScenario/Done_Lock-24                           903           1267249 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                  1040           1167875 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                   937           1198665 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                  1010           1177499 ns/op
BenchmarkDoneScenario/Done_RWLock-24                        1023           1195476 ns/op
BenchmarkDoneScenario/Done_RWLock-24                        1021           1217294 ns/op
BenchmarkDoneScenario/Done_RWLock-24                        1014           1195385 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                1034           1188411 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                1002           1183191 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                1033           1217335 ns/op
BenchmarkDoneScenario/Done_Channel-24                        254           4716841 ns/op
BenchmarkDoneScenario/Done_Channel-24                        252           4830221 ns/op
BenchmarkDoneScenario/Done_Channel-24                        242           4826638 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24               1017           1183201 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24               1026           1172154 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24                966           1202709 ns/op
PASS
ok      scenario        49.370s
```

The results shown that the `sync.RWMutex` gives the best throughput: for 10,000 goroutines 
it is about 3ms (or about 60%) faster than other methods.

However, the `sync.RWMutex` method is also more complex than the other methods.
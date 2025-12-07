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
BenchmarkStandardScenario/Once-24                      1        2004738100 ns/op
BenchmarkStandardScenario/Once-24                      1        2006183300 ns/op
BenchmarkStandardScenario/Once-24                      1        2005251100 ns/op
BenchmarkStandardScenario/Lock-24                      1        2005081500 ns/op
BenchmarkStandardScenario/Lock-24                      1        2005319800 ns/op
BenchmarkStandardScenario/Lock-24                      1        2005249200 ns/op
BenchmarkStandardScenario/RWLock-24                    1        2002888600 ns/op
BenchmarkStandardScenario/RWLock-24                    1        2002059100 ns/op
BenchmarkStandardScenario/RWLock-24                    1        2002008900 ns/op
BenchmarkStandardScenario/Channel-24                   1        2004769000 ns/op
BenchmarkStandardScenario/Channel-24                   1        2004357600 ns/op
BenchmarkStandardScenario/Channel-24                   1        2004456700 ns/op
BenchmarkStandardScenario/AtomicSwap-24                1        1999966100 ns/op
BenchmarkStandardScenario/AtomicSwap-24                1        2000061400 ns/op
BenchmarkStandardScenario/AtomicSwap-24                1        2000888800 ns/op
BenchmarkStandardScenario/AtomicCAS-24                 1        2000488700 ns/op
BenchmarkStandardScenario/AtomicCAS-24                 1        2000083600 ns/op
BenchmarkStandardScenario/AtomicCAS-24                 1        2000785200 ns/op
BenchmarkDoneScenario/Done_Once-24                  1051           1131385 ns/op
BenchmarkDoneScenario/Done_Once-24                  1059           1131741 ns/op
BenchmarkDoneScenario/Done_Once-24                  1062           1123734 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                  1084           1117444 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                   996           1252450 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                  1028           1131280 ns/op
BenchmarkDoneScenario/Done_Lock-24                          1015           1199054 ns/op
BenchmarkDoneScenario/Done_Lock-24                          1021           1198459 ns/op
BenchmarkDoneScenario/Done_Lock-24                           993           1193328 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                  1041           1132922 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                  1063           1124003 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                  1058           1192098 ns/op
BenchmarkDoneScenario/Done_RWLock-24                         978           1161744 ns/op
BenchmarkDoneScenario/Done_RWLock-24                        1060           1201274 ns/op
BenchmarkDoneScenario/Done_RWLock-24                         986           1207492 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                 920           1179833 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                 866           1208491 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                1078           1128873 ns/op
BenchmarkDoneScenario/Done_Channel-24                        258           4661335 ns/op
BenchmarkDoneScenario/Done_Channel-24                        254           5243544 ns/op
BenchmarkDoneScenario/Done_Channel-24                        240           5015273 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24                951           1174689 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24                961           1145474 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24               1054           1249393 ns/op
BenchmarkDoneScenario/Done_AtomicSwap-24                    1052           1134125 ns/op
BenchmarkDoneScenario/Done_AtomicSwap-24                    1036           1159636 ns/op
BenchmarkDoneScenario/Done_AtomicSwap-24                    1065           1142673 ns/op
BenchmarkDoneScenario/Done_AtomicSwapPrecheck-24            1062           1130601 ns/op
BenchmarkDoneScenario/Done_AtomicSwapPrecheck-24            1018           1125005 ns/op
BenchmarkDoneScenario/Done_AtomicSwapPrecheck-24            1081           1132058 ns/op
BenchmarkDoneScenario/Done_AtomicCAS-24                     1064           1154833 ns/op
BenchmarkDoneScenario/Done_AtomicCAS-24                     1045           1160291 ns/op
BenchmarkDoneScenario/Done_AtomicCAS-24                     1066           1237422 ns/op
BenchmarkDoneScenario/Done_AtomicCASPrecheck-24             1063           1131400 ns/op
BenchmarkDoneScenario/Done_AtomicCASPrecheck-24             1052           1141671 ns/op
BenchmarkDoneScenario/Done_AtomicCASPrecheck-24             1044           1129641 ns/op
PASS
ok      scenario        79.546s
```

The results showed that the `sync.RWMutex` gives the best throughput: for 10,000 goroutines 
it is about 3ms (or about 60%) faster than other methods.

However, the `sync.RWMutex` method is also more complex than the other methods.
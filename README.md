# Go Once Benchmark

Benchmark various ways to ensure that a task is done only once:

- Using `sync.Once`
- Using `sync.Mutex`
- Using `sync.RWMutex`
- Using `chan`
- Using `atomic.Bool`

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
BenchmarkStandardScenario/Once-24                      1        2004703900 ns/op
BenchmarkStandardScenario/Once-24                      1        2005059000 ns/op
BenchmarkStandardScenario/Once-24                      1        2005643200 ns/op
BenchmarkStandardScenario/Lock-24                      1        2005060900 ns/op
BenchmarkStandardScenario/Lock-24                      1        2004998600 ns/op
BenchmarkStandardScenario/Lock-24                      1        2005020100 ns/op
BenchmarkStandardScenario/RWLock-24                    1        2002263800 ns/op
BenchmarkStandardScenario/RWLock-24                    1        2002042000 ns/op
BenchmarkStandardScenario/RWLock-24                    1        2001978100 ns/op
BenchmarkStandardScenario/Channel-24                   1        2004577800 ns/op
BenchmarkStandardScenario/Channel-24                   1        2004936200 ns/op
BenchmarkStandardScenario/Channel-24                   1        2005001000 ns/op
BenchmarkStandardScenario/AtomicSwap-24                1        2000423900 ns/op
BenchmarkStandardScenario/AtomicSwap-24                1        2000185500 ns/op
BenchmarkStandardScenario/AtomicSwap-24                1        2000109700 ns/op
BenchmarkStandardScenario/AtomicCAS-24                 1        2000361500 ns/op
BenchmarkStandardScenario/AtomicCAS-24                 1        2000223700 ns/op
BenchmarkStandardScenario/AtomicCAS-24                 1        2000147300 ns/op
BenchmarkDoneScenario/Done_Once-24                  1054           1127393 ns/op
BenchmarkDoneScenario/Done_Once-24                  1070           1134143 ns/op
BenchmarkDoneScenario/Done_Once-24                  1082           1123561 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                  1057           1118513 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                  1040           1126664 ns/op
BenchmarkDoneScenario/Done_OncePrecheck-24                  1051           1137375 ns/op
BenchmarkDoneScenario/Done_Lock-24                          1012           1200171 ns/op
BenchmarkDoneScenario/Done_Lock-24                           915           1274216 ns/op
BenchmarkDoneScenario/Done_Lock-24                           896           1225920 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                  1066           1131149 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                  1053           1127758 ns/op
BenchmarkDoneScenario/Done_LockPrecheck-24                  1047           1127130 ns/op
BenchmarkDoneScenario/Done_RWLock-24                        1032           1157363 ns/op
BenchmarkDoneScenario/Done_RWLock-24                         974           1164877 ns/op
BenchmarkDoneScenario/Done_RWLock-24                        1058           1145716 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                1070           1127320 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                1047           1130950 ns/op
BenchmarkDoneScenario/Done_RWLockPreCheck-24                1062           1128635 ns/op
BenchmarkDoneScenario/Done_Channel-24                        256           4586462 ns/op
BenchmarkDoneScenario/Done_Channel-24                        261           4684500 ns/op
BenchmarkDoneScenario/Done_Channel-24                        256           4643004 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24               1053           1129189 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24               1046           1138408 ns/op
BenchmarkDoneScenario/Done_ChannelPrecheck-24               1039           1123737 ns/op
BenchmarkDoneScenario/Done_AtomicSwap-24                    1045           1126974 ns/op
BenchmarkDoneScenario/Done_AtomicSwap-24                    1058           1144896 ns/op
BenchmarkDoneScenario/Done_AtomicSwap-24                    1052           1159469 ns/op
BenchmarkDoneScenario/Done_AtomicSwapPrecheck-24            1026           1151909 ns/op
BenchmarkDoneScenario/Done_AtomicSwapPrecheck-24            1052           1138924 ns/op
BenchmarkDoneScenario/Done_AtomicSwapPrecheck-24            1051           1174992 ns/op
BenchmarkDoneScenario/Done_AtomicCAS-24                     1040           1153370 ns/op
BenchmarkDoneScenario/Done_AtomicCAS-24                     1056           1152302 ns/op
BenchmarkDoneScenario/Done_AtomicCAS-24                     1046           1160185 ns/op
BenchmarkDoneScenario/Done_AtomicCASPrecheck-24             1052           1145626 ns/op
BenchmarkDoneScenario/Done_AtomicCASPrecheck-24             1072           1126586 ns/op
BenchmarkDoneScenario/Done_AtomicCASPrecheck-24             1041           1138697 ns/op
PASS
ok      scenario        79.273s
```

The results showed that the `sync.RWMutex` gives the best throughput: for 10,000 goroutines 
it is about 3ms (or about 60%) faster than other methods.

However, the `sync.RWMutex` method is also more complex than the other methods.
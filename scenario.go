package scenario

import (
	"sync"
	"sync/atomic"
	"time"
)

type Scenario struct {
	GoRoutinesCount int
	DoStuffTook     time.Duration
	CheckTook       time.Duration
	IsDone          bool
	Touched         int
	wg              sync.WaitGroup
}

func (s *Scenario) doStuff() {
	if s.DoStuffTook > 0 {
		time.Sleep(s.DoStuffTook)
	}
	s.IsDone = true
	s.Touched++
}

func (s *Scenario) check() bool {
	if s.CheckTook > 0 {
		time.Sleep(s.CheckTook)
	}
	return s.IsDone
}

func (s *Scenario) Verify() {
	s.wg.Wait()
	if s.Touched > 1 {
		panic("Touched more than once")
	}
	if s.Touched < 1 {
		panic("did not touch")
	}
}

func (s *Scenario) VerifyAndReset() {
	s.Verify()
	s.IsDone = false
	s.Touched = 0
}

// Lock simple use of Lock.
func (s *Scenario) Lock() {
	lock := &sync.Mutex{}
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			lock.Lock()
			defer lock.Unlock()
			if s.check() {
				return
			}
			s.doStuff()
		}()
	}
}

// LockPrecheck pre-check the IsDone flag before acquiring the lock.
func (s *Scenario) LockPrecheck() {
	lock := &sync.Mutex{}
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if s.check() {
				return
			}
			lock.Lock()
			defer lock.Unlock()
			if s.check() {
				return
			}
			s.doStuff()
		}()
	}
}

// RWLock use RLock to check the IsDone flag before locking.
func (s *Scenario) RWLock() {
	lock := &sync.RWMutex{}
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			lock.RLock()
			if s.check() {
				lock.RUnlock()
				return
			}
			lock.RUnlock()

			lock.Lock()
			defer lock.Unlock()
			// Check again, because another goroutine may have changed the IsDone flag.
			if s.check() {
				return
			}
			s.doStuff()
		}()
	}
}

// RWLockPreCheck pre-check the IsDone flag before using RLock to check the IsDone flag.
func (s *Scenario) RWLockPreCheck() {
	lock := &sync.RWMutex{}
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			// Check without locking.
			if s.check() {
				return
			}
			lock.RLock()
			if s.check() {
				lock.RUnlock()
				return
			}
			lock.RUnlock()

			lock.Lock()
			defer lock.Unlock()
			// Check again, because another goroutine may have changed the IsDone flag.
			if s.check() {
				return
			}
			s.doStuff()
		}()
	}
}

func (s *Scenario) OncePrecheck() {
	once := &sync.Once{}
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if s.check() {
				return
			}
			once.Do(s.doStuff)
		}()
	}
}

func (s *Scenario) Once(once *sync.Once) {
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			once.Do(s.doStuff)
		}()
	}
}

func (s *Scenario) Channel() {
	lock := make(chan struct{}, 1)
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			lock <- struct{}{}
			defer func() {
				<-lock
			}()
			if s.check() {
				return
			}
			s.doStuff()
		}()
	}
}

func (s *Scenario) ChannelPrecheck() {
	lock := make(chan struct{}, 1)
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if s.check() {
				return
			}
			lock <- struct{}{}
			defer func() {
				<-lock
			}()
			if s.check() {
				return
			}
			s.doStuff()
		}()
	}
}

func (s *Scenario) AtomicSwap(isOk *atomic.Bool) {
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			alreadyRun := isOk.Swap(true)
			if alreadyRun {
				return
			}
			s.doStuff()
		}()
	}
}

func (s *Scenario) AtomicSwapPrecheck() {
	isOk := atomic.Bool{}
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if s.check() {
				return
			}
			alreadyRun := isOk.Swap(true)
			if alreadyRun {
				return
			}
			s.doStuff()
		}()
	}
}

func (s *Scenario) AtomicCAS(isOk *atomic.Bool) {
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			canRun := isOk.CompareAndSwap(false, true)
			if !canRun {
				return
			}
			s.doStuff()
		}()
	}
}

func (s *Scenario) AtomicCASPrecheck() {
	isOk := atomic.Bool{}
	for i := 0; i < s.GoRoutinesCount; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if s.check() {
				return
			}
			canRun := isOk.CompareAndSwap(false, true)
			if !canRun {
				return
			}
			s.doStuff()
		}()
	}
}

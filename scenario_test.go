package scenario_test

import (
	"scenario"
	"testing"
	"time"
)

const (
	// DoStuffTook controls how long doStuff takes.
	// You should subtract this duration from the benchmark results when running BenchmarkStandardScenario.
	DoStuffTook = 2 * time.Second
	// GoRoutinesCount controls how many goroutines to create.
	GoRoutinesCount = 10_000
)

func BenchmarkStandardScenario(b *testing.B) {
	b.Run("OncePrecheck", func(b *testing.B) {
		s := newScenario()
		for b.Loop() {
			s.OncePrecheck()
			s.VerifyAndReset()
		}
	})

	b.Run("Lock", func(b *testing.B) {
		s := newScenario()
		for b.Loop() {
			s.Lock()
			s.VerifyAndReset()
		}
	})

	b.Run("RWLock", func(b *testing.B) {
		s := newScenario()
		for b.Loop() {
			s.RWLock()
			s.VerifyAndReset()
		}
	})

	b.Run("Channel", func(b *testing.B) {
		s := newScenario()
		for b.Loop() {
			s.Channel()
			s.VerifyAndReset()
		}
	})
}

func BenchmarkDoneScenario(b *testing.B) {
	b.Run("Done OncePrecheck", func(b *testing.B) {
		for b.Loop() {
			s := newDoneScenario()
			s.OncePrecheck()
			s.Verify()
		}
	})

	b.Run("Done Lock", func(b *testing.B) {
		for b.Loop() {
			s := newDoneScenario()
			s.Lock()
			s.Verify()
		}
	})

	b.Run("Done LockPrecheck", func(b *testing.B) {
		for b.Loop() {
			s := newDoneScenario()
			s.LockPrecheck()
			s.Verify()
		}
	})

	b.Run("Done RWLock", func(b *testing.B) {
		for b.Loop() {
			s := newDoneScenario()
			s.RWLock()
			s.Verify()
		}
	})

	b.Run("Done RWLockPreCheck", func(b *testing.B) {
		for b.Loop() {
			s := newDoneScenario()
			s.RWLockPreCheck()
			s.Verify()
		}
	})

	b.Run("Done Channel", func(b *testing.B) {
		for b.Loop() {
			s := newDoneScenario()
			s.Channel()
			s.Verify()
		}
	})

	b.Run("Done ChannelPrecheck", func(b *testing.B) {
		s := newDoneScenario()
		for b.Loop() {
			s.ChannelPrecheck()
			s.Verify()
		}
	})
}

func newScenario() scenario.Scenario {
	return scenario.Scenario{
		GoRoutinesCount: GoRoutinesCount,
		DoStuffTook:     DoStuffTook,
	}
}

func newDoneScenario() scenario.Scenario {
	return scenario.Scenario{
		GoRoutinesCount: GoRoutinesCount,
		IsDone:          true,
	}
}

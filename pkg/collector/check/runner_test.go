package check

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRunner(t *testing.T) {
	r := NewRunner(1)
	assert.NotNil(t, r.pending)
	assert.NotNil(t, r.runningChecks)
}

func TestStop(t *testing.T) {
	r := NewRunner(1)
	r.Stop()
	_, ok := <-r.pending
	assert.False(t, ok)

	// calling Stop on a stopped runner should be a noop
	r.Stop()
}

func TestGetChan(t *testing.T) {
	r := NewRunner(1)
	assert.NotNil(t, r.GetChan())
}

func TestWork(t *testing.T) {
	r := NewRunner(1)
	c1 := TestCheck{}
	c2 := TestCheck{doErr: true}

	r.pending <- &c1
	r.pending <- &c2
	assert.True(t, c1.hasRun)
	r.Stop()

	// fake a check is already running
	r = NewRunner(1)
	c3 := new(TestCheck)
	r.runningChecks[c3.ID()] = c3
	r.pending <- c3
	// wait to be sure the worker tried to run the check
	time.Sleep(100 * time.Millisecond)
	assert.False(t, c3.hasRun)
}

type TimingoutCheck struct {
	TestCheck
}

func (tc *TimingoutCheck) Stop() {
	for {
	}
}

func TestStopCheck(t *testing.T) {
	r := NewRunner(1)
	err := r.StopCheck("foo")
	assert.Nil(t, err)

	c1 := &TestCheck{}
	r.runningChecks[c1.ID()] = c1
	err = r.StopCheck(c1.ID())
	assert.Nil(t, err)

	c2 := &TimingoutCheck{}
	r.runningChecks[c2.ID()] = c2
	err = r.StopCheck(c2.ID())
	assert.Equal(t, "timeout during stop operation on check id TestCheck", err.Error())
}

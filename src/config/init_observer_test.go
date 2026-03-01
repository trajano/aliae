package config

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testInitObserver struct {
	started []InitPhase
	ended   []InitPhase
	errs    []error
}

func (t *testInitObserver) OnInitPhaseStart(phase InitPhase) {
	t.started = append(t.started, phase)
}

func (t *testInitObserver) OnInitPhaseEnd(phase InitPhase, _ time.Duration, err error) {
	t.ended = append(t.ended, phase)
	t.errs = append(t.errs, err)
}

func TestRunInitPhaseNotifiesObserverOnSuccess(t *testing.T) {
	observer := &testInitObserver{}
	phase := InitPhaseRenderEnv

	err := runInitPhase(observer, phase, func() error { return nil })
	require.NoError(t, err)

	assert.Equal(t, []InitPhase{phase}, observer.started)
	assert.Equal(t, []InitPhase{phase}, observer.ended)
	assert.Len(t, observer.errs, 1)
	assert.NoError(t, observer.errs[0])
}

func TestRunInitPhaseNotifiesObserverOnError(t *testing.T) {
	observer := &testInitObserver{}
	phase := InitPhaseRenderScript
	expected := errors.New("boom")

	err := runInitPhase(observer, phase, func() error { return expected })
	require.ErrorIs(t, err, expected)

	assert.Equal(t, []InitPhase{phase}, observer.started)
	assert.Equal(t, []InitPhase{phase}, observer.ended)
	assert.Len(t, observer.errs, 1)
	assert.ErrorIs(t, observer.errs[0], expected)
}

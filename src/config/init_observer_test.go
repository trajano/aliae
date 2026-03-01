package config

import (
	"errors"
	"testing"
	"time"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testInitObserver struct {
	started []InitPhase
	ended   []InitPhase
	errs    []error
	visits  []string
}

func (t *testInitObserver) OnInitPhaseStart(phase InitPhase) {
	t.started = append(t.started, phase)
}

func (t *testInitObserver) OnInitPhaseEnd(phase InitPhase, _ time.Duration, err error) {
	t.ended = append(t.ended, phase)
	t.errs = append(t.errs, err)
}

func (t *testInitObserver) OnInitVisitStart(section InitSection, key string) {
	t.visits = append(t.visits, "start:"+string(section)+":"+key)
}

func (t *testInitObserver) OnInitVisitEnd(section InitSection, key string, _ time.Duration) {
	t.visits = append(t.visits, "end:"+string(section)+":"+key)
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

func TestEmitInitVisitsUsesWalkerOrder(t *testing.T) {
	observer := &testInitObserver{}
	cfg := &Aliae{
		Extends: []Extend{{Path: "./base.yaml"}},
		Vars:    Vars{{Name: "V"}},
		Envs:    shell.Envs{{Name: "E"}},
		Paths:   shell.Paths{{Value: "/bin"}},
		CDPaths: shell.CDPaths{{Value: "/tmp"}},
		Aliae:   shell.Aliae{{Name: "a"}},
		Links:   shell.Links{{Name: "l"}},
		Scripts: shell.Scripts{{Value: "echo hi"}},
	}

	emitInitVisits(cfg, observer)

	assert.Equal(t, []string{
		"start:extends:./base.yaml", "end:extends:./base.yaml",
		"start:var:V", "end:var:V",
		"start:env:E", "end:env:E",
		"start:path:/bin", "end:path:/bin",
		"start:cdpath:/tmp", "end:cdpath:/tmp",
		"start:alias:a", "end:alias:a",
		"start:link:l", "end:link:l",
		"start:script:echo hi", "end:script:echo hi",
	}, observer.visits)
}

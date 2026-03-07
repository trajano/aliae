package init

import (
	"errors"
	"testing"
	"time"

	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testObserver struct {
	started []Phase
	ended   []Phase
	errs    []error
	visits  []string
}

func (t *testObserver) OnPhaseStart(phase Phase) {
	t.started = append(t.started, phase)
}

func (t *testObserver) OnPhaseEnd(phase Phase, _ time.Duration, err error) {
	t.ended = append(t.ended, phase)
	t.errs = append(t.errs, err)
}

func (t *testObserver) OnVisitStart(section Section, key string) {
	t.visits = append(t.visits, "start:"+string(section)+":"+key)
}

func (t *testObserver) OnVisitEnd(section Section, key string, _ time.Duration) {
	t.visits = append(t.visits, "end:"+string(section)+":"+key)
}

func TestRunPhaseNotifiesObserverOnSuccess(t *testing.T) {
	observer := &testObserver{}
	phase := PhaseRenderEnv

	err := runPhase(observer, phase, func() error { return nil })
	require.NoError(t, err)

	assert.Equal(t, []Phase{phase}, observer.started)
	assert.Equal(t, []Phase{phase}, observer.ended)
	assert.Len(t, observer.errs, 1)
	assert.NoError(t, observer.errs[0])
}

func TestRunPhaseNotifiesObserverOnError(t *testing.T) {
	observer := &testObserver{}
	phase := PhaseRenderScript
	expected := errors.New("boom")

	err := runPhase(observer, phase, func() error { return expected })
	require.ErrorIs(t, err, expected)

	assert.Equal(t, []Phase{phase}, observer.started)
	assert.Equal(t, []Phase{phase}, observer.ended)
	assert.Len(t, observer.errs, 1)
	assert.ErrorIs(t, observer.errs[0], expected)
}

func TestEmitVisitsUsesWalkerOrder(t *testing.T) {
	observer := &testObserver{}
	aliae := &cfg.Aliae{
		Extends: []cfg.Extend{{Path: "./base.yaml"}},
		Vars:    cfg.Vars{{Name: "V"}},
		Envs:    shell.Envs{{Name: "E"}},
		Paths:   shell.Paths{{Value: "/bin"}},
		CDPaths: shell.CDPaths{{Value: "/tmp"}},
		Aliae:   shell.Aliae{{Name: "a"}},
		Links:   shell.Links{{Name: "l"}},
		Scripts: shell.Scripts{{Value: "echo hi"}},
	}

	emitVisits(aliae, observer)

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

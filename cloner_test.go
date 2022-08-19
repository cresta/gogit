package gogit

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

const readOnlyPublicRepo = "git@github.com:cresta/gitops-autobot-reference.git"

type TestLogger struct {
	t *testing.T
}

func (t *TestLogger) Debug(_ context.Context, _ string, _ map[string]string, _ map[string]int64) {
}

func (t *TestLogger) Info(_ context.Context, msg string, strings map[string]string, ints map[string]int64) {
	t.t.Logf("%s %s", msg, strings)
}

var _ Logger = &TestLogger{}

func makeTestLogger(t *testing.T) Logger {
	return &TestLogger{
		t: t,
	}
}

func MakeTestCloner(t *testing.T) Cloner {
	return Cloner{
		Logger: makeTestLogger(t),
	}
}

func TestCloner_Clone(t *testing.T) {
	c := MakeTestCloner(t)
	repo, err := c.Clone(context.Background(), readOnlyPublicRepo)
	require.NoError(t, err)
	defer cleanupRepository(t, repo)
	require.NotNil(t, repo)
	require.NotEmpty(t, repo.Location())
	s, err := os.Stat(filepath.Join(repo.Location(), ".git"))
	require.NoError(t, err)
	require.True(t, s.IsDir())
}

func TestCloner_CloneInto(t *testing.T) {
	into, err := os.MkdirTemp("", "gogit")
	require.NoError(t, err)
	c := MakeTestCloner(t)
	repo, err := c.CloneInto(context.Background(), readOnlyPublicRepo, into)
	require.NoError(t, err)
	require.NotNil(t, repo)
	s, err := os.Stat(filepath.Join(repo.Location(), ".git"))
	require.NoError(t, err)
	require.True(t, s.IsDir())
	require.NoError(t, os.RemoveAll(repo.Location()))
}

package gogit

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func makeTestRepository(t *testing.T) *Repository {
	c := MakeTestCloner(t)
	repo, err := c.Clone(context.Background(), readOnlyPublicRepo)
	require.NoError(t, err)
	return repo
}

func cleanupRepository(t *testing.T, repo *Repository) {
	require.NoError(t, os.RemoveAll(repo.Location()))
}

func TestRepository_Location(t *testing.T) {
	repo := makeTestRepository(t)
	defer cleanupRepository(t, repo)
	require.NotEmpty(t, repo.Location())
	s, err := os.Stat(repo.Location())
	require.NoError(t, err)
	require.True(t, s.IsDir())
}

func TestRepository_AreThereUncommittedChanges(t *testing.T) {
	repo := makeTestRepository(t)
	defer cleanupRepository(t, repo)
	changes, err := repo.AreThereUncommittedChanges(context.Background())
	require.NoError(t, err)
	require.False(t, changes)
	require.NoError(t, os.WriteFile(filepath.Join(repo.Location(), "test.txt"), []byte("test"), 0644))
	changes, err = repo.AreThereUncommittedChanges(context.Background())
	require.NoError(t, err)
	require.True(t, changes)
}

func TestRepository_GuessRemoteName(t *testing.T) {
	repo := makeTestRepository(t)
	defer cleanupRepository(t, repo)
	remoteName, err := repo.GuessRemoteName(context.Background())
	require.NoError(t, err)
	require.Equal(t, "origin", remoteName)
}

func TestRepository_GuessRemoteHead(t *testing.T) {
	repo := makeTestRepository(t)
	defer cleanupRepository(t, repo)
	remoteName, err := repo.GuessRemoteHead(context.Background(), "origin")
	require.NoError(t, err)
	require.Equal(t, "main", remoteName)
}

func TestRepository_CheckoutNewBranch(t *testing.T) {
	repo := makeTestRepository(t)
	defer cleanupRepository(t, repo)
	require.NoError(t, repo.CheckoutNewBranch(context.Background(), "test"))
	currentBranch, err := repo.CurrentBranchName(context.Background())
	require.NoError(t, err)
	require.Equal(t, "test", currentBranch)
}

func TestRepository_CommitAll(t *testing.T) {
	repo := makeTestRepository(t)
	defer cleanupRepository(t, repo)
	require.NoError(t, os.WriteFile(filepath.Join(repo.Location(), "test.txt"), []byte("test"), 0644))
	require.NoError(t, repo.CommitAll(context.Background(), "test"))
	changes, err := repo.AreThereUncommittedChanges(context.Background())
	require.NoError(t, err)
	require.False(t, changes)
}

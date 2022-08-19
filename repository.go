package gogit

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Repository struct {
	logger                Logger
	location              string
	cachedGuessRemoteLock sync.Mutex
	cachedGuessRemote     string
	mu                    sync.RWMutex
}

func (r *Repository) Location() string {
	return r.location
}

func (r *Repository) AreThereUncommittedChanges(ctx context.Context) (bool, error) {
	var res *execResult
	var err error
	res, err = wrappedExec(ctx, r.location, r.logger, "git", "status", "--short")
	if err != nil {
		return false, fmt.Errorf("git status failed: %w", err)
	}
	if res.stdout.Len() > 0 {
		return true, nil
	}
	return false, nil
}

func (r *Repository) GuessRemoteName(ctx context.Context) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	execResult, err := wrappedExec(ctx, r.location, r.logger, "git", "remote", "show")
	if err != nil {
		return "", fmt.Errorf("git remote show failed: %w", err)
	}
	return strings.TrimSpace(execResult.stdout.String()), nil
}

func (r *Repository) GetUserEmail(ctx context.Context) (string, error) {
	execResult, err := wrappedExec(ctx, r.location, r.logger, "git", "config", "--get", "user.email")
	if err != nil {
		return "", fmt.Errorf("git config --get user.email failed: %w", err)
	}
	return strings.TrimSpace(execResult.stdout.String()), nil
}

func (r *Repository) GetUserName(ctx context.Context) (string, error) {
	execResult, err := wrappedExec(ctx, r.location, r.logger, "git", "config", "--get", "user.name")
	if err != nil {
		return "", fmt.Errorf("git config --get user.name failed: %w", err)
	}
	return strings.TrimSpace(execResult.stdout.String()), nil
}

func (r *Repository) SetUserNameAndEmailIfUnset(ctx context.Context, name string, email string) error {
	existingEmail, err := r.GetUserEmail(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user email: %w", err)
	}
	if existingEmail == "" {
		if _, err := wrappedExec(ctx, r.location, r.logger, "git", "config", "user.email", email); err != nil {
			return fmt.Errorf("git config user.email failed: %w", err)
		}
	}
	existingName, err := r.GetUserName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user name: %w", err)
	}
	if existingName == "" {
		if _, err := wrappedExec(ctx, r.location, r.logger, "git", "config", "user.name", name); err != nil {
			return fmt.Errorf("git config user.name failed: %w", err)
		}
	}
	return nil
}

func (r *Repository) GuessRemoteHead(ctx context.Context, remoteName string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	execResult, err := wrappedExec(ctx, r.location, r.logger, "git", "remote", "show", remoteName)
	if err != nil {
		return "", fmt.Errorf("git remote show failed: %w", err)
	}
	lines := strings.Split(execResult.stdout.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "HEAD branch:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "HEAD branch:")), nil
		}
	}
	return "", fmt.Errorf("failed to guess remote head")
}

func (r *Repository) GuessRemote() (string, error) {
	r.cachedGuessRemoteLock.Lock()
	defer r.cachedGuessRemoteLock.Unlock()
	if r.cachedGuessRemote != "" {
		return r.cachedGuessRemote, nil
	}
	remoteName, err := r.GuessRemoteName(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to guess remote name: %w", err)
	}
	remoteHead, err := r.GuessRemoteHead(context.Background(), remoteName)
	if err != nil {
		return "", fmt.Errorf("failed to guess remote head: %w", err)
	}
	r.cachedGuessRemote = fmt.Sprintf("%s/%s", remoteName, remoteHead)
	return r.cachedGuessRemote, nil
}

func (r *Repository) CheckoutNewBranch(ctx context.Context, branch string) error {
	guessedRemote, err := r.GuessRemote()
	if err != nil {
		return fmt.Errorf("failed to guess remote: %w", err)
	}
	if _, err := wrappedExec(ctx, r.location, r.logger, "git", "checkout", "-b", branch, guessedRemote); err != nil {
		return fmt.Errorf("git checkout failed: %w", err)
	}
	return nil
}

func (r *Repository) CommitAll(ctx context.Context, message string) error {
	if _, err := wrappedExec(ctx, r.location, r.logger, "git", "add", "."); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}
	if _, err := wrappedExec(ctx, r.location, r.logger, "git", "commit", "-a", "-m", message); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}
	return nil
}

func (r *Repository) CurrentBranchName(ctx context.Context) (string, error) {
	if r, err := wrappedExec(ctx, r.location, r.logger, "git", "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return "", fmt.Errorf("failed to get current branch name: %w", err)
	} else {
		if r.stdout.String() == "" {
			return "", fmt.Errorf("got an empty branch name")
		}
		return strings.TrimSpace(r.stdout.String()), nil
	}
}

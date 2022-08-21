package gogit

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cresta/pipe"
	"os"
	"strings"
)

type Logger interface {
	Debug(ctx context.Context, msg string, strings map[string]string, ints map[string]int64)
	Info(ctx context.Context, msg string, strings map[string]string, ints map[string]int64)
}

type SilentLogger struct{}

func (n SilentLogger) Debug(_ context.Context, _ string, _ map[string]string, _ map[string]int64) {
}

func (n SilentLogger) Info(_ context.Context, _ string, _ map[string]string, _ map[string]int64) {
}

var _ Logger = &SilentLogger{}

type Cloner struct {
	Logger  Logger
	TempDir string
}

func (p *Cloner) Clone(ctx context.Context, origin string) (*Repository, error) {
	into, err := os.MkdirTemp(p.TempDir, "gogit")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary directory: %w", err)
	}
	return p.CloneInto(ctx, origin, into)
}

func (p *Cloner) CloneInto(ctx context.Context, origin string, into string) (*Repository, error) {
	if _, err := wrappedExec(ctx, "", p.Logger, "git", "clone", origin, into); err != nil {
		return nil, fmt.Errorf("cannot clone repo: %w", err)
	}
	return &Repository{
		location: into,
		logger:   p.Logger,
	}, nil
}

type execResult struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func wrappedExec(ctx context.Context, dir string, logger Logger, cmd string, args ...string) (*execResult, error) {
	var res execResult
	fullCmd := fmt.Sprintf("call %s %s", cmd, strings.Join(args, " "))
	logger.Debug(ctx, fmt.Sprintf("call %s", fullCmd), nil, nil)
	err := pipe.NewPiped(cmd, args...).WithDir(dir).Execute(ctx, nil, &res.stdout, &res.stderr)
	logger.Debug(ctx, fmt.Sprintf("done %s", fullCmd), map[string]string{
		"stdout": res.stdout.String(),
		"stderr": res.stderr.String(),
	}, nil)
	if err != nil {
		return nil, execErr(fmt.Sprintf("exec error %s", fullCmd), res.stdout, res.stderr, err)
	}
	return &res, nil
}

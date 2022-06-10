package gogit

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cresta/pipe"
	"io/ioutil"
)

type Logger interface {
	Debug(ctx context.Context, msg string, strings map[string]string, ints map[string]int64)
	Info(ctx context.Context, msg string, strings map[string]string, ints map[string]int64)
}

type Porcelain struct {
	Logger Logger
}

func (p *Porcelain) Clone(ctx context.Context, origin string) (*Repository, error) {
	into, err := ioutil.TempDir("", "gogit")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary directory: %w", err)
	}
	return p.CloneInto(ctx, origin, into)
}

func (p *Porcelain) CloneInto(ctx context.Context, origin string, into string) (*Repository, error) {
	var stdout, stderr bytes.Buffer
	if err := pipe.NewPiped("git", "clone", origin, into).Execute(ctx, nil, &stdout, &stderr); err != nil {
		return nil, execErr("cannot clone remote", stdout, stderr, err)
	}
	return &Repository{
		location: into,
		logger:   p.Logger,
	}, nil
}

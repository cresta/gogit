package gogit

import "sync"

type Worktree struct {
	mu sync.Mutex
}

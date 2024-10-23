package nef

import (
	"sync"
)

type (
	creator[I any] interface {
		create(root string, overwrite bool, F ChangerFS) I
	}

	creatorFunc[I, F any] func(root string, overwrite bool, fS F) I
)

type lazyObject[I, F any] struct {
	once    sync.Once
	command I
	creator creatorFunc[I, F]
}

func (l *lazyObject[I, F]) instance(root string, overwrite bool, fS F) I {
	l.once.Do(func() {
		l.command = l.creator(root, overwrite, fS)
	})

	return l.command
}

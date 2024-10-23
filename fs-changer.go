package nef

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/snivilised/nefilim/internal/third/lo"
)

const (
	changeOpName = "Change"
)

type (
	changer interface {
		create() changer
		change(from, to string) error
	}

	changeFunc func(from, to string) error

	changers map[bitmask]changeFunc

	baseChanger struct {
		baseOp[ChangerFS]
		actions changers
	}
)

func (m *baseChanger) guard(_, to string) error {
	if strings.Contains(to, "/") {
		return NewInvalidPathError(
			"move rejected, change 'to' path can't contain separator", to,
		)
	}

	return nil
}

func (m *baseChanger) change(from, to string) error {
	if err := m.guard(from, to); err != nil {
		return err
	}

	mask := m.query(from, to)

	if action, exists := m.actions[mask]; exists {
		return action(from, to)
	}

	return NewInvalidBinaryFsOpError(moveOpName, from, to)
}

func (m *baseChanger) query(from, to string) bitmask {
	fromExists, fromIsDir := m.peek(from)
	toExists, toIsDir := m.peek(m.fill(from, to))

	return bitmask{
		fromExists: fromExists,
		toExists:   toExists,
		fromIsDir:  fromIsDir,
		toIsDir:    toIsDir,
	}
}

func (m *baseChanger) fill(from, to string) string {
	// returns the parent from 'from' combined with 'to', ie
	// given: from: 'foo/bar/baz.txt', to: 'pez.txt'
	// returns 'foo/bar/pez.txt'
	//
	return Join(Parent(from), to)
}

func (m *baseChanger) rename(from, to string) error {
	destination := m.fill(from, to)

	if from == destination {
		return nil
	}

	return os.Rename(
		filepath.Join(m.root, from),
		filepath.Join(m.root, destination),
	)
}

type lazyChanger struct {
	once    sync.Once
	changer changer
}

func (l *lazyChanger) instance(root string, overwrite bool, fS ChangerFS) changer {
	l.once.Do(func() {
		l.changer = l.create(root, overwrite, fS)
	})

	return l.changer
}

func (l *lazyChanger) create(root string, overwrite bool, fS ChangerFS) changer {
	// create an interface for this function
	//
	return lo.TernaryF(overwrite,
		func() changer {
			return &overwriteChanger{
				baseChanger: baseChanger{
					baseOp: baseOp[ChangerFS]{
						fS:   fS,
						root: root,
					},
				},
			}
		},
		func() changer {
			return &tentativeChanger{
				baseChanger: baseChanger{
					baseOp: baseOp[ChangerFS]{
						fS:   fS,
						root: root,
					},
				},
			}
		},
	).create()
}

type overwriteChanger struct {
	baseChanger
}

func (m *overwriteChanger) create() changer {
	m.actions = changers{
		{true, false, false, false}: m.rename, // from exists as file, to does not exist
		{true, false, true, false}:  m.rename, // from exists as dir, to does not exist
		{true, true, true, true}:    m.rename, // from exists as dir, to exists as dir
		{true, true, false, false}:  m.rename, // from and to refer to the same existing file
	}

	return m
}

type tentativeChanger struct {
	baseChanger
}

func (m *tentativeChanger) create() changer {
	m.actions = changers{
		{true, false, false, false}: m.rename,              // from exists as file, to does not exist
		{true, false, true, false}:  m.rename,              // from exists as dir, to does not exist
		{true, true, true, true}:    m.rename,              // from exists as dir, to exists as dir
		{true, true, false, false}:  m.rejectFileOverwrite, // from and to may refer to the same existing file
	}

	return m
}

func (m *tentativeChanger) rejectFileOverwrite(from, to string) error {
	// to file already exists
	//
	return NewInvalidBinaryFsOpError("Change", from, to)
}

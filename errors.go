package nef

import (
	"errors"
	"fmt"
)

// IsInvalidPathError
func IsInvalidPathError(err error) bool {
	return errors.Is(errors.Unwrap(err), ErrCoreInvalidPath)
}

// NewInvalidPathError
func NewInvalidPathError(path string) error {
	// should this have a description too to describe the path?
	// we probably need a CoreInvalidPathError which InvalidPathError
	// wraps.
	return fmt.Errorf("path: %q, %w", path, ErrCoreInvalidPath)
}

// IsBinaryFsOpError
func IsBinaryFsOpError(err error) bool {
	return errors.Is(err, ErrCoreBinaryFsOp)
}

// NewInvalidBinaryFsOpError
func NewInvalidBinaryFsOpError(op, from, to string) error {
	return fmt.Errorf("op: %q, from %q, to: %q %w", op, from, to, ErrCoreBinaryFsOp)
}

// IsRejectSameDirMoveError
func IsRejectSameDirMoveError(err error) bool {
	return errors.Is(err, ErrCoreRejectSameDirMove)
}

// NewRejectSameDirMoveError
func NewRejectSameDirMoveError(op, from, to string) error {
	return fmt.Errorf("op: %q, from %q, to: %q %w", op, from, to, ErrCoreRejectSameDirMove)
}

// IsRejectDifferentDirChangeError
func IsRejectDifferentDirChangeError(err error) bool {
	return errors.Is(err, ErrCoreRejectDifferentDirChange)
}

// NewRejectDifferentDirChangeError
func NewRejectDifferentDirChangeError(op, from, to string) error {
	return fmt.Errorf("op: %q, from %q, to: %q %w", op, from, to, ErrCoreRejectDifferentDirChange)
}

// these errors are deliberately being exported, so that client libraries (eg traverse) that do support i18n
// can wrap them and make them translate-able.
var (
	ErrCoreInvalidPath              = errors.New("invalid path")
	ErrCoreBinaryFsOp               = errors.New("invalid binary file system operation")
	ErrCoreRejectSameDirMove        = errors.New("same directory move rejected, use move instead")
	ErrCoreRejectDifferentDirChange = errors.New("different directory change rejected, use move instead")
)

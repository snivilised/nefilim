package nef

import (
	"errors"
	"fmt"
)

// IsInvalidPathError reports whether err is or wraps the invalid path error.
func IsInvalidPathError(err error) bool {
	return errors.Is(errors.Unwrap(err), ErrCoreInvalidPath)
}

// NewInvalidPathError returns an error indicating an invalid path,
// with the given description and path.
func NewInvalidPathError(description, path string) error {
	// should this have a description too to describe the path?
	// we probably need a CoreInvalidPathError which InvalidPathError
	// wraps.
	return fmt.Errorf("description: %q, path: %q, %w",
		description, path, ErrCoreInvalidPath,
	)
}

// IsBinaryFsOpError reports whether err is or wraps the invalid binary
// file system operation error.
func IsBinaryFsOpError(err error) bool {
	return errors.Is(err, ErrCoreBinaryFsOp)
}

// NewInvalidBinaryFsOpError returns an error for an invalid binary FS
// operation (op, from, to).
func NewInvalidBinaryFsOpError(op, from, to string) error {
	return fmt.Errorf("op: %q, from %q, to: %q %w", op, from, to, ErrCoreBinaryFsOp)
}

// IsRejectSameDirMoveError reports whether err is or wraps the same-directory
// move rejection error.
func IsRejectSameDirMoveError(err error) bool {
	return errors.Is(err, ErrCoreRejectSameDirMove)
}

// NewRejectSameDirMoveError returns an error when a move within the same
// directory is rejected.
func NewRejectSameDirMoveError(op, from, to string) error {
	return fmt.Errorf("op: %q, from %q, to: %q %w", op, from, to, ErrCoreRejectSameDirMove)
}

// IsRejectDifferentDirChangeError reports whether err is or wraps the different-directory
// change rejection error.
func IsRejectDifferentDirChangeError(err error) bool {
	return errors.Is(err, ErrCoreRejectDifferentDirChange)
}

// NewRejectDifferentDirChangeError returns an error when a change across
// different directories is rejected.
func NewRejectDifferentDirChangeError(op, from, to string) error {
	return fmt.Errorf("op: %q, from %q, to: %q %w", op, from, to, ErrCoreRejectDifferentDirChange)
}

// these errors are deliberately being exported, so that client libraries
// (eg traverse) that do support i18n
// can wrap them and make them translate-able.
var (
	// ErrCoreInvalidPath indicates an invalid path	
	ErrCoreInvalidPath              = errors.New("invalid path")
	// ErrCoreBinaryFsOp indicates an invalid binary file system operation
	ErrCoreBinaryFsOp               = errors.New("invalid binary file system operation")
	// ErrCoreRejectSameDirMove indicates a same directory move is rejected
	ErrCoreRejectSameDirMove        = errors.New("same directory move rejected, use move instead")
	// ErrCoreRejectDifferentDirChange indicates a different directory change is rejected
	ErrCoreRejectDifferentDirChange = errors.New("different directory change rejected, use move instead")
)

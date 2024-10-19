package nef

import (
	"errors"
	"fmt"
)

// IsInvalidPathError
func IsInvalidPathError(err error) bool {
	return errors.Is(err, ErrCoreInvalidPath)
}

// NewInvalidPathError
func NewInvalidPathError(path string) error {
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

// these errors are deliberately being exported, so that client libraries (eg traverse) that do support i18n
// can wrap them and make them translate-able.
var (
	ErrCoreInvalidPath       = errors.New("invalid path")
	ErrCoreBinaryFsOp        = errors.New("invalid binary file system operation")
	ErrCoreRejectSameDirMove = errors.New("same directory move rejected, use move instead")
)

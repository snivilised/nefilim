package nef

import (
	"errors"
	"fmt"
)

// IsInvalidPathError
func IsInvalidPathError(err error) bool {
	return errors.Is(err, errCoreInvalidPath)
}

// NewInvalidPathError
func NewInvalidPathError(path string) error {
	return fmt.Errorf("path: %q, %w", path, errCoreInvalidPath)
}

// IsBinaryFsOpError
func IsBinaryFsOpError(err error) bool {
	return errors.Is(err, errCoreBinaryFsOp)
}

// NewInvalidBinaryFsOpError
func NewInvalidBinaryFsOpError(op, from, to string) error {
	return fmt.Errorf("op: %q, from %q, to: %q %w", op, from, to, errCoreBinaryFsOp)
}

// IsRejectSameDirMoveError
func IsRejectSameDirMoveError(err error) bool {
	return errors.Is(err, errCoreRejectSameDirMove)
}

// NewRejectSameDirMoveError
func NewRejectSameDirMoveError(op, from, to string) error {
	return fmt.Errorf("op: %q, from %q, to: %q %w", op, from, to, errCoreRejectSameDirMove)
}

var (
	errCoreInvalidPath       = errors.New("invalid path")
	errCoreBinaryFsOp        = errors.New("invalid binary file system operation")
	errCoreRejectSameDirMove = errors.New("same directory move rejected, use move instead")
)

package collections_test

import (
	"fmt"

	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/types"
	"github.com/snivilised/nefilim/collections"
)

// HaveSizeMatcher asserts that a stack has the expected size.
type HaveSizeMatcher struct {
	size uint
}

// HaveSize returns a Gomega matcher that expects a *collections.Stack to have the given size.
func HaveSize(size uint) types.GomegaMatcher {
	return &HaveSizeMatcher{
		size: size,
	}
}

// Match runs the matcher; actual must be *collections.Stack[string].
func (m *HaveSizeMatcher) Match(actual interface{}) (bool, error) {
	stack, ok := actual.(*collections.Stack[string])

	if !ok {
		return false, fmt.Errorf("matcher expected a *collections.Stack[T] (%T)", stack)
	}

	return stack.Size() == m.size, nil
}

// FailureMessage returns the message when the stack size does not match.
func (m *HaveSizeMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack to have size: %v\n", m.size)
}

// NegatedFailureMessage returns the message when the negated matcher fails.
func (m *HaveSizeMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack NOT to have size: %v\n", m.size)
}

// HaveCurrentMatcher asserts that a stack's top (current) element equals the expected value.
type HaveCurrentMatcher struct {
	current string
}

// HaveCurrent returns a Gomega matcher that expects the stack's current (top) value to equal current.
func HaveCurrent(current string) types.GomegaMatcher {
	return &HaveCurrentMatcher{
		current: current,
	}
}

// Match runs the matcher; actual must be *collections.Stack[string].
func (m *HaveCurrentMatcher) Match(actual interface{}) (bool, error) {
	stack, ok := actual.(*collections.Stack[string])

	if !ok {
		return false, fmt.Errorf("matcher expected a *collections.Stack[T] (%T)", stack)
	}

	current, _ := stack.Current()

	return current == m.current, nil
}

// FailureMessage returns the message when the stack's current value does not match.
func (m *HaveCurrentMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack to have current value of: %v\n", m.current)
}

// NegatedFailureMessage returns the message when the negated matcher fails.
func (m *HaveCurrentMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack NOT to have current value of: %v\n", m.current)
}

// BeInCorrectState returns a matcher that expects the stack to have the given size and current value.
func BeInCorrectState(size uint, current string) types.GomegaMatcher {
	return And(
		HaveSize(size),
		HaveCurrent(current),
	)
}

// HavePoppedMatcher asserts that a stack had the expected size and popped value after a pop.
type HavePoppedMatcher struct {
	size       uint
	actualItem string
}

// WithExpectedPop holds the stack and the value that was popped, for use with HavePopped.
type WithExpectedPop struct {
	stack  *collections.Stack[string]
	popped string
}

// HavePopped returns a Gomega matcher that expects actual (*WithExpectedPop) to have the given size and popped item.
func HavePopped(size uint, actual string) types.GomegaMatcher {
	return &HavePoppedMatcher{
		size:       size,
		actualItem: actual,
	}
}

// Match runs the matcher; actual must be *WithExpectedPop.
func (m *HavePoppedMatcher) Match(actual interface{}) (bool, error) {
	expectation, ok := actual.(*WithExpectedPop)

	if !ok {
		return false, fmt.Errorf("matcher expected a *ExpectedPop (%T)", expectation)
	}

	result := expectation.stack.Size() == m.size && m.actualItem == expectation.popped

	return result, nil
}

// FailureMessage returns the message when the stack size or popped item does not match.
func (m *HavePoppedMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack to\n\thave size: %v\n\tand popped item: %v\n",
		m.size, m.actualItem,
	)
}

// NegatedFailureMessage returns the message when the negated matcher fails.
func (m *HavePoppedMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack NOT to\n\thave size: %v\n\tand popped item: %v\n",
		m.size, m.actualItem,
	)
}

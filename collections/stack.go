package collections

// NewStack returns a new empty stack of type T.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		content: []T{},
	}
}

// NewStackWith returns a new stack pre-filled with the given items (top is last element of with).
func NewStackWith[T any](with []T) *Stack[T] {
	return &Stack[T]{
		content: with,
	}
}

// Stack is a generic LIFO stack; the top element is the last one pushed.
type Stack[T any] struct {
	content []T
}

// Push adds item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.content = append(s.content, item)
}

// Pop removes and returns the top element, or ErrStackIsEmpty if the stack is empty.
func (s *Stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, ErrStackIsEmpty
	}

	item := s.pop()

	return item, nil
}

// MustPop removes and returns the top element; it panics with ErrStackIsEmpty if the stack is empty.
func (s *Stack[T]) MustPop() T {
	if s.IsEmpty() {
		panic(ErrStackIsEmpty)
	}

	return s.pop()
}

// Current returns the top element without removing it, or ErrStackIsEmpty if the stack is empty.
func (s *Stack[T]) Current() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, ErrStackIsEmpty
	}

	return s.content[s.top()], nil
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() uint {
	return uint(len(s.content))
}

// IsEmpty reports whether the stack has no elements.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.content) == 0
}

// Content returns a copy of the stack's elements from bottom to top.
func (s *Stack[T]) Content() []T {
	return s.content
}

func (s *Stack[T]) top() int {
	return len(s.content) - 1
}

func (s *Stack[T]) pop() T {
	t := s.top()
	item := s.content[t]
	s.content = s.content[:t]

	return item
}

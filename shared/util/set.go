package util

type Set[T comparable] map[T]bool

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(value T) {
	s[value] = true
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s Set[T]) ToSlice() []T {
	keys := make([]T, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}

func (s Set[T]) Union(s2 Set[T]) Set[T] {
	r := NewSet[T]()

	helper := func(st Set[T]) {
		for k, _ := range st {
			r.Add(k)
		}
	}

	helper(s)
	helper(s2)

	return r
}

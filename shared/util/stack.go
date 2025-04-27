package util

import (
	"fmt"
	"reflect"
)

type Stack[T any] []T

func (s *Stack[T]) Push(value T) {
	*s = append(*s, value)
}

func (s *Stack[T]) Top() (T, error) {
	var zero T // use go's default value
	if len(*s) == 0 {
		return zero, fmt.Errorf("stack is empty")
	}
	return (*s)[len(*s)-1], nil
}

func (s *Stack[T]) Pop() (T, error) {
	var zero T // use go's default value
	if len(*s) == 0 {
		return zero, fmt.Errorf("stack is empty")
	}
	top := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return top, nil
}

func (s *Stack[T]) Contains(value T) bool {
	for _, el := range *s {
		if reflect.DeepEqual(value, el) {
			return true
		}
	}
	return false
}

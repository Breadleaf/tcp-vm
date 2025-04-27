package util

import (
	"fmt"
)

type Queue[T comparable] []T

// enqueue
func (q *Queue[T]) Push(value T) {
	*q = append(*q, value)
}

// dequeue
func (q *Queue[T]) Pop() (T, error) {
	var zero T // use go's default value
	if len(*q) == 0 {
		return zero, fmt.Errorf("queue is empty")
	}
	front := (*q)[0]
	*q = (*q)[1:]
	return front, nil
}

func (q Queue[T]) Peek() (T, error) {
	var zero T // use go's default value
	if len(q) == 0 {
		return zero, fmt.Errorf("queue is empty")
	}
	return q[0], nil
}

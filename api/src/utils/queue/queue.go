package queue

import (
	q "github.com/emirpasic/gods/queues/arrayqueue"
)

// ------------------------------------------------------------
// : Queue
// ------------------------------------------------------------
type Iterator = q.Iterator

type Queue struct {
	data *q.Queue
}
// ------------------------------------------------------------
// : Constructor
// ------------------------------------------------------------
func New() *Queue {
	return &Queue{data: q.New()}
}

// ------------------------------------------------------------
// : Statics
// ------------------------------------------------------------
func FromJSON(data []byte) (*Queue, error) {
	q   := New()
	err := q.data.FromJSON(data)

	return q, err
}

// ------------------------------------------------------------
// : Methods
// ------------------------------------------------------------
func (q *Queue) Enqueue(v any) {
	q.data.Enqueue(v)
}

func (q *Queue) Dequeue() (any, bool) {
	return q.data.Dequeue()
}

func (q *Queue) Peek() (any, bool) {
	return q.data.Peek()
}

func (q *Queue) Empty() bool {
	return q.data.Empty()
}

func (q *Queue) Size() int {
	return q.data.Size()
}

func (q *Queue) Clear() {
	q.data.Clear()
}

func (q *Queue) Iterator() Iterator {
	return q.data.Iterator()
}

func (q *Queue) Values() []any {
	return q.data.Values()
}

func (q *Queue) String() string {
	return q.data.String()
}

func (q *Queue) ToJSON() ([]byte, error) {
	return q.data.ToJSON()
}

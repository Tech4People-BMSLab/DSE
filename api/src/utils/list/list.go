package list

import (
	c "github.com/emirpasic/gods/lists/arraylist"
)

// ------------------------------------------------------------
// : List
// ------------------------------------------------------------
type List struct {
	data *c.List
}
// ------------------------------------------------------------
// : Constructor
// ------------------------------------------------------------
func New() *List {
	return &List{
		data: c.New(),
	}
}
// ------------------------------------------------------------
// : Statics
// ------------------------------------------------------------
func FromJSON(data []byte) (*List, error) {
	l  := New()
	err := l.FromJSON(data)
	return l, err
}

func FromSlice(data []any) *List {
	l := New()

	for _, v := range data {
		l.Add(v)
	}

	return l
}
// ------------------------------------------------------------
// : Methods
// ------------------------------------------------------------
func (l *List) FromJSON(data []byte) error {
	return l.data.FromJSON(data)
}

func (l *List) ToJSON() ([]byte, error) {
	return l.data.ToJSON()
}

func (l *List) Contains(v any) bool {
	return l.data.Contains(v)
}

func (l *List) Add(v any) {
	l.data.Add(v)
}

func (l *List) Get(i int) (any, bool) {
	return l.data.Get(i)
}

func (l *List) Set(i int, v any) {
	l.data.Set(i, v)
}

func (l *List) Remove(i int) {
	l.data.Remove(i)
}

func (l *List) IndexOf(v any) int {
	return l.data.IndexOf(v)
}

func (l *List) Iterator() c.Iterator {
    return l.data.Iterator()
}

func (l *List) Size() int {
	return l.data.Size()
}

func (l *List) IsEmpty() bool {
	return l.data.Empty()
}

func (l *List) Find(predicate func(int, any) bool) (int, any) {
	return l.data.Find(predicate)
}

func (l *List) Filter(predicate func(int, any) bool) *List {
	f    := l.data.Select(predicate)
	iter := f.Iterator()

	fl := New()

	for iter.Next() {
		fl.Add(iter.Value())
	}

	return fl
}

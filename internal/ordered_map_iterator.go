package internal

import (
    "fmt"
)

type OrderedMapIterator[V any] struct {
    index int
    data  *OrderedMap[V]
}

func (m *OrderedMap[V]) CreateIterator() OrderedMapIterator[V] {
    return OrderedMapIterator[V]{
        index: -1,
        data:  m,
    }
}

func (i *OrderedMapIterator[V]) SkipUntil(searchKey string) bool {
    for i.Next() {
        if key, exists := i.Key(); exists && key == searchKey {
            return true
        }
    }

    return false
}

func (i *OrderedMapIterator[V]) HasNext() bool {
    return i.index <= len(i.data.keys)
}

func (i *OrderedMapIterator[V]) Next() bool {
    hasNext := i.HasNext()
    if hasNext {
        i.index++
    }

    return hasNext
}

func (i *OrderedMapIterator[V]) HasPrevious() bool {
    return i.index > 0
}

func (i *OrderedMapIterator[V]) Previous() bool {
    hasPrev := i.HasPrevious()
    if hasPrev {
        i.index--
    }

    return hasPrev
}

func (i *OrderedMapIterator[V]) Skip(count int) bool {
    for idx := 0; idx < count; idx++ {
        if !i.Next() {
            return false
        }
    }

    return true
}

func (i *OrderedMapIterator[V]) Rewind(count int) bool {
    if i.index-count >= -1 {
        i.index -= count
        return true
    } else {
        i.index = -1
        return false
    }
}

func (i *OrderedMapIterator[V]) Take(count int) []V {
    var data []V
    for idx := 0; idx < count; idx++ {
        if i.Next() {
            if val, exists := i.Value(); exists {
                data = append(data, val)
            }
        }
    }

    return data
}

func (i *OrderedMapIterator[V]) Key() (string, bool) {
    if i.index > len(i.data.keys)-1 || i.index < 0 {
        return "", false
    }
    return i.data.keys[i.index], true
}

func (i *OrderedMapIterator[V]) Value() (V, bool) {
    var res V
    if key, exists := i.Key(); !exists {
        return res, false
    } else {
        return i.data.Get(key)
    }
}

func (i *OrderedMapIterator[V]) Pair() (string, V, bool) {
    var res V
    if key, exists := i.Key(); !exists {
        return "", res, false
    } else {
        value, exists := i.data.Get(key)
        return key, value, exists
    }
}

func (i *OrderedMapIterator[V]) Reset() {
    i.index = 0
}

func (i *OrderedMapIterator[V]) End() {
    i.index = len(i.data.keys) - 1
}

func (i *OrderedMapIterator[V]) Set(index int) error {
    if index < 0 || index > len(i.data.keys)-1 {
        return fmt.Errorf("index out of range")
    }

    i.index = index

    return nil
}

package internal

import (
    "strings"
    "sync"
)

type OrderedMap[V any] struct {
    inner map[string]V
    keys  []string

    rw *sync.RWMutex
}

func CreateOrderedMap[V any]() OrderedMap[V] {
    return OrderedMap[V]{
        inner: map[string]V{},
        keys:  []string{},
        rw:    &sync.RWMutex{},
    }
}

func (m *OrderedMap[V]) Get(key string) (V, bool) {
    var val V

    m.rw.RLock()
    defer m.rw.RUnlock()

    val, exists := m.inner[key]
    return val, exists
}

func (m *OrderedMap[V]) Set(key string, value V) {
    m.rw.Lock()
    defer m.rw.Unlock()

    if _, exists := m.inner[key]; !exists {
        //m.keys = append(m.keys, key)
        //sort.Strings(m.keys)
        if len(m.keys) > 0 {
            for i, target := range m.keys {
                if strings.Compare(target, key) == 1 {
                    m.keys = append(m.keys[:i+1], m.keys[i:]...)
                    m.keys[i] = key
                    break
                }

                if len(m.keys)-1 == i {
                    m.keys = append(m.keys, key)
                }
            }
        } else {
            m.keys = append(m.keys, key)
        }
    }

    m.inner[key] = value
}

func (m *OrderedMap[V]) Length() int {
    return len(m.keys)
}

func (m *OrderedMap[V]) First() V {
    var val V
    if len(m.keys) == 0 {
        return val
    }

    val, _ = m.Get(m.keys[0])
    return val
}

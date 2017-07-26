package atomicutil

import (
	"fmt"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	m := NewMapStringUint64()
	m.Set("foo", 33)
	v, ok := m.Get("foo")
	if v != 33 {
		t.Errorf("expected 33, got %v", v)
	}
	if !ok {
		t.Errorf("expected ok, got !ok")
	}
	m.Set("bar", 0)
	v, ok = m.Get("bar")
	if v != 0 {
		t.Errorf("expected 0, got %v", v)
	}
	if !ok {
		t.Errorf("expected ok, got !ok")
	}
	m.Reset()
	v, ok = m.Get("foo")
	if v != 0 {
		t.Errorf("expected 0, got %v", v)
	}
	if ok {
		t.Errorf("expected !ok, got ok")
	}
	v, ok = m.Get("bar")
	if v != 0 {
		t.Errorf("expected 0, got %v", v)
	}
	if ok {
		t.Errorf("expected !ok, got ok")
	}
}

func TestInc(t *testing.T) {
	m := NewMapStringUint64()
	v, ok := m.Get("foo")
	if ok {
		t.Errorf("expected !ok, got ok")
	}
	if v != 0 {
		t.Errorf("expected 0, got %v", v)
	}
	m.Inc("foo")
	v, ok = m.Get("foo")
	if !ok {
		t.Errorf("expected ok, got !ok")
	}
	if v != 1 {
		t.Errorf("expected 1, got %v", v)
	}
	m.IncN("foo", 500)
	v, ok = m.Get("foo")
	if !ok {
		t.Errorf("expected ok, got !ok")
	}
	if v != 501 {
		t.Errorf("expected 501, got %v", v)
	}
}

func BenchmarkConcurrentReadWrite(b *testing.B) {
	m := NewMapStringUint64()
	var names = []string{}
	b.StopTimer()
	for i := 0; i < 100000; i++ {
		names = append(names, fmt.Sprintf("%v", i))
	}
	b.StartTimer()
	b.ResetTimer()
	go func() {
		i := 0
		for {
			i++
			m.Get(names[i%len(names)])
		}
	}()
	for i := 0; i < b.N; i++ {
		m.Inc(names[i%len(names)])
	}
	if m.Len() == 0 {
		b.Errorf("expected > 0, got 0")
	}
}

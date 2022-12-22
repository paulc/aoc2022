package set

import (
	"testing"
)

func TestSetBasic(t *testing.T) {
	s := NewSetFrom([]string{"AA", "BB", "CC", "AA"})
	if s.Len() != 3 {
		t.Error(s)
	}
	if !s.Has("AA") {
		t.Error(s)
	}
	if s.Has("ZZ") {
		t.Error(s)
	}
	s.Add("XX")
	s.Add("XX")
	if !s.Equals(NewSetFrom([]string{"AA", "BB", "CC", "XX"})) {
		t.Error(s)
	}
	s.Remove("XX")
	s.Remove("XX")
	if !s.Equals(NewSetFrom([]string{"AA", "BB", "CC"})) {
		t.Error(s)
	}
	l := s.Len()
	for i := 0; i < l; i++ {
		v, ok := s.Pop()
		if !ok {
			t.Error("Pop", s, v)
		}
	}
	v, ok := s.Pop()
	if ok {
		t.Error("Pop", s, v)
	}
}

func TestSetIterator(t *testing.T) {
	s1 := NewSetFrom([]string{"AA", "BB", "CC", "AA"})
	s2 := s1.Copy()
	for v := range s1.Iter() {
		s2.Remove(v)
	}
	if s2.Len() != 0 {
		t.Error(s2)
	}
}

func TestSetIntersection(t *testing.T) {
	s1 := NewSetFrom([]string{"AA", "BB", "CC"})
	s2 := NewSetFrom([]string{"BB", "CC", "DD"})
	s3 := s1.Intersection(s2)
	if !s3.Equals(NewSetFrom([]string{"BB", "CC"})) {
		t.Error(s3)
	}
}

func TestSetUnion(t *testing.T) {
	s1 := NewSetFrom([]string{"AA", "BB", "CC"})
	s2 := NewSetFrom([]string{"BB", "CC", "DD"})
	s3 := s1.Union(s2)
	if !s3.Equals(NewSetFrom([]string{"AA", "BB", "CC", "DD"})) {
		t.Error(s3)
	}
}

func TestSetDifference(t *testing.T) {
	s1 := NewSetFrom([]string{"AA", "BB", "CC"})
	s2 := NewSetFrom([]string{"BB", "CC", "DD"})
	s3 := s1.Difference(s2)
	if !s3.Equals(NewSetFrom([]string{"AA"})) {
		t.Error(s3)
	}
}

package venn

import (
	"iter"
	"slices"
	"testing"
)

// onceSeq returns a single-use iter.Seq that fails the test
// if consumed more than once. It verifies that Seq methods
// do at most one pass over their argument.
func onceSeq[E any](t *testing.T, elems ...E) iter.Seq[E] {
	t.Helper()
	consumed := false
	return func(yield func(E) bool) {
		if consumed {
			t.Fatal("sequence consumed more than once")
		}
		consumed = true
		for _, e := range elems {
			if !yield(e) {
				return
			}
		}
	}
}

func TestMakeLenContains(t *testing.T) {
	s := Make(1, 2, 3, 2)
	if s.Len() != 3 {
		t.Errorf("Len() = %d; want 3", s.Len())
	}
	if !s.ContainsAll(1, 2, 3) {
		t.Errorf("ContainsAll(1, 2, 3) = false; want true")
	}
	if s.Contains(4) {
		t.Errorf("Contains(4) = true; want false")
	}
}

func TestAddRemove(t *testing.T) {
	s := Make[string]()
	s.Add("a")
	s.AddAll("b", "c", "b")
	if !s.Equal(Make("a", "b", "c")) {
		t.Errorf("after adds: %v; want {a b c}", s)
	}
	s.Remove("a")
	s.Remove("nope") // absent: no-op
	s.RemoveAll("b", "nope")
	if !s.Equal(Make("c")) {
		t.Errorf("after removes: %v; want {c}", s)
	}
}

func TestPop(t *testing.T) {
	s := Make(42)
	e, ok := s.Pop()
	if !ok || e != 42 {
		t.Errorf("Pop() = %v, %v; want 42, true", e, ok)
	}
	if _, ok := s.Pop(); ok {
		t.Error("Pop() on empty set: ok = true; want false")
	}
}

func TestClearCopy(t *testing.T) {
	s := Make(1, 2)
	c := s.Copy()
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("after Clear: Len() = %d; want 0", s.Len())
	}
	if !c.Equal(Make(1, 2)) {
		t.Errorf("copy affected by Clear on original: %v", c)
	}
}

func TestToSliceAll(t *testing.T) {
	s := Make(3, 1, 2)
	got := s.ToSlice()
	slices.Sort(got)
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("ToSlice() sorted = %v; want [1 2 3]", got)
	}
	got = slices.Sorted(s.All())
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("Sorted(All()) = %v; want [1 2 3]", got)
	}
}

func TestPredicates(t *testing.T) {
	s := Make(1, 2)
	z := Make(1, 2, 3)
	if !s.SubsetOf(z) || s.SupersetOf(z) {
		t.Errorf("SubsetOf = %v, SupersetOf = %v; want true, false",
			s.SubsetOf(z), s.SupersetOf(z))
	}
	if !z.SupersetOf(s) || z.SubsetOf(s) {
		t.Errorf("SupersetOf = %v, SubsetOf = %v; want true, false",
			z.SupersetOf(s), z.SubsetOf(s))
	}
	if !s.SubsetOf(s.Copy()) || !s.SupersetOf(s.Copy()) {
		t.Error("set must be subset and superset of an equal set")
	}
	if s.Equal(z) || !z.Equal(Make(3, 2, 1)) {
		t.Errorf("Equal: got %v and %v; want false and true",
			s.Equal(z), z.Equal(Make(3, 2, 1)))
	}
}

var binaryOpCases = []struct {
	name    string
	s, z    []int
	op      func(s, z Set[int]) Set[int]
	want    []int
}{
	{"Union", []int{1, 2}, []int{2, 3}, Set[int].Union, []int{1, 2, 3}},
	{"Intersection", []int{1, 2, 3}, []int{2, 3, 4}, Set[int].Intersection, []int{2, 3}},
	{"Difference", []int{1, 2, 3}, []int{2, 4}, Set[int].Difference, []int{1, 3}},
	{"SymmetricDifference", []int{1, 2, 3}, []int{2, 3, 4}, Set[int].SymmetricDifference, []int{1, 4}},
}

func TestBinaryOps(t *testing.T) {
	for _, c := range binaryOpCases {
		t.Run(c.name, func(t *testing.T) {
			s, z := Make(c.s...), Make(c.z...)
			got := c.op(s, z)
			if !got.Equal(Make(c.want...)) {
				t.Errorf("%s: got %v; want %v", c.name, got, Make(c.want...))
			}
			if !s.Equal(Make(c.s...)) || !z.Equal(Make(c.z...)) {
				t.Errorf("%s modified its operands: s=%v z=%v", c.name, s, z)
			}
		})
	}
}

func TestUpdateOps(t *testing.T) {
	cases := []struct {
		name string
		s, z []int
		op   func(s, z Set[int])
		want []int
	}{
		{"UnionUpdate", []int{1, 2}, []int{2, 3}, Set[int].UnionUpdate, []int{1, 2, 3}},
		{"IntersectionUpdate", []int{1, 2, 3}, []int{2, 3, 4}, Set[int].IntersectionUpdate, []int{2, 3}},
		{"DifferenceUpdate", []int{1, 2, 3}, []int{2, 4}, Set[int].DifferenceUpdate, []int{1, 3}},
		{"SymmetricDifferenceUpdate", []int{1, 2, 3}, []int{2, 3, 4}, Set[int].SymmetricDifferenceUpdate, []int{1, 4}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Make(c.s...)
			c.op(s, Make(c.z...))
			if !s.Equal(Make(c.want...)) {
				t.Errorf("%s: got %v; want %v", c.name, s, Make(c.want...))
			}
		})
	}
}

func TestCollect(t *testing.T) {
	s := Collect(onceSeq(t, 1, 2, 2, 3, 1))
	if !s.Equal(Make(1, 2, 3)) {
		t.Errorf("Collect: got %v; want {1 2 3}", s)
	}
	if got := Collect(onceSeq[int](t)); got.Len() != 0 {
		t.Errorf("Collect of empty seq: got %v; want empty set", got)
	}
}

func TestAddSeqRemoveSeq(t *testing.T) {
	s := Make(1, 2)
	s.AddSeq(onceSeq(t, 2, 3, 3, 4))
	if !s.Equal(Make(1, 2, 3, 4)) {
		t.Errorf("AddSeq: got %v; want {1 2 3 4}", s)
	}
	s.RemoveSeq(onceSeq(t, 3, 4, 4, 99))
	if !s.Equal(Make(1, 2)) {
		t.Errorf("RemoveSeq: got %v; want {1 2}", s)
	}
}

// Seq variants must give the same results as their set-taking
// counterparts, even when the sequence yields duplicates,
// and must not modify the receiver.
func TestSeqOpsMatchSetOps(t *testing.T) {
	cases := []struct {
		name string
		op   func(s Set[int], seq iter.Seq[int]) Set[int]
	}{
		{"UnionSeq", Set[int].UnionSeq},
		{"IntersectionSeq", Set[int].IntersectionSeq},
		{"DifferenceSeq", Set[int].DifferenceSeq},
		{"SymmetricDifferenceSeq", Set[int].SymmetricDifferenceSeq},
	}
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Make(1, 2, 3)
			// seq yields the elements of {2 3 4}, with duplicates
			seq := onceSeq(t, 2, 3, 2, 4, 3)
			want := binaryOpCases[i].op(s.Copy(), Make(2, 3, 4))
			got := c.op(s, seq)
			if !got.Equal(want) {
				t.Errorf("%s: got %v; want %v", c.name, got, want)
			}
			if !s.Equal(Make(1, 2, 3)) {
				t.Errorf("%s modified the receiver: %v", c.name, s)
			}
		})
	}
}

func TestSeqOpsOnEmptySeq(t *testing.T) {
	s := Make(1, 2)
	if got := s.UnionSeq(onceSeq[int](t)); !got.Equal(s) {
		t.Errorf("UnionSeq(empty): got %v; want %v", got, s)
	}
	if got := s.IntersectionSeq(onceSeq[int](t)); got.Len() != 0 {
		t.Errorf("IntersectionSeq(empty): got %v; want empty set", got)
	}
	if got := s.DifferenceSeq(onceSeq[int](t)); !got.Equal(s) {
		t.Errorf("DifferenceSeq(empty): got %v; want %v", got, s)
	}
	if got := s.SymmetricDifferenceSeq(onceSeq[int](t)); !got.Equal(s) {
		t.Errorf("SymmetricDifferenceSeq(empty): got %v; want %v", got, s)
	}
}

func TestString(t *testing.T) {
	if got := Make[int]().String(); got != "{}" {
		t.Errorf("empty set String() = %q; want %q", got, "{}")
	}
	if got := Make(7).String(); got != "{7}" {
		t.Errorf("String() = %q; want %q", got, "{7}")
	}
}

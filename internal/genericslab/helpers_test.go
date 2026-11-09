package genericslab

import (
	"reflect"
	"testing"
)

func TestUniqueUserIDs(t *testing.T) {
	got := Unique([]int64{7, 7, 8, 7, 9})
	want := []int64{7, 8, 9}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestFirst(t *testing.T) {
	v, ok := First([]string{"a", "b"})
	if !ok || v != "a" {
		t.Fatalf("got=%q ok=%v", v, ok)
	}
	_, ok = First([]int{})
	if ok {
		t.Fatal("empty should be false")
	}
}

func TestMapIDs(t *testing.T) {
	type en struct{ UserID int64 }
	ids := MapIDs([]en{{1}, {2}}, func(e en) int64 { return e.UserID })
	if !reflect.DeepEqual(ids, []int64{1, 2}) {
		t.Fatalf("ids=%v", ids)
	}
}

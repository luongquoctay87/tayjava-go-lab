package course

import "testing"

func BenchmarkDescribeSlow(b *testing.B) {
	c := New(1, "Golang from zero", 30)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = DescribeSlow(c, 10)
	}
}

func BenchmarkDescribeFast(b *testing.B) {
	c := New(1, "Golang from zero", 30)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = DescribeFast(c, 10)
	}
}

func TestDescribeSameShape(t *testing.T) {
	c := New(1, "Go", 30)
	if DescribeSlow(c, 5) != DescribeFast(c, 5) {
		t.Fatal("slow/fast labels differ")
	}
}

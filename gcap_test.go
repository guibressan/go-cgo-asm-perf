package gcap

import (
	"gcap/addasm"
	"gcap/addc"
	"gcap/addgo"
	"testing"
)

func BenchmarkWarmup(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var v [1 << 10]byte
			for i := range v {
				v[i] = 0xFF
			}
		}
	})
}

func BenchmarkAsmAdd(b *testing.B) {
	for range b.N {
		r := addasm.AsmAdd(2, 3)
		if r != 5 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkCAdd(b *testing.B) {
	for range b.N {
		r := addc.CAdd(2, 3)
		if r != 5 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkGoAdd(b *testing.B) {
	for range b.N {
		r := addgo.GoAdd(2, 3)
		if r != 5 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkAsmAddParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := addasm.AsmAdd(2, 3)
			if r != 5 {
				b.Fatal("unexpected")
			}
		}
	})
}

func BenchmarkCAddParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := addc.CAdd(2, 3)
			if r != 5 {
				b.Fatal("unexpected")
			}
		}
	})
}

func BenchmarkGoAddParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := addgo.GoAdd(2, 3)
			if r != 5 {
				b.Fatal("unexpected")
			}
		}
	})
}

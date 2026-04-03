package addc

/*
long c_add(long a, long b) {
	return a + b;
}
*/
import "C"

//go:noinline
func CAdd(a, b int) int {
	return int(C.c_add(C.long(a), C.long(b)))
}

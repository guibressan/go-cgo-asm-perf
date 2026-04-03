package addasm

import _ "unsafe"

//go:noescape
func AsmAdd(a, b int) int

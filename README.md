# Go/CGO/ASM performance comparison

Compare call overhead between Go,CGO and ASM on amd64

```
$ go version
go version go1.26.1 linux/amd64
```

AsmAdd assembly:
```
TEXT gcap/addasm.AsmAdd.abi0(SB) addasm/add_amd64.s
  add_amd64.s:4		0x12a79c0		488b442408		MOVQ 0x8(SP), AX	
  add_amd64.s:5		0x12a79c5		488b5c2410		MOVQ 0x10(SP), BX	
  add_amd64.s:6		0x12a79ca		4801d8			ADDQ BX, AX		
  add_amd64.s:7		0x12a79cd		4889442418		MOVQ AX, 0x18(SP)	
  add_amd64.s:8		0x12a79d2		c3			RET			
```

Go assumes ABI0 for user written assembly, so arguments and returns must be in
the stack. Go also has ABIInternal, which is used internally between Go
functions and runtime assembly and pass arguments and returns via registers. But
we, mere mortals, can't use it. By consequence, we have 5 instructions in this
procedure.

GoAdd assembly:
```
TEXT gcap/addgo.GoAdd(SB) addgo/addgo.go
  addgo.go:5		0x12a77a0		4801d8			ADDQ BX, AX		
  addgo.go:6		0x12a77a3		c3			RET			
```

In this case, this leaf function can't grow the stack, so this procedure is very 
cheap, only 2 instructions.

CAdd assembly:
```
TEXT gcap/addc.CAdd(SB) addc/addc.go
  addc.go:11		0x12a79c0		493b6610		CMPQ SP, 0x10(R14)			
  addc.go:11		0x12a79c4		762e			JBE 0x12a79f4				
  addc.go:11		0x12a79c6		55			PUSHQ BP				
  addc.go:11		0x12a79c7		4889e5			MOVQ SP, BP				
  addc.go:11		0x12a79ca		4883ec18		SUBQ $0x18, SP				
  addc.go:12		0x12a79ce		48890424		MOVQ AX, 0(SP)				
  addc.go:12		0x12a79d2		48895c2408		MOVQ BX, 0x8(SP)			
  addc.go:12		0x12a79d7		e804ffffff		CALL gcap/addc._Cfunc_c_add.abi0(SB)	
  addc.go:12		0x12a79dc		450f57ff		XORPS X15, X15				
  addc.go:12		0x12a79e0		644c8b3425f8ffffff	MOVQ FS:0xfffffff8, R14			
  addc.go:12		0x12a79e9		488b442410		MOVQ 0x10(SP), AX			
  addc.go:12		0x12a79ee		4883c418		ADDQ $0x18, SP				
  addc.go:12		0x12a79f2		5d			POPQ BP					
  addc.go:12		0x12a79f3		c3			RET					
  addc.go:11		0x12a79f4		4889442408		MOVQ AX, 0x8(SP)			
  addc.go:11		0x12a79f9		48895c2410		MOVQ BX, 0x10(SP)			
  addc.go:11		0x12a79fe		6690			NOPW					
  addc.go:11		0x12a7a00		e81ba1f5ff		CALL runtime.morestack_noctxt.abi0(SB)	
  addc.go:11		0x12a7a05		488b442408		MOVQ 0x8(SP), AX			
  addc.go:11		0x12a7a0a		488b5c2410		MOVQ 0x10(SP), BX			
  addc.go:11		0x12a7a0f		ebaf			JMP gcap/addc.CAdd(SB)			


TEXT gcap/addc._Cfunc_c_add.abi0(SB) _cgo_gotypes.go
  _cgo_gotypes.go:47	0x12a78e0		644c8b3425f8ffffff	MOVQ FS:0xfffffff8, R14					
  _cgo_gotypes.go:47	0x12a78e9		493b6610		CMPQ SP, 0x10(R14)					
  _cgo_gotypes.go:47	0x12a78ed		0f86ab000000		JBE 0x12a799e						
  _cgo_gotypes.go:47	0x12a78f3		55			PUSHQ BP						
  _cgo_gotypes.go:47	0x12a78f4		4889e5			MOVQ SP, BP						
  _cgo_gotypes.go:47	0x12a78f7		4883ec18		SUBQ $0x18, SP						
  _cgo_gotypes.go:47	0x12a78fb		48c744243800000000	MOVQ $0x0, 0x38(SP)					
  _cgo_gotypes.go:48	0x12a7904		488d5c2428		LEAQ 0x28(SP), BX					
  _cgo_gotypes.go:48	0x12a7909		48895c2410		MOVQ BX, 0x10(SP)					
  _cgo_gotypes.go:48	0x12a790e		488b0583470100		MOVQ gcap/addc._cgo_2ff6efcb9736_Cfunc_c_add(SB), AX	
  _cgo_gotypes.go:48	0x12a7915		450f57ff		XORPS X15, X15						
  _cgo_gotypes.go:48	0x12a7919		644c8b3425f8ffffff	MOVQ FS:0xfffffff8, R14					
  _cgo_gotypes.go:48	0x12a7922		e8b92bf5ff		CALL runtime.cgocall(SB)				
  _cgo_gotypes.go:49	0x12a7927		803d047d050000		CMPB runtime.cgoAlwaysFalse(SB), $0x0			
  _cgo_gotypes.go:49	0x12a792e		7468			JE 0x12a7998						
  _cgo_gotypes.go:50	0x12a7930		488b442428		MOVQ 0x28(SP), AX					
  _cgo_gotypes.go:50	0x12a7935		450f57ff		XORPS X15, X15						
  _cgo_gotypes.go:50	0x12a7939		644c8b3425f8ffffff	MOVQ FS:0xfffffff8, R14					
  _cgo_gotypes.go:50	0x12a7942		e8f930f5ff		CALL runtime.convT64(SB)				
  _cgo_gotypes.go:50	0x12a7947		4889c3			MOVQ AX, BX						
  _cgo_gotypes.go:50	0x12a794a		488d050fced6ff		LEAQ 0xffd6ce0f(IP), AX					
  _cgo_gotypes.go:50	0x12a7951		450f57ff		XORPS X15, X15						
  _cgo_gotypes.go:50	0x12a7955		644c8b3425f8ffffff	MOVQ FS:0xfffffff8, R14					
  _cgo_gotypes.go:50	0x12a795e		6690			NOPW							
  _cgo_gotypes.go:50	0x12a7960		e83b2bf5ff		CALL runtime.cgoUse(SB)					
  _cgo_gotypes.go:51	0x12a7965		488b442430		MOVQ 0x30(SP), AX					
  _cgo_gotypes.go:51	0x12a796a		450f57ff		XORPS X15, X15						
  _cgo_gotypes.go:51	0x12a796e		644c8b3425f8ffffff	MOVQ FS:0xfffffff8, R14					
  _cgo_gotypes.go:51	0x12a7977		e8c430f5ff		CALL runtime.convT64(SB)				
  _cgo_gotypes.go:51	0x12a797c		4889c3			MOVQ AX, BX						
  _cgo_gotypes.go:51	0x12a797f		488d05dacdd6ff		LEAQ 0xffd6cdda(IP), AX					
  _cgo_gotypes.go:51	0x12a7986		450f57ff		XORPS X15, X15						
  _cgo_gotypes.go:51	0x12a798a		644c8b3425f8ffffff	MOVQ FS:0xfffffff8, R14					
  _cgo_gotypes.go:51	0x12a7993		e8082bf5ff		CALL runtime.cgoUse(SB)					
  _cgo_gotypes.go:53	0x12a7998		4883c418		ADDQ $0x18, SP						
  _cgo_gotypes.go:53	0x12a799c		5d			POPQ BP							
  _cgo_gotypes.go:53	0x12a799d		c3			RET							
  _cgo_gotypes.go:47	0x12a799e		6690			NOPW							
  _cgo_gotypes.go:47	0x12a79a0		e87ba1f5ff		CALL runtime.morestack_noctxt.abi0(SB)			
  _cgo_gotypes.go:47	0x12a79a5		e936ffffff		JMP gcap/addc._Cfunc_c_add.abi0(SB)			
TEXT _cgo_2ff6efcb9736_Cfunc_c_add(SB) 
  :0			0x12b8710		55			PUSHQ BP			
  :0			0x12b8711		4889e5			MOVQ SP, BP			
  :0			0x12b8714		4157			PUSHQ R15			
  :0			0x12b8716		4156			PUSHQ R14			
  :0			0x12b8718		53			PUSHQ BX			
  :0			0x12b8719		50			PUSHQ AX			
  :0			0x12b871a		4889fb			MOVQ DI, BX			
  :0			0x12b871d		e87eaef4ff		CALL _cgo_topofstack(SB)	
  :0			0x12b8722		4989c6			MOVQ AX, R14			
  :0			0x12b8725		4c8b7b08		MOVQ 0x8(BX), R15		
  :0			0x12b8729		4c033b			ADDQ 0(BX), R15			
  :0			0x12b872c		e86faef4ff		CALL _cgo_topofstack(SB)	
  :0			0x12b8731		4c29f0			SUBQ R14, AX			
  :0			0x12b8734		4c897c0310		MOVQ R15, 0x10(BX)(AX*1)	
  :0			0x12b8739		4883c408		ADDQ $0x8, SP			
  :0			0x12b873d		5b			POPQ BX				
  :0			0x12b873e		415e			POPQ R14			
  :0			0x12b8740		415f			POPQ R15			
  :0			0x12b8742		5d			POPQ BP				
  :0			0x12b8743		c3			RET				
```

And this is the monster generated by the CGO call. There's a **lot** of things
going on here. Which adds overhead. There are reasons for that, CGO must handle
basically infinite cases. If we know what we are doing, we can just use 
assembly to call into C.

Benchmark run0
```
BenchmarkWarmup
BenchmarkWarmup-4                  50000               223.4 ns/op             0 B/op          0 allocs/op
BenchmarkAsmAdd
BenchmarkAsmAdd-4                  50000                 3.297 ns/op           0 B/op          0 allocs/op
BenchmarkCAdd
BenchmarkCAdd-4                    50000                62.45 ns/op            0 B/op          0 allocs/op
BenchmarkGoAdd
BenchmarkGoAdd-4                   50000                 4.367 ns/op           0 B/op          0 allocs/op
BenchmarkAsmAddParallel
BenchmarkAsmAddParallel-4          50000                 3.059 ns/op           0 B/op          0 allocs/op
BenchmarkCAddParallel
BenchmarkCAddParallel-4            50000                32.18 ns/op            0 B/op          0 allocs/op
BenchmarkGoAddParallel
BenchmarkGoAddParallel-4           50000                 4.430 ns/op           0 B/op          0 allocs/op
PASS
```

Benchmark run1
```
BenchmarkWarmup
BenchmarkWarmup-4                  50000               223.2 ns/op             0 B/op          0 allocs/op
BenchmarkAsmAdd
BenchmarkAsmAdd-4                  50000                 3.214 ns/op           0 B/op          0 allocs/op
BenchmarkCAdd
BenchmarkCAdd-4                    50000                59.12 ns/op            0 B/op          0 allocs/op
BenchmarkGoAdd
BenchmarkGoAdd-4                   50000                 3.475 ns/op           0 B/op          0 allocs/op
BenchmarkAsmAddParallel
BenchmarkAsmAddParallel-4          50000                 9.453 ns/op           0 B/op          0 allocs/op
BenchmarkCAddParallel
BenchmarkCAddParallel-4            50000                32.18 ns/op            0 B/op          0 allocs/op
BenchmarkGoAddParallel
BenchmarkGoAddParallel-4           50000                 4.606 ns/op           0 B/op          0 allocs/op
PASS
```

Benchmark run2
```
BenchmarkWarmup
BenchmarkWarmup-4                  50000               222.1 ns/op             0 B/op          0 allocs/op
BenchmarkAsmAdd
BenchmarkAsmAdd-4                  50000                 3.214 ns/op           0 B/op          0 allocs/op
BenchmarkCAdd
BenchmarkCAdd-4                    50000                62.38 ns/op            0 B/op          0 allocs/op
BenchmarkGoAdd
BenchmarkGoAdd-4                   50000                 3.214 ns/op           0 B/op          0 allocs/op
BenchmarkAsmAddParallel
BenchmarkAsmAddParallel-4          50000                10.02 ns/op            0 B/op          0 allocs/op
BenchmarkCAddParallel
BenchmarkCAddParallel-4            50000                37.94 ns/op            0 B/op          0 allocs/op
BenchmarkGoAddParallel
BenchmarkGoAddParallel-4           50000                 4.658 ns/op           0 B/op          0 allocs/op
PASS
```

There's a lot of deviation in the benchmarks, given that we have the Go
scheduler/runtime doing it's work in the background and also my operating 
system and other programs running in my machine.

In summary

| Implementation        | Average Latency   |
|-----------------------|-----------------  |
| Go Serialized         | 3.68 ms           |
| Assembly Serialized   | 3.24 ms           |
| C Serialized          | 61,31 ms          |
| Go Parallel           | 4,56ms            |
| Assembly Parallel     | 7.51ms            |
| C Parallel            | 34.1ms            |


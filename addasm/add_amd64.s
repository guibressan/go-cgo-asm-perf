#include "textflag.h"

TEXT ·AsmAdd(SB),NOSPLIT,$0-24 
	MOVQ a+0(FP), AX
	MOVQ b+8(FP), BX
	ADDQ BX, AX
	MOVQ AX, r+16(FP)
	RET

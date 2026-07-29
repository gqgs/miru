//go:build amd64 && !noasm && !appengine

#include "textflag.h"

// func compareHistAVX512(h1r, h2r, h1g, h2g, h1b, h2b, result unsafe.Pointer)
TEXT ·compareHistAVX512(SB), NOSPLIT, $0-56
	MOVQ h1r+0(FP), DI
	MOVQ h2r+8(FP), SI
	MOVQ h1g+16(FP), DX
	MOVQ h2g+24(FP), CX
	MOVQ h1b+32(FP), R8
	MOVQ h2b+40(FP), R9
	MOVQ result+48(FP), R10

	VPXORD Z0, Z0, Z0
	VPXORD Z15, Z15, Z15
	KXNORW K1, K1, K1
	XORQ AX, AX

loop:
	// Red channel.
	VMOVUPS (DI)(AX*1), Z1
	VMOVUPS (SI)(AX*1), Z2
	VADDPS Z2, Z1, Z3
	VSUBPS Z2, Z1, Z4
	VMULPS Z4, Z4, Z4
	VCMPPS $12, Z15, Z3, K1, K2
	VDIVPS Z3, Z4, K2, Z4
	VADDPS Z4, Z0, Z0

	// Green channel.
	VMOVUPS (DX)(AX*1), Z1
	VMOVUPS (CX)(AX*1), Z2
	VADDPS Z2, Z1, Z3
	VSUBPS Z2, Z1, Z4
	VMULPS Z4, Z4, Z4
	VCMPPS $12, Z15, Z3, K1, K2
	VDIVPS Z3, Z4, K2, Z4
	VADDPS Z4, Z0, Z0

	// Blue channel.
	VMOVUPS (R8)(AX*1), Z1
	VMOVUPS (R9)(AX*1), Z2
	VADDPS Z2, Z1, Z3
	VSUBPS Z2, Z1, Z4
	VMULPS Z4, Z4, Z4
	VCMPPS $12, Z15, Z3, K1, K2
	VDIVPS Z3, Z4, K2, Z4
	VADDPS Z4, Z0, Z0

	ADDQ $64, AX
	CMPQ AX, $1024
	JB loop

	VMOVUPS Z0, (R10)
	VZEROUPPER
	RET

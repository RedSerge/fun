alias Euler

alias zerosd (reg) {
	pxor reg, reg
}

alias reset () {
	mov count, len
	mov psrc, src
	mov pdst, dst
	alias zerosd (temp1)
}

Euler:

	alias [linux] {
		alias src, rdi
		alias dst, rsi
		alias len, edx
	}
	
	alias [windows] {
		alias src, rcx
		alias dst, rdx
		alias len, r8d
	}
	
	alias EulerCount, xmm0
	alias EulerStep, xmm1
	alias EulerMax, xmm2
	alias EulerHalf, xmm3
	
	alias temp1, xmm4
	alias temp2, xmm5
	
	alias count, r9d
	alias psrc, r10
	alias pdst, r11
	
	alias Kolmogorov, L10
	alias Leave_Kolmogorov, L11
	alias Kolmogorov2, L12
	alias Leave_Kolmogorov2, L13
	alias QuitProcedure, L14
	
	movsd EulerStep, [dst]
	movsd EulerMax, 8[dst]
	movsd EulerHalf, 16[dst]
	
	alias zerosd (EulerCount)
	alias reset()
	
.Kolmogorov:

	movsd temp2, 16[psrc]
	movsd [pdst], temp2
	add pdst, 8
	
	movsd temp2, [psrc]
	addsd temp2, 8[psrc]
	mulsd temp2, 16[psrc]
	subsd temp1, temp2

	add psrc, 24
	sub count, 1
	
	cmp count, 1
	jl .Leave_Kolmogorov
	
	movsd temp2, 8[psrc]
	mulsd temp2, 16[psrc]
	addsd temp2, temp1
	
	movsd [pdst], temp2
	add pdst, 8
	
	movsd temp1, -16[pdst]
	mulsd temp2, EulerStep
	addsd temp1, temp2
	movsd [pdst], temp1
	add pdst, 8
	
	movsd temp1, -24[psrc]
	mulsd temp1, -8[psrc]
	
	jmp .Kolmogorov
.Leave_Kolmogorov:

	alias reset()
	
.Kolmogorov2:
	
	movsd temp2, [psrc]
	addsd temp2, 8[psrc]
	mulsd temp2, 16[pdst]
	subsd temp1, temp2

	add psrc, 24
	add pdst, 24
	sub count, 1
	
	cmp count, 1
	jl .Leave_Kolmogorov2
	
	movsd temp2, 8[psrc]
	mulsd temp2, 16[pdst]
	addsd temp2, temp1
	
	addsd temp2, -16[pdst]
	mulsd temp2, EulerHalf
	addsd temp2, -24[pdst]
	
	movsd -8[psrc], temp2
	
	movsd temp1, -24[psrc]
	mulsd temp1, -8[pdst]
	
	jmp .Kolmogorov2
.Leave_Kolmogorov2:
	
	addsd EulerCount, EulerStep
	ucomisd EulerMax, EulerCount
	jbe .QuitProcedure

	alias reset()
	jmp .Kolmogorov
	
.QuitProcedure:
	movsd [dst], EulerStep
	movsd 8[dst], EulerMax
	movsd 16[dst], EulerHalf
	ret
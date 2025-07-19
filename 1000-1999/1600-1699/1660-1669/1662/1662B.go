package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min3(a, b, c int) int {
	return min(a, min(b, c))
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func solve(ABC, AB, AC, BC, A, B, C []int) [][2]int {
	var ret [][2]int
	for _, ch := range ABC {
		ret = append(ret, [2]int{ch, ch})
	}
	for len(BC) > 0 && len(A) > 0 {
		u := BC[len(BC)-1]
		v := A[len(A)-1]
		ret = append(ret, [2]int{u, v})
		BC = BC[:len(BC)-1]
		A = A[:len(A)-1]
	}
	for _, ch := range BC {
		B = append(B, ch)
		C = append(C, ch)
	}
	for len(AC) > 0 && len(B) > 0 {
		u := AC[len(AC)-1]
		v := B[len(B)-1]
		ret = append(ret, [2]int{u, v})
		AC = AC[:len(AC)-1]
		B = B[:len(B)-1]
	}
	for _, ch := range AC {
		A = append(A, ch)
		C = append(C, ch)
	}
	for len(AB) > 0 && len(C) > 0 {
		u := AB[len(AB)-1]
		v := C[len(C)-1]
		ret = append(ret, [2]int{u, v})
		AB = AB[:len(AB)-1]
		C = C[:len(C)-1]
	}
	for _, ch := range AB {
		A = append(A, ch)
		B = append(B, ch)
	}
	for (boolToInt(len(A) == 0) + boolToInt(len(B) == 0) + boolToInt(len(C) == 0)) < 2 {
		if len(C) == min3(len(A), len(B), len(C)) {
			u := A[len(A)-1]
			v := B[len(B)-1]
			ret = append(ret, [2]int{u, v})
			A = A[:len(A)-1]
			B = B[:len(B)-1]
		} else if len(B) == min3(len(A), len(B), len(C)) {
			u := A[len(A)-1]
			v := C[len(C)-1]
			ret = append(ret, [2]int{u, v})
			A = A[:len(A)-1]
			C = C[:len(C)-1]
		} else {
			u := B[len(B)-1]
			v := C[len(C)-1]
			ret = append(ret, [2]int{u, v})
			B = B[:len(B)-1]
			C = C[:len(C)-1]
		}
	}
	for _, ch := range A {
		ret = append(ret, [2]int{ch, ch})
	}
	for _, ch := range B {
		ret = append(ret, [2]int{ch, ch})
	}
	for _, ch := range C {
		ret = append(ret, [2]int{ch, ch})
	}
	return ret
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b, c string
	fmt.Fscan(reader, &a, &b, &c)
	cnta := make([]int, 26)
	cntb := make([]int, 26)
	cntc := make([]int, 26)
	for i := 0; i < len(a); i++ {
		cnta[a[i]-'A']++
	}
	for i := 0; i < len(b); i++ {
		cntb[b[i]-'A']++
	}
	for i := 0; i < len(c); i++ {
		cntc[c[i]-'A']++
	}
	var ABC, AB, AC, BC, A_list, B_list, C_list []int
	for i := 0; i < 26; i++ {
		ins := func(vec *[]int, x int) {
			for j := 0; j < x; j++ {
				*vec = append(*vec, i)
			}
		}
		fuck := func(ABCp, XYp, Xp, XZp, YZp *[]int, cntX, cntY, cntZ int) {
			mn := min(cntZ, cntX-cntY)
			ins(ABCp, (cntZ-mn)%2)
			ins(XYp, cntY-cntZ+mn+(cntZ-mn)/2)
			ins(Xp, cntX-cntY-mn)
			ins(XZp, mn+(cntZ-mn)/2)
			ins(YZp, (cntZ-mn)/2)
		}
		mx := max(cnta[i], max(cntb[i], cntc[i]))
		if cnta[i] == mx {
			if cntb[i] > cntc[i] {
				fuck(&ABC, &AB, &A_list, &AC, &BC, cnta[i], cntb[i], cntc[i])
			} else {
				fuck(&ABC, &AC, &A_list, &AB, &BC, cnta[i], cntc[i], cntb[i])
			}
		} else if cntb[i] == mx {
			if cnta[i] > cntc[i] {
				fuck(&ABC, &AB, &B_list, &BC, &AC, cntb[i], cnta[i], cntc[i])
			} else {
				fuck(&ABC, &BC, &B_list, &AB, &AC, cntb[i], cntc[i], cnta[i])
			}
		} else {
			if cnta[i] > cntb[i] {
				fuck(&ABC, &AC, &C_list, &BC, &AB, cntc[i], cnta[i], cntb[i])
			} else {
				fuck(&ABC, &BC, &C_list, &AC, &AB, cntc[i], cntb[i], cnta[i])
			}
		}
	}
	vec := solve(ABC, AB, AC, BC, A_list, B_list, C_list)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, len(vec))
	for _, uv := range vec {
		u, v := uv[0], uv[1]
		writer.WriteByte(byte(u + 'A'))
		writer.WriteByte(byte(v + 'A'))
		writer.WriteByte('\n')
	}
}

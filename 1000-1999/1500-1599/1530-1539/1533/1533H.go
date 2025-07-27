package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func countSubmatrices(matrix [][]byte, mask int) int64 {
	n := len(matrix)
	m := len(matrix[0])
	heights := make([]int, m)
	stackH := make([]int, m)
	stackC := make([]int, m)
	var total int64
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if (1<<(matrix[i][j]-'A'))&mask != 0 {
				heights[j]++
			} else {
				heights[j] = 0
			}
		}
		top := 0
		var sum int64
		for j := 0; j < m; j++ {
			h := heights[j]
			cnt := 1
			for top > 0 && stackH[top-1] >= h {
				cnt += stackC[top-1]
				sum -= int64(stackH[top-1]) * int64(stackC[top-1])
				top--
			}
			sum += int64(h) * int64(cnt)
			stackH[top] = h
			stackC[top] = cnt
			top++
			total += sum
		}
	}
	return total
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	matrix := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		matrix[i] = []byte(s)
	}
	F := make([]int64, 32)
	for mask := 1; mask < 32; mask++ {
		F[mask] = countSubmatrices(matrix, mask)
	}
	exact := make([]int64, 32)
	for mask := 1; mask < 32; mask++ {
		exact[mask] = F[mask]
		for sub := (mask - 1) & mask; sub > 0; sub = (sub - 1) & mask {
			exact[mask] -= exact[sub]
		}
	}
	ans := make([]int64, 6)
	for mask := 1; mask < 32; mask++ {
		k := bits.OnesCount(uint(mask))
		ans[k] += exact[mask]
	}
	for k := 1; k <= 5; k++ {
		if k > 1 {
			fmt.Print(" ")
		}
		fmt.Print(ans[k])
	}
	fmt.Println()
}

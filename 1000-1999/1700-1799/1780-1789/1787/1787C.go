package main

import (
	"bufio"
	"fmt"
	"os"
)

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, s int
		fmt.Fscan(reader, &n, &s)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// precompute two possible splits for each internal element
		xA := make([]int, n)
		yA := make([]int, n)
		xB := make([]int, n)
		yB := make([]int, n)
		for i := 1; i <= n-2; i++ {
			ai := a[i]
			L := s
			if ai < s {
				L = ai
			}
			R := 0
			if ai > s {
				R = ai - s
			}
			xA[i] = L
			yA[i] = ai - L
			xB[i] = R
			yB[i] = ai - R
		}
		dpA := int64(a[0] * xA[1])
		dpB := int64(a[0] * xB[1])
		for i := 2; i <= n-2; i++ {
			aiA := xA[i]
			aiB := xB[i]
			yPrevA := yA[i-1]
			yPrevB := yB[i-1]
			newA := min64(dpA+int64(yPrevA*aiA), dpB+int64(yPrevB*aiA))
			newB := min64(dpA+int64(yPrevA*aiB), dpB+int64(yPrevB*aiB))
			dpA, dpB = newA, newB
		}
		yPrevA := yA[n-2]
		yPrevB := yB[n-2]
		ans := min64(dpA+int64(yPrevA*a[n-1]), dpB+int64(yPrevB*a[n-1]))
		fmt.Fprintln(writer, ans)
	}
}

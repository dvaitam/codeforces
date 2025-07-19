package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &l[i], &r[i])
	}

	const INF = 1000000005
	L, R := -INF, INF
	posL, posR := 0, 0
	for i := 0; i < n; i++ {
		if l[i] > L {
			L = l[i]
			posL = i
		}
		if r[i] < R {
			R = r[i]
			posR = i
		}
	}
	ans := 0
	if R-L > ans {
		ans = R - L
	}
	// exclude posL
	L1, R1 := -INF, INF
	for i := 0; i < n; i++ {
		if i == posL {
			continue
		}
		if l[i] > L1 {
			L1 = l[i]
		}
		if r[i] < R1 {
			R1 = r[i]
		}
	}
	if R1-L1 > ans {
		ans = R1 - L1
	}
	// exclude posR
	L2, R2 := -INF, INF
	for i := 0; i < n; i++ {
		if i == posR {
			continue
		}
		if l[i] > L2 {
			L2 = l[i]
		}
		if r[i] < R2 {
			R2 = r[i]
		}
	}
	if R2-L2 > ans {
		ans = R2 - L2
	}
	if ans < 0 {
		ans = 0
	}
	fmt.Fprintln(writer, ans)
}

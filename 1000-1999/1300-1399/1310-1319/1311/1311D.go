package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	INF  = 1000000000
	MAXA = 10000
)

var (
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	defer out.Flush()
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		solve()
	}
}

func solve() {
	var a, b, c int
	fmt.Fscan(in, &a, &b, &c)
	ans := INF
	var ansA, ansB, ansC int

	// Process A <= a
	for A := a; A >= 1; A-- {
		da := a - A
		if da > ans {
			break
		}
		processA(A, a, b, c, da, &ans, &ansA, &ansB, &ansC)
	}
	// Process A > a
	for A := a + 1; A <= MAXA; A++ {
		da := A - a
		if da > ans {
			break
		}
		processA(A, a, b, c, da, &ans, &ansA, &ansB, &ansC)
	}
	// Output result
	fmt.Fprintln(out, ans)
	fmt.Fprintf(out, "%d %d %d\n", ansA, ansB, ansC)
}

func processA(A, a, b, c, costA int, ans, ansA, ansB, ansC *int) {
	// Consider B = A * k around b
	k := b / A
	// B lower (k)
	if k >= 1 {
		B := A * k
		costB := abs(B - b)
		newNum := costA + costB
		if newNum <= *ans {
			processB(A, B, c, newNum, ans, ansA, ansB, ansC)
		}
	}
	// B upper (k+1)
	B2 := A * (k + 1)
	costB2 := abs(B2 - b)
	newNum2 := costA + costB2
	if newNum2 <= *ans {
		processB(A, B2, c, newNum2, ans, ansA, ansB, ansC)
	}
}

func processB(A, B, c, num int, ans, ansA, ansB, ansC *int) {
	// Consider C = B * j around c
	jc := c / B
	// C lower (jc)
	if jc >= 1 {
		C := B * jc
		total := num + abs(C-c)
		if total < *ans {
			*ans = total
			*ansA = A
			*ansB = B
			*ansC = C
		}
	}
	// C upper (jc+1)
	C2 := B * (jc + 1)
	total2 := num + abs(C2-c)
	if total2 < *ans {
		*ans = total2
		*ansA = A
		*ansB = B
		*ansC = C2
	}
}

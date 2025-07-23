package main

import (
	"bufio"
	"fmt"
	"os"
)

func commonDigit(a, b [2]int) (int, bool) {
	count := 0
	var val int
	if a[0] == b[0] || a[0] == b[1] {
		count++
		val = a[0]
	}
	if a[1] == b[0] || a[1] == b[1] {
		count++
		val = a[1]
	}
	if count == 1 {
		return val, true
	}
	return 0, false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	A := make([][2]int, n)
	B := make([][2]int, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &A[i][0], &A[i][1])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &B[i][0], &B[i][1])
	}

	var union [10]bool
	aDigits := make([][10]bool, n)
	bDigits := make([][10]bool, m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if d, ok := commonDigit(A[i], B[j]); ok {
				union[d] = true
				aDigits[i][d] = true
				bDigits[j][d] = true
			}
		}
	}

	ans := 0
	cnt := 0
	for d := 1; d <= 9; d++ {
		if union[d] {
			cnt++
			ans = d
		}
	}
	if cnt == 1 {
		fmt.Fprintln(writer, ans)
		return
	}

	for i := 0; i < n; i++ {
		c := 0
		for d := 1; d <= 9; d++ {
			if aDigits[i][d] {
				c++
			}
		}
		if c > 1 {
			fmt.Fprintln(writer, -1)
			return
		}
	}

	for j := 0; j < m; j++ {
		c := 0
		for d := 1; d <= 9; d++ {
			if bDigits[j][d] {
				c++
			}
		}
		if c > 1 {
			fmt.Fprintln(writer, -1)
			return
		}
	}

	fmt.Fprintln(writer, 0)
}

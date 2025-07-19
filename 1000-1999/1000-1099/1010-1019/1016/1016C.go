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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n+2)
	b := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		a[i] = x
	}
	for i := 1; i <= n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		b[i] = x
	}

	s := make([]int64, n+3)
	for i := n; i >= 1; i-- {
		s[i] = a[i] + b[i] + s[i+1]
	}

	var nowL, nowR1, nowR2, ans int64
	var k int64
	// compute initial nowR1: a[1..n], then b[n..1]
	for i := 1; i <= n; i++ {
		nowR1 += a[i] * k
		k++
	}
	for i := n; i >= 1; i-- {
		nowR1 += b[i] * k
		k++
	}
	// compute initial nowR2: b[1..n], then a[n..1]
	k = 0
	for i := 1; i <= n; i++ {
		nowR2 += b[i] * k
		k++
	}
	for i := n; i >= 1; i-- {
		nowR2 += a[i] * k
		k++
	}

	nowL = 0
	ans = 0
	k = 0
	// iterate and update
	totalSteps := int64(n*2 - 1)
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			if nowL+nowR1 > ans {
				ans = nowL + nowR1
			}
			nowL += a[i] * k
			k++
			nowL += b[i] * k
			k++
		} else {
			if nowL+nowR2 > ans {
				ans = nowL + nowR2
			}
			nowL += b[i] * k
			k++
			nowL += a[i] * k
			k++
		}
		// update right sums
		nowR1 -= totalSteps * b[i]
		nowR1 -= int64(i-1) * 2 * a[i]
		nowR1 += s[i+1]

		nowR2 -= totalSteps * a[i]
		nowR2 -= int64(i-1) * 2 * b[i]
		nowR2 += s[i+1]
	}
	fmt.Fprintln(writer, ans)
}

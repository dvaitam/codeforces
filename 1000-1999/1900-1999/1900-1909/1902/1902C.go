package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int64, n)
	maxv := int64(-1 << 63)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > maxv {
			maxv = a[i]
		}
	}
	if n == 1 {
		fmt.Fprintln(writer, 1)
		return
	}
	g := int64(0)
	for i := 0; i < n; i++ {
		g = gcd(g, maxv-a[i])
	}
	sumOps := int64(0)
	set := make(map[int64]struct{}, n)
	for i := 0; i < n; i++ {
		sumOps += (maxv - a[i]) / g
		set[a[i]] = struct{}{}
	}
	k := int64(1)
	for {
		candidate := maxv - k*g
		if _, ok := set[candidate]; !ok {
			break
		}
		k++
	}
	fmt.Fprintln(writer, sumOps+k)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}

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
	return a
}

func solve(r *bufio.Reader, w *bufio.Writer) {
	var n int
	fmt.Fscan(r, &n)
	arr := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &arr[i])
		total += arr[i]
	}
	var pref int64
	best := int64(1)
	for i := 0; i < n-1; i++ {
		pref += arr[i]
		g := gcd(pref, total)
		if g > best {
			best = g
		}
	}
	fmt.Fprintln(w, best)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}

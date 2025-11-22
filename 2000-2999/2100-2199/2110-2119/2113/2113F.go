package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		lim := 2*n + 2
		seenA := make([]bool, lim)
		seenB := make([]bool, lim)

		for i := 0; i < n; i++ {
			x, y := a[i], b[i]
			gainOrig := 0
			if !seenA[x] {
				gainOrig++
			}
			if !seenB[y] {
				gainOrig++
			}

			gainSwap := 0
			if !seenA[y] {
				gainSwap++
			}
			if !seenB[x] {
				gainSwap++
			}

			if gainSwap > gainOrig {
				a[i], b[i] = b[i], a[i]
				x, y = a[i], b[i]
			}

			seenA[x] = true
			seenB[y] = true
		}

		ans := 0
		for i := 0; i < lim; i++ {
			if seenA[i] {
				ans++
			}
			if seenB[i] {
				ans++
			}
		}

		fmt.Fprintln(out, ans)
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, a[i])
		}
		fmt.Fprintln(out)
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, b[i])
		}
		fmt.Fprintln(out)
	}
}

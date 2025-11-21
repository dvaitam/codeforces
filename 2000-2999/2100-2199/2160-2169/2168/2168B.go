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

	var phase string
	if _, err := fmt.Fscan(in, &phase); err != nil {
		return
	}

	switch phase {
	case "first":
		handleFirstRun(in, out)
	case "second":
		handleSecondRun(in, out)
	}
}

func handleFirstRun(in *bufio.Reader, out *bufio.Writer) {
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		pos1, posN := -1, -1
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)
			if v == 1 {
				pos1 = i
			}
			if v == n {
				posN = i
			}
		}
		if pos1 < posN {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, 1)
		}
	}
}

func handleSecondRun(in *bufio.Reader, out *bufio.Writer) {
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(in, &n, &x)
		ask := func(l, r int) int {
			fmt.Fprintf(out, "? %d %d\n", l, r)
			out.Flush()
			var res int
			if _, err := fmt.Fscan(in, &res); err != nil {
				os.Exit(0)
			}
			if res == -1 {
				os.Exit(0)
			}
			return res
		}

		pos := 0
		if x == 0 {
			pos = searchPrefix(n, ask)
		} else {
			pos = searchSuffix(n, ask)
		}

		fmt.Fprintf(out, "! %d\n", pos)
		out.Flush()
	}
}

func searchPrefix(n int, ask func(int, int) int) int {
	// Smallest prefix whose range difference hits n-1 already contains both 1 and n.
	lo, hi := 1, n
	for lo < hi {
		mid := (lo + hi) / 2
		if ask(1, mid) == n-1 {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func searchSuffix(n int, ask func(int, int) int) int {
	// Smallest suffix start whose range difference hits n-1 also captures both extremes.
	lo, hi := 1, n
	for lo < hi {
		mid := (lo + hi) / 2
		if ask(mid, n) == n-1 {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

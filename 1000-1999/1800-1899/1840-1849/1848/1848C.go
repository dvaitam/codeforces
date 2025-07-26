package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// residue returns the first step modulo 3 at which a becomes zero during the
// process described in problemC. If both a and b are zero, the pair imposes no
// restriction and ok is false.
func residue(a, b int64) (int, bool) {
	if a == 0 && b == 0 {
		return 0, false
	}
	g := gcd64(a, b)
	a /= g
	b /= g
	if a%2 == 0 {
		return 0, true
	}
	if b%2 == 0 {
		return 1, true
	}
	return 2, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		first := true
		res := 0
		ok := true
		for i := 0; i < n; i++ {
			r, used := residue(a[i], b[i])
			if !used {
				continue
			}
			if first {
				first = false
				res = r
			} else if r != res {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

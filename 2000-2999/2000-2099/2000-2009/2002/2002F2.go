package main

import (
	"bufio"
	"fmt"
	"os"
)

const searchLimit = 256

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func bestCoprime(maxVal, other int) int {
	if maxVal <= 0 {
		return 0
	}
	start := maxVal
	minVal := start - searchLimit + 1
	if minVal < 1 {
		minVal = 1
	}
	for v := start; v >= minVal; v-- {
		if gcd(v, other) == 1 {
			return v
		}
	}
	for v := minVal - 1; v >= 1; v-- {
		if gcd(v, other) == 1 {
			return v
		}
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		var lVal, fVal int64
		fmt.Fscan(in, &n, &m, &lVal, &fVal)

		best := int64(0)

		seenA := make(map[int]struct{})
		processA := func(a int) {
			if a < 1 || a > n {
				return
			}
			if _, ok := seenA[a]; ok {
				return
			}
			seenA[a] = struct{}{}
			b := bestCoprime(m, a)
			if b == 0 {
				return
			}
			val := int64(a)*lVal + int64(b)*fVal
			if val > best {
				best = val
			}
		}

		startA := n - searchLimit + 1
		if startA < 1 {
			startA = 1
		}
		for a := startA; a <= n; a++ {
			processA(a)
		}
		processA(1)
		processA(2)

		seenB := make(map[int]struct{})
		processB := func(b int) {
			if b < 1 || b > m {
				return
			}
			if _, ok := seenB[b]; ok {
				return
			}
			seenB[b] = struct{}{}
			a := bestCoprime(n, b)
			if a == 0 {
				return
			}
			val := int64(a)*lVal + int64(b)*fVal
			if val > best {
				best = val
			}
		}

		startB := m - searchLimit + 1
		if startB < 1 {
			startB = 1
		}
		for b := startB; b <= m; b++ {
			processB(b)
		}
		processB(1)
		processB(2)

		fmt.Fprintln(out, best)
	}
}

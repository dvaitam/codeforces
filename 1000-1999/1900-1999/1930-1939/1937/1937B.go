package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solve(n int, a1, a2 string) (string, int) {
	b1 := []byte(a1)
	b2 := []byte(a2)
	res := make([]byte, 0, n+1)
	res = append(res, b1[0])

	l1, r1 := 1, 0 // prefix interval [l1,r1], empty initially
	l2, r2 := 1, n // suffix interval [l2,r2]

	for i := 2; i <= n+1; i++ {
		if l2 <= i-1 {
			if l1 > r1 {
				l1 = l2
				r1 = i - 1
			} else {
				if l2 <= r1+1 {
					r1 = i - 1
				} else {
					l1 = l2
					r1 = i - 1
				}
			}
			l2 = i
		}

		prefixExists := l1 <= r1
		suffixExists := l2 <= r2
		pc := byte('2')
		sc := byte('2')
		if prefixExists {
			pc = b2[i-2]
		}
		if suffixExists {
			sc = b1[i-1]
		}
		minc := pc
		if sc < minc {
			minc = sc
		}
		res = append(res, minc)

		if prefixExists && pc == minc {
			if r1 > i-1 {
				r1 = i - 1
			}
		} else {
			l1, r1 = 1, 0
		}

		if suffixExists && sc == minc {
			if l2 < i {
				l2 = i
			}
		} else {
			l2 = r2 + 1
		}
	}

	count := 0
	if l1 <= r1 {
		count += r1 - l1 + 1
	}
	if l2 <= r2 {
		count += r2 - l2 + 1
	}
	if l1 <= r1 && l2 <= r2 {
		overlap := min(r1, r2) - max(l1, l2) + 1
		if overlap > 0 {
			count -= overlap
		}
	}

	return string(res), count
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
		var s1, s2 string
		fmt.Fscan(in, &s1)
		fmt.Fscan(in, &s2)
		ans, cnt := solve(n, s1, s2)
		fmt.Fprintln(out, ans)
		fmt.Fprintln(out, cnt)
	}
}

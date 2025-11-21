package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCase(n, m, k int) (bool, string) {
	if k > max(n, m) {
		return false, ""
	}
	lower := max(0, max(k-m, n-m))
	upper := min(k, min(n, n-m+k))
	if lower > upper {
		return false, ""
	}
	dMax := upper
	dMin := dMax - k

	cur := 0
	zeros := n
	ones := m
	desired := n - m
	res := make([]byte, 0, n+m)

	maxHit := dMax == 0
	minHit := dMin == 0

	for cur < dMax {
		if zeros == 0 {
			return false, ""
		}
		res = append(res, '0')
		cur++
		zeros--
	}
	if cur == dMax {
		maxHit = true
	}

	for zeros > 0 || ones > 0 {
		if !minHit && cur > dMin {
			if ones == 0 {
				return false, ""
			}
			res = append(res, '1')
			cur--
			ones--
			if cur == dMin {
				minHit = true
			}
			continue
		}
		if cur > desired {
			if ones == 0 {
				return false, ""
			}
			res = append(res, '1')
			cur--
			ones--
			if cur == dMin {
				minHit = true
			}
			continue
		}
		if cur < desired {
			if zeros == 0 || cur == dMax {
				return false, ""
			}
			res = append(res, '0')
			cur++
			zeros--
			if cur == dMax {
				maxHit = true
			}
			continue
		}
		if zeros > 0 && cur < dMax {
			res = append(res, '0')
			cur++
			zeros--
			if cur == dMax {
				maxHit = true
			}
			continue
		}
		if ones > 0 && cur > dMin {
			res = append(res, '1')
			cur--
			ones--
			if cur == dMin {
				minHit = true
			}
			continue
		}
		return false, ""
	}

	if cur != desired || !maxHit || !minHit {
		return false, ""
	}

	return true, string(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		ok, s := solveCase(n, m, k)
		if !ok {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, s)
		}
	}
}

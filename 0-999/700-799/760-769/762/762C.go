package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b string
	if _, err := fmt.Fscan(in, &a, &b); err != nil {
		return
	}
	n := len(a)
	m := len(b)
	pref := make([]int, m)
	pos := 0
	for i := 0; i < m; i++ {
		for pos < n && a[pos] != b[i] {
			pos++
		}
		if pos == n {
			pref[i] = n
		} else {
			pref[i] = pos
			pos++
		}
	}
	suff := make([]int, m+1)
	suff[m] = n
	pos = n - 1
	for i := m - 1; i >= 0; i-- {
		for pos >= 0 && a[pos] != b[i] {
			pos--
		}
		if pos < 0 {
			suff[i] = -1
		} else {
			suff[i] = pos
			pos--
		}
	}
	bestL, bestR, bestLen := 0, m, m
	r := 0
	for l := 0; l <= m; l++ {
		prefixPos := -1
		if l > 0 {
			if pref[l-1] == n {
				break
			}
			prefixPos = pref[l-1]
		}
		if r < l {
			r = l
		}
		for r <= m {
			if r == m {
				break
			}
			if suff[r] == -1 || suff[r] <= prefixPos {
				r++
				continue
			}
			break
		}
		if r > m {
			break
		}
		delLen := r - l
		if delLen < bestLen {
			bestLen = delLen
			bestL = l
			bestR = r
		}
	}
	res := b[:bestL] + b[bestR:]
	if len(res) == 0 {
		fmt.Println("-")
	} else {
		fmt.Println(res)
	}
}

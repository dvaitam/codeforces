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
	var T int
	fmt.Fscan(reader, &T)
	for tc := 0; tc < T; tc++ {
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n, Q int
	fmt.Fscan(reader, &n, &Q)
	v := make([]int, n+2)
	psum := make([]int64, n+2)
	xsum := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &v[i])
		psum[i] = psum[i-1] + int64(v[i])
		xsum[i] = xsum[i-1] ^ int64(v[i])
	}
	nxt := make([]int, n+2)
	id := n + 1
	for i := n; i >= 1; i-- {
		nxt[i] = id
		if v[i] != 0 {
			id = i
		}
	}
	for q := 0; q < Q; q++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		getval := func(i, j int) int64 {
			return (psum[j] - psum[i-1]) - (xsum[j] ^ xsum[i-1])
		}
		mxval := getval(a, b)
		s := a
		if s <= n && v[s] == 0 {
			s = nxt[s]
		}
		ansA, ansB := a, b
		for i := 0; i < 31; i++ {
			if s > b {
				s = b
				i = 31
			}
			if getval(s, b) != mxval {
				break
			}
			l, r := s-1, b
			for l+1 < r {
				m := (l + r) >> 1
				if getval(s, m) == mxval {
					r = m
				} else {
					l = m
				}
			}
			if (ansB - ansA) > (r - s) {
				ansA = s
				ansB = r
			}
			if s <= n {
				s = nxt[s]
			}
		}
		fmt.Fprint(writer, ansA, " ", ansB, "\n")
	}
}

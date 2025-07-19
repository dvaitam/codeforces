package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader = bufio.NewReader(os.Stdin)

func readInt() int {
	sign := 1
	c, err := reader.ReadByte()
	for err == nil && (c < '0' || c > '9') && c != '-' {
		c, err = reader.ReadByte()
	}
	if err != nil {
		return 0
	}
	if c == '-' {
		sign = -1
		c, _ = reader.ReadByte()
	}
	val := 0
	for err == nil && c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = reader.ReadByte()
	}
	return sign * val
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	n := readInt()
	c := readInt()
	a := make([]int, n+1)
	s := make([]int, n+1)
	maxVal := 0
	for i := 1; i <= n; i++ {
		a[i] = readInt()
		if a[i] > maxVal {
			maxVal = a[i]
		}
		s[i] = s[i-1]
		if a[i] == c {
			s[i]++
		}
	}
	f := make([]int, n+1)
	l := make([]int, maxVal+1)
	for i := 1; i <= n; i++ {
		v := a[i]
		prev := l[v]
		t1 := f[prev] + 1
		t2 := s[i-1] + 1
		cmx := max(t1, t2)
		if prev > 0 {
			cand := s[prev-1] + 1
			f[i] = max(cmx, cand)
		} else {
			f[i] = cmx
		}
		l[v] = i
	}
	ans := 0
	nw := 0
	for i := n; i >= 1; i-- {
		if f[i]+nw > ans {
			ans = f[i] + nw
		}
		if a[i] == c {
			nw++
		}
	}
	// output answer
	fmt.Println(ans)
}

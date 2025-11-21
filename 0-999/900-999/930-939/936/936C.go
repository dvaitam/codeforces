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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var sStr, tStr string
	fmt.Fscan(in, &sStr, &tStr)

	if len(sStr) != n || len(tStr) != n {
		fmt.Fprintln(out, -1)
		return
	}

	s := []byte(sStr)
	t := []byte(tStr)

	if !sameMultiset(s, t) {
		fmt.Fprintln(out, -1)
		return
	}

	ans := make([]int, 0, 3*n)
	for i := 0; i < n; i++ {
		target := t[n-1-i]
		j := i
		for j < n && s[j] != target {
			j++
		}
		if j == n {
			fmt.Fprintln(out, -1)
			return
		}
		ans = append(ans, n, j, 1)
		reverseSuffix(s, j)
		moveLastToFront(s)
	}

	fmt.Fprintln(out, len(ans))
	for i, v := range ans {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, v)
	}
	out.WriteByte('\n')
}

func sameMultiset(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var cnt [26]int
	for _, ch := range a {
		cnt[ch-'a']++
	}
	for _, ch := range b {
		cnt[ch-'a']--
	}
	for _, v := range cnt {
		if v != 0 {
			return false
		}
	}
	return true
}

func reverseSuffix(s []byte, start int) {
	for l, r := start, len(s)-1; l < r; l, r = l+1, r-1 {
		s[l], s[r] = s[r], s[l]
	}
}

func moveLastToFront(s []byte) {
	if len(s) <= 1 {
		return
	}
	last := s[len(s)-1]
	copy(s[1:], s[:len(s)-1])
	s[0] = last
}

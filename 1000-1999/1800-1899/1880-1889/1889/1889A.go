package main

import (
	"bufio"
	"fmt"
	"os"
)

func isGood(s string) bool {
	n := len(s)
	if n%2 == 1 {
		return false
	}
	for i := 0; i < n/2; i++ {
		if s[i] == s[n-1-i] {
			return false
		}
	}
	return true
}

func solveCase(s string) (int, []int) {
	count0, count1 := 0, 0
	for _, ch := range s {
		if ch == '0' {
			count0++
		} else {
			count1++
		}
	}
	if count0 != count1 {
		return -1, nil
	}
	if isGood(s) {
		return 0, nil
	}
	n := len(s)
	// try one insertion
	for i := 0; i <= n; i++ {
		t := s[:i] + "01" + s[i:]
		if isGood(t) {
			return 1, []int{i}
		}
	}
	// try two insertions
	for i := 0; i <= n; i++ {
		t1 := s[:i] + "01" + s[i:]
		m := len(t1)
		for j := 0; j <= m; j++ {
			t2 := t1[:j] + "01" + t1[j:]
			if isGood(t2) {
				return 2, []int{i, j}
			}
		}
	}
	return -1, nil
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n, &s)
		opsCount, ops := solveCase(s)
		if opsCount == -1 {
			fmt.Fprintln(out, -1)
			continue
		}
		fmt.Fprintln(out, opsCount)
		if opsCount > 0 {
			for i, v := range ops {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			out.WriteByte('\n')
		} else {
			out.WriteByte('\n')
		}
	}
}

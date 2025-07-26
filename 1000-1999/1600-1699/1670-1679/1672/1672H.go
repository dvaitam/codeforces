package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	var a string
	fmt.Fscan(reader, &a)

	s := make([]byte, n+2)
	s[0] = a[0]
	for i := 0; i < n; i++ {
		s[i+1] = a[i]
	}
	s[n+1] = a[n-1]

	l := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if s[i-1] == s[i] {
			l[i] = i
		} else {
			l[i] = l[i-1]
		}
	}

	r := make([]int, n+2)
	psum := make([]int, n+2)
	balance := make([]int, n+2)
	for i := n; i >= 1; i-- {
		if s[i] == s[i+1] {
			r[i] = i
			psum[i] = 1
			if (i-l[i])%2 == 0 {
				if s[i] == '1' {
					balance[i] = 1
				} else {
					balance[i] = -1
				}
			}
		} else {
			r[i] = r[i+1]
		}
	}
	for i := 1; i <= n; i++ {
		psum[i] += psum[i-1]
		balance[i] += balance[i-1]
	}

	for ; q > 0; q-- {
		var L, R int
		fmt.Fscan(reader, &L, &R)
		aIdx := l[L]
		bIdx := r[R]
		bl := balance[bIdx] - balance[aIdx-1]
		sum := psum[bIdx] - psum[aIdx-1]
		ans := (sum + abs(bl)) / 2
		if (sum+abs(bl))%2 == 1 {
			ans++
		} else if abs(bl) == 0 {
			ans++
		} else if (bl > 0) != (s[aIdx] == '1') {
			ans++
		}
		fmt.Fprintln(writer, ans)
	}
}

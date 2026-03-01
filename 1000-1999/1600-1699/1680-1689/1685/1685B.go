package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func solve(a, b, c, d int, s string) bool {
	n := len(s)
	countA := 0
	for i := 0; i < n; i++ {
		if s[i] == 'A' {
			countA++
		}
	}
	countB := n - countA
	if countA != a+c+d || countB != b+c+d {
		return false
	}

	tot := 0
	l := 0
	var A, B []int
	for i := 0; i < n; i++ {
		if i == n-1 || s[i] == s[i+1] {
			length := i - l + 1
			if length%2 == 1 {
				tot += length / 2
			} else {
				x := length / 2
				if s[l] == 'A' {
					A = append(A, x)
				} else {
					B = append(B, x)
				}
			}
			l = i + 1
		}
	}

	sort.Ints(A)
	sort.Ints(B)

	for _, x := range A {
		if c >= x {
			c -= x
		} else {
			tot += x - 1
		}
	}
	for _, x := range B {
		if d >= x {
			d -= x
		} else {
			tot += x - 1
		}
	}

	return c+d <= tot
}

func main() {
	in := bufio.NewReader(os.Stdin)
	t := 0
	fmt.Fscan(in, &t)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var a, b, c, d int
		var s string
		fmt.Fscan(in, &a, &b, &c, &d)
		fmt.Fscan(in, &s)
		if solve(a, b, c, d, s) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

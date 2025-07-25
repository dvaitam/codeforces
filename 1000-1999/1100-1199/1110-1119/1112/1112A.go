package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	s := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &s[i])
	}
	c := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &c[i])
	}
	// compute max power per school
	maxp := make([]int, m+1)
	for i := 1; i <= n; i++ {
		si := s[i]
		if p[i] > maxp[si] {
			maxp[si] = p[i]
		}
	}
	// count chosen ones that are not strongest in their school
	ans := 0
	for _, ci := range c {
		if p[ci] < maxp[s[ci]] {
			ans++
		}
	}
	fmt.Println(ans)
}

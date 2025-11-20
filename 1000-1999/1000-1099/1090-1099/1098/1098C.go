package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	n int
	s int64
	a []int
	b []int
)

func chk(x int) bool {
	for i := 0; i <= n+1; i++ {
		b[i] = 0
	}
	a[1] = 1
	i := 2
	depth := 1
	t := int64(1)
	rem := s - 1
	for i <= n {
		t *= int64(x)
		if t > int64(n) {
			t = int64(n)
		}
		depth++
		b[depth] = 0
		limit := int(t)
		for j := 0; j < limit && i <= n; j++ {
			a[i] = depth
			rem -= int64(depth)
			b[depth]++
			i++
			if rem < 0 {
				return false
			}
		}
	}
	if rem < 0 {
		return false
	}
	j := n
	for rem > 0 {
		depth++
		for j > 1 && b[a[j]] == 1 {
			j--
		}
		if j <= 1 {
			return false
		}
		inc := depth - a[j]
		if int64(inc) > rem {
			inc = int(rem)
		}
		rem -= int64(inc)
		b[a[j]]--
		a[j] += inc
		j--
	}
	return rem == 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}
	minSum := int64(2*n - 1)
	maxSum := int64(n) * int64(n+1) / 2
	if s < minSum || s > maxSum {
		fmt.Fprintln(writer, "No")
		return
	}

	a = make([]int, n+2)
	b = make([]int, n+2)

	l, r := 1, n
	ans := n
	for l <= r {
		mid := (l + r) / 2
		if chk(mid) {
			ans = mid
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	chk(ans)

	a[1] = 1
	depths := make([]int, n-1)
	for i := 2; i <= n; i++ {
		depths[i-2] = a[i]
	}
	sort.Ints(depths)
	for i := 2; i <= n; i++ {
		a[i] = depths[i-2]
	}

	children := make([]int, n+2)
	parents := make([]int, n+1)
	ptr := 1
	for i := 2; i <= n; i++ {
		for a[ptr] != a[i]-1 || children[ptr] == ans {
			ptr++
		}
		parents[i] = ptr
		children[ptr]++
	}

	fmt.Fprintln(writer, "Yes")
	for i := 2; i <= n; i++ {
		if i > 2 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, parents[i])
	}
	fmt.Fprintln(writer)
}

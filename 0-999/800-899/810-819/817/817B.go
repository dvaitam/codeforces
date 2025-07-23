package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func choose(n int64, k int) int64 {
	if k == 0 {
		return 1
	}
	switch k {
	case 1:
		return n
	case 2:
		return n * (n - 1) / 2
	case 3:
		return n * (n - 1) * (n - 2) / 6
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sort.Ints(a)

	v1 := a[0]
	cnt1 := 0
	for cnt1 < n && a[cnt1] == v1 {
		cnt1++
	}
	if cnt1 >= 3 {
		fmt.Fprintln(out, choose(int64(cnt1), 3))
		return
	}

	v2 := a[cnt1]
	cnt2 := 0
	idx := cnt1
	for idx < n && a[idx] == v2 {
		cnt2++
		idx++
	}
	if cnt1 == 2 {
		fmt.Fprintln(out, choose(2, 2)*int64(cnt2))
		return
	}
	// cnt1 == 1
	if cnt2 >= 2 {
		fmt.Fprintln(out, choose(int64(cnt2), 2))
		return
	}
	// cnt1 == 1, cnt2 == 1
	cnt3 := 0
	v3 := a[idx]
	for idx < n && a[idx] == v3 {
		cnt3++
		idx++
	}
	fmt.Fprintln(out, int64(cnt3))
}

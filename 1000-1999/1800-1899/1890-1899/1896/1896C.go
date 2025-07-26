package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, x int
		fmt.Fscan(in, &n, &x)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		ok, ans := solveCase(n, x, a, b)
		if !ok {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			for i := 0; i < n; i++ {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, ans[i])
			}
			fmt.Fprintln(out)
		}
	}
}

func solveCase(n, x int, a, b []int) (bool, []int) {
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return a[idx[i]] < a[idx[j]] })
	A := make([]int, n)
	for i, pos := range idx {
		A[i] = a[pos]
	}
	B := append([]int(nil), b...)
	sort.Ints(B)

	diff := make([]int, n+1)
	for i, ai := range A {
		hi := sort.Search(len(B), func(j int) bool { return B[j] >= ai }) - 1
		if hi >= 0 {
			l := (n - i) % n
			r := (hi - i + n) % n
			if l <= r {
				diff[l]++
				diff[r+1]--
			} else {
				diff[0]++
				diff[r+1]--
				diff[l]++
			}
		}
	}
	beauties := make([]int, n)
	cur := 0
	for s := 0; s < n; s++ {
		cur += diff[s]
		beauties[s] = cur
	}
	shift := -1
	for s := 0; s < n; s++ {
		if beauties[s] == x {
			shift = s
			break
		}
	}
	if shift == -1 {
		return false, nil
	}
	ans := make([]int, n)
	for i, pos := range idx {
		ans[pos] = B[(i+shift)%n]
	}
	return true, ans
}

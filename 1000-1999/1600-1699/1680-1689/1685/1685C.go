package main

import (
	"bufio"
	"fmt"
	"os"
)

func isBalancedBytes(s []byte) bool {
	bal := 0
	for _, ch := range s {
		if ch == '(' {
			bal++
		} else {
			bal--
		}
		if bal < 0 {
			return false
		}
	}
	return bal == 0
}

func solve(s string) [][2]int {
	n := len(s) / 2
	pref := make([]int, 2*n+1)
	good := true
	for i := 0; i < 2*n; i++ {
		if s[i] == '(' {
			pref[i+1] = pref[i] + 1
		} else {
			pref[i+1] = pref[i] - 1
		}
		if pref[i+1] < 0 {
			good = false
		}
	}
	if good {
		return nil
	}
	first := 0
	for first <= 2*n && pref[first] >= 0 {
		first++
	}
	last := 2 * n
	for last >= 0 && pref[last] >= 0 {
		last--
	}
	left := 0
	for i := 0; i < first; i++ {
		if pref[i] > pref[left] {
			left = i
		}
	}
	right := last
	for i := last; i <= 2*n; i++ {
		if pref[i] > pref[right] {
			right = i
		}
	}
	arr := []byte(s)
	for l, r := left, right-1; l < r; l, r = l+1, r-1 {
		arr[l], arr[r] = arr[r], arr[l]
	}
	if isBalancedBytes(arr) {
		return [][2]int{{left + 1, right}}
	}
	maxIdx := 0
	for i := 0; i <= 2*n; i++ {
		if pref[i] > pref[maxIdx] {
			maxIdx = i
		}
	}
	return [][2]int{{1, maxIdx}, {maxIdx + 1, 2 * n}}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		ops := solve(s)
		fmt.Fprintln(out, len(ops))
		for _, p := range ops {
			fmt.Fprintln(out, p[0], p[1])
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func nextPerm(a []int) bool {
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] > a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for a[j] < a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func checkInequalities(p []int, s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			if !(p[i] < p[i+1]) {
				return false
			}
		} else {
			if !(p[i] > p[i+1]) {
				return false
			}
		}
	}
	return true
}

func isGood(p []int) bool {
	n := len(p)
	for l := 0; l < n; l++ {
		for r := l + 2; r < n; r++ {
			// check subarray [l,r]
			arr := make([]int, r-l+1)
			copy(arr, p[l:r+1])
			// find first, second, third minimums
			// simple selection for small arrays
			for i := 0; i < 3; i++ {
				minIdx := i
				for j := i + 1; j < len(arr); j++ {
					if arr[j] < arr[minIdx] {
						minIdx = j
					}
				}
				arr[i], arr[minIdx] = arr[minIdx], arr[i]
			}
			first, second, third := arr[0], arr[1], arr[2]
			if (p[l] == first && p[r] == second) || (p[l] == second && p[r] == first) {
				if third != p[l+1] && third != p[r-1] {
					return false
				}
			}
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	var s string
	fmt.Fscan(in, &s)

	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	count := 0
	for {
		if checkInequalities(p, s) && isGood(p) {
			count = (count + 1) % MOD
		}
		if !nextPerm(p) {
			break
		}
	}
	fmt.Println(count)
}

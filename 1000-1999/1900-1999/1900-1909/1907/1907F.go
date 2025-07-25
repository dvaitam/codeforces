package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// kmpAll finds all occurrences of pattern in text and returns their starting indices.
func kmpAll(text, pattern []int) []int {
	m := len(pattern)
	if m == 0 {
		return []int{}
	}
	pi := make([]int, m)
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && pattern[i] != pattern[j] {
			j = pi[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		pi[i] = j
	}
	res := make([]int, 0)
	j := 0
	for i := 0; i < len(text); i++ {
		for j > 0 && text[i] != pattern[j] {
			j = pi[j-1]
		}
		if text[i] == pattern[j] {
			j++
		}
		if j == m {
			res = append(res, i-m+1)
			j = pi[j-1]
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, n)
		copy(b, a)
		sort.Ints(b)
		desc := make([]int, n)
		for i := 0; i < n; i++ {
			desc[i] = b[n-1-i]
		}
		asc2 := append(append([]int{}, b...), b...)
		desc2 := append(append([]int{}, desc...), desc...)

		const INF int = int(1e9)
		ans := INF
		// check rotations of ascending order
		for _, pos := range kmpAll(asc2, a) {
			if pos >= n {
				pos -= n
			}
			r := pos
			cost := r
			if n-r+2 < cost {
				cost = n - r + 2
			}
			if cost < ans {
				ans = cost
			}
		}
		// check rotations of descending order
		for _, pos := range kmpAll(desc2, a) {
			if pos >= n {
				pos -= n
			}
			r := pos
			cost := 1 + r
			if 1+n-r < cost {
				cost = 1 + n - r
			}
			if cost < ans {
				ans = cost
			}
		}
		if ans == INF {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, ans)
		}
	}
}

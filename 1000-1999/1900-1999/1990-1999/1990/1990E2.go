package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	N = 5007
	B = 70
)

func query(x int) int {
	fmt.Printf("? %d\n", x)
	os.Stdout.Sync()
	var b int
	fmt.Scan(&b)
	return b
}

func dfs1(u, fa int, to [][]int, h []int) {
	h[u] = 1
	for _, v := range to[u] {
		if v != fa {
			dfs1(v, u, to, h)
			if h[v]+1 > h[u] {
				h[u] = h[v] + 1
			}
		}
	}
}

func search(u, fa int, to [][]int, h []int, a []int, k *int) {
	*k++
	a[*k] = u
	vec := []int{}
	for _, v := range to[u] {
		if v != fa && h[v] > B {
			vec = append(vec, v)
		}
	}
	if len(vec) == 0 {
		return
	}
	for len(vec) > 1 {
		if query(vec[len(vec)-1]) != 0 {
			search(vec[len(vec)-1], u, to, h, a, k)
			return
		}
		vec = vec[:len(vec)-1]
	}
	search(vec[len(vec)-1], u, to, h, a, k)
}

func solve(n int, to [][]int, h []int, a []int) int {
	var k int
	dfs1(1, 0, to, h)

	for i := 1; i <= n; i++ {
		if h[i] == 1 {
			for c := 1; c <= B; c++ {
				if query(i) != 0 {
					return i
				}
			}
			break
		}
	}

	search(1, 0, to, h, a, &k)

	l := 1
	r := k
	for l < r {
		mid := (l + r + 1) >> 1
		if query(a[mid]) != 0 {
			l = mid
		} else {
			if l-1 > 1 {
				l = l - 1
			} else {
				l = 1
			}
			if mid-2 > 1 {
				r = mid - 2
			} else {
				r = 1
			}
		}
	}
	return a[l]
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	var T int
	fmt.Scan(&T)

	for T > 0 {
		scanner.Scan()
		var n int
		fmt.Sscan(scanner.Text(), &n)

		to := make([][]int, n+1)
		h := make([]int, n+1)
		a := make([]int, n+1)

		for i := 1; i < n; i++ {
			scanner.Scan()
			var u int
			fmt.Sscan(scanner.Text(), &u)
			scanner.Scan()
			var v int
			fmt.Sscan(scanner.Text(), &v)
			to[u] = append(to[u], v)
			to[v] = append(to[v], u)
		}

		x := solve(n, to, h, a)
		fmt.Printf("! %d\n", x)
		os.Stdout.Sync()

		T--
	}
}

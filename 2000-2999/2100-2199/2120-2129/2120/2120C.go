package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var m int64
		fmt.Fscan(in, &n, &m)

		root := -1
		var minSum, maxSum int64
		for r := 1; r <= n; r++ {
			minSum = int64(n - 1 + r)
			maxSum = int64(r*(r+1)/2 + (n-r)*r)
			if m >= minSum && m <= maxSum {
				root = r
				break
			}
		}
		if root == -1 {
			fmt.Fprintln(out, -1)
			continue
		}

		l := root
		var base int64
		for l > 1 {
			chainLen := root - (l - 1) + 1
			sumChain := int64((root + (l - 1)) * chainLen / 2)
			nextBase := sumChain + int64(n-chainLen)
			if nextBase <= m {
				l--
			} else {
				break
			}
		}
		chainLen := root - l + 1
		sumChain := int64((root + l) * chainLen / 2)
		base = sumChain + int64(n-chainLen)
		delta := m - base

		parent := make([]int, n+1)
		for i := root; i > l; i-- {
			parent[i-1] = i
		}

		pool := make([]int, 0, n-chainLen)
		for i := 1; i < l; i++ {
			pool = append(pool, i)
		}
		for i := root + 1; i <= n; i++ {
			pool = append(pool, i)
		}

		pointer := 0
		remaining := len(pool)
		for v := root; v >= l; v-- {
			if delta == 0 {
				break
			}
			if v == 1 {
				continue
			}
			inc := int64(v - 1)
			if inc == 0 || remaining == 0 {
				continue
			}
			take := delta / inc
			if take > int64(remaining) {
				take = int64(remaining)
			}
			for ; take > 0; take-- {
				node := pool[pointer]
				pointer++
				parent[node] = v
				delta -= inc
				remaining--
				if delta == 0 {
					break
				}
			}
		}

		for pointer < len(pool) {
			node := pool[pointer]
			pointer++
			if node == 1 && l > 1 && parent[node] == 0 {
				parent[node] = l
			} else {
				parent[node] = 1
			}
		}

		if parent[1] == 0 && root != 1 {
			parent[1] = l
		}

		fmt.Fprintln(out, root)
		for i := 1; i <= n; i++ {
			if i == root {
				continue
			}
			if parent[i] == 0 {
				parent[i] = root
			}
			fmt.Fprintln(out, i, parent[i])
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	words := make([][]int, n)
	for i := 0; i < n; i++ {
		var l int
		fmt.Fscan(reader, &l)
		w := make([]int, l)
		for j := 0; j < l; j++ {
			fmt.Fscan(reader, &w[j])
		}
		words[i] = w
	}

	adj := make([][]int, m+1)
	mustUpper := make([]bool, m+1)
	forbidUpper := make([]bool, m+1)

	for i := 0; i < n-1; i++ {
		a := words[i]
		b := words[i+1]
		j := 0
		for j < len(a) && j < len(b) && a[j] == b[j] {
			j++
		}
		if j == len(a) || j == len(b) {
			if len(a) > len(b) {
				fmt.Fprintln(writer, "No")
				return
			}
			continue
		}
		x := a[j]
		y := b[j]
		if x < y {
			adj[y] = append(adj[y], x)
		} else if x > y {
			mustUpper[x] = true
			forbidUpper[y] = true
		}
	}

	// propagate mustUpper through edges
	queue := make([]int, 0)
	upper := make([]bool, m+1)
	for i := 1; i <= m; i++ {
		if mustUpper[i] {
			upper[i] = true
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, u := range adj[v] {
			if !upper[u] {
				upper[u] = true
				queue = append(queue, u)
			}
		}
	}

	// check conflicts with forbidUpper
	for i := 1; i <= m; i++ {
		if upper[i] && forbidUpper[i] {
			fmt.Fprintln(writer, "No")
			return
		}
	}

	// verify order after applying capitalization
	cmpLetter := func(x, y int) int {
		if upper[x] && upper[y] {
			if x < y {
				return -1
			} else if x > y {
				return 1
			}
			return 0
		}
		if upper[x] && !upper[y] {
			return -1
		}
		if !upper[x] && upper[y] {
			return 1
		}
		if x < y {
			return -1
		} else if x > y {
			return 1
		}
		return 0
	}

	for i := 0; i < n-1; i++ {
		a := words[i]
		b := words[i+1]
		len1 := len(a)
		len2 := len(b)
		minLen := len1
		if len2 < minLen {
			minLen = len2
		}
		k := 0
		for k < minLen {
			if a[k] != b[k] {
				cmp := cmpLetter(a[k], b[k])
				if cmp > 0 {
					fmt.Fprintln(writer, "No")
					return
				}
				break
			}
			k++
		}
		if k == minLen {
			if len1 > len2 {
				fmt.Fprintln(writer, "No")
				return
			}
		}
	}

	fmt.Fprintln(writer, "Yes")
	var letters []int
	for i := 1; i <= m; i++ {
		if upper[i] {
			letters = append(letters, i)
		}
	}
	fmt.Fprintln(writer, len(letters))
	for i, v := range letters {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	if len(letters) > 0 {
		fmt.Fprintln(writer)
	}
}

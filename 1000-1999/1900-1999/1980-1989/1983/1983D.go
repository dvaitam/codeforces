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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		pos := make(map[int]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
			pos[b[i]] = i
		}
		perm := make([]int, n)
		ok := true
		for i := 0; i < n; i++ {
			p, found := pos[a[i]]
			if !found {
				ok = false
				break
			}
			perm[i] = p
		}
		if !ok {
			fmt.Fprintln(out, "NO")
			continue
		}
		visited := make([]bool, n)
		cycles := 0
		for i := 0; i < n; i++ {
			if !visited[i] {
				cycles++
				j := i
				for !visited[j] {
					visited[j] = true
					j = perm[j]
				}
			}
		}
		parity := (n - cycles) % 2
		if parity == 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

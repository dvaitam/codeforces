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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a := make([]int, n)
	b := make([]int, n)

	for i := 0; i < n; i++ {
		fmt.Fprintf(writer, "? %d %d\n", i, 0)
		writer.Flush()
		if _, err := fmt.Fscan(reader, &a[i]); err != nil {
			return
		}
	}

	for j := 0; j < n; j++ {
		fmt.Fprintf(writer, "? %d %d\n", 0, j)
		writer.Flush()
		if _, err := fmt.Fscan(reader, &b[j]); err != nil {
			return
		}
	}

	count := 0
	var ans []int

	for pos0 := 0; pos0 < n; pos0++ {
		perm := make([]int, n)
		used := make([]bool, n)
		ok := true
		for i := 0; i < n; i++ {
			v := a[i] ^ pos0
			if v < 0 || v >= n || used[v] {
				ok = false
				break
			}
			perm[i] = v
			used[v] = true
		}
		if !ok {
			continue
		}
		p0 := perm[0]
		inv := make([]int, n)
		for j := 0; j < n; j++ {
			inv[j] = b[j] ^ p0
		}
		for i := 0; i < n && ok; i++ {
			if inv[perm[i]] != i {
				ok = false
			}
		}
		if ok {
			count++
			if ans == nil {
				ans = perm
			}
		}
	}

	fmt.Fprintln(writer, "!")
	fmt.Fprintln(writer, count)
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}

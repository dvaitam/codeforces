package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	mat := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		mat[i] = []byte(s)
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	used := make([]bool, n)
	for {
		idx := -1
		for i := 0; i < n; i++ {
			if b[i] == 0 {
				idx = i
				break
			}
		}
		if idx == -1 {
			break
		}
		used[idx] = true
		for j := 0; j < n; j++ {
			if mat[idx][j] == '1' {
				b[j]--
			}
		}
	}

	res := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if used[i] {
			res = append(res, i+1)
		}
	}
	sort.Ints(res)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	fmt.Fprintln(w, len(res))
	if len(res) == 0 {
		fmt.Fprintln(w)
		return
	}
	for i, v := range res {
		if i > 0 {
			w.WriteByte(' ')
		}
		fmt.Fprint(w, v)
	}
	w.WriteByte('\n')
}

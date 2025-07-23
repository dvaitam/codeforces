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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	pos := make(map[int64]int)
	removed := make([]bool, n)

	for i := 0; i < n; i++ {
		v := a[i]
		for {
			if j, ok := pos[v]; ok {
				delete(pos, v)
				removed[j] = true
				v *= 2
			} else {
				break
			}
		}
		a[i] = v
		pos[v] = i
	}

	cnt := 0
	for i := 0; i < n; i++ {
		if !removed[i] {
			cnt++
		}
	}
	fmt.Fprintln(out, cnt)
	first := true
	for i := 0; i < n; i++ {
		if !removed[i] {
			if !first {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, a[i])
			first = false
		}
	}
	if cnt > 0 {
		fmt.Fprintln(out)
	}
}

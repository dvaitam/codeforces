package main

import (
	"bufio"
	"fmt"
	"os"
)

func better(a, b [5]int) bool {
	cnt := 0
	for i := 0; i < 5; i++ {
		if a[i] < b[i] {
			cnt++
		}
	}
	return cnt >= 3
}

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
		athletes := make([][5]int, n)
		for i := 0; i < n; i++ {
			for j := 0; j < 5; j++ {
				fmt.Fscan(in, &athletes[i][j])
			}
		}
		candidate := 0
		for i := 1; i < n; i++ {
			if better(athletes[i], athletes[candidate]) {
				candidate = i
			}
		}
		ok := true
		for i := 0; i < n; i++ {
			if i == candidate {
				continue
			}
			if !better(athletes[candidate], athletes[i]) {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, candidate+1)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}

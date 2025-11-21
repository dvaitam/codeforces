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
		fmt.Fscan(in, &n)

		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			var val int
			fmt.Fscan(in, &val)
			pos[val] = i
		}

		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = pos[a[i]]
		}

		inc := true
		dec := true
		if order[0] != 0 {
			inc = false
		}
		if order[0] != n-1 {
			dec = false
		}
		for i := 1; i < n; i++ {
			if inc && order[i] != order[i-1]+1 {
				inc = false
			}
			if dec && order[i] != order[i-1]-1 {
				dec = false
			}
			if !inc && !dec {
				break
			}
		}

		if inc || dec {
			fmt.Fprintln(out, "Bob")
		} else {
			fmt.Fprintln(out, "Alice")
		}
	}
}

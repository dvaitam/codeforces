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
	x := make([]int, n)
	y := make([]int, n)
	present := make(map[int]struct{}, 2*n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i])
		present[x[i]] = struct{}{}
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &y[i])
		present[y[i]] = struct{}{}
	}
	count := 0
	for _, a := range x {
		for _, b := range y {
			if _, ok := present[a^b]; ok {
				count++
			}
		}
	}
	if count%2 == 0 {
		fmt.Fprintln(out, "Karen")
	} else {
		fmt.Fprintln(out, "Koyomi")
	}
}

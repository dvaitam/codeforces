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

	var n, d int
	fmt.Fscan(in, &n, &d)
	cities := make(map[int]struct{}, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		cities[x] = struct{}{}
	}

	count := 0
	for x := range cities {
		if _, ok := cities[x+d]; ok {
			count++
		}
	}
	fmt.Fprintln(out, count)
}

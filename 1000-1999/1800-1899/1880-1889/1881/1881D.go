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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		counts := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			v := x
			for p := 2; p*p <= v; p++ {
				for v%p == 0 {
					counts[p]++
					v /= p
				}
			}
			if v > 1 {
				counts[v]++
			}
		}
		ok := true
		for _, c := range counts {
			if c%n != 0 {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

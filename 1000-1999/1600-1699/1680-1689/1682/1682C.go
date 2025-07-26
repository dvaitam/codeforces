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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make(map[int]int, n)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
		}
		singles, dup := 0, 0
		for _, c := range freq {
			if c == 1 {
				singles++
			} else {
				dup++
			}
		}
		ans := dup + (singles+1)/2
		fmt.Fprintln(writer, ans)
	}
}

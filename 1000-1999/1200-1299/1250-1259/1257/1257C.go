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
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		last := make(map[int]int)
		best := n + 10
		for i, v := range arr {
			if prev, ok := last[v]; ok {
				if i-prev+1 < best {
					best = i - prev + 1
				}
			}
			last[v] = i
		}
		if best > n {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, best)
		}
	}
}

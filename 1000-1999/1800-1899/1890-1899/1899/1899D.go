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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
		}
		var ans int64
		for _, c := range freq {
			ans += int64(c*(c-1)) / 2
		}
		if c1, ok1 := freq[1]; ok1 {
			if c2, ok2 := freq[2]; ok2 {
				ans += int64(c1 * c2)
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

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
		var a, b int64
		fmt.Fscan(reader, &a, &b)
		if a < b {
			a, b = b, a
		}
		diff := a - b
		if diff == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		var k int64 = 0
		var sum int64 = 0
		for {
			k++
			sum += k
			if sum >= diff && (sum-diff)%2 == 0 {
				fmt.Fprintln(writer, k)
				break
			}
		}
	}
}

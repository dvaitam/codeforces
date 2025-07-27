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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		p := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		sum := p[0]
		var add int64
		for i := 1; i < n; i++ {
			required := (p[i]*100 + int64(k) - 1) / int64(k)
			if sum < required {
				diff := required - sum
				add += diff
				sum += diff
			}
			sum += p[i]
		}
		fmt.Fprintln(writer, add)
	}
}

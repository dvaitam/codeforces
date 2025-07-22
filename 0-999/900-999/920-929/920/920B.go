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
		l := make([]int, n)
		r := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &l[i], &r[i])
		}
		cur := 1
		for i := 0; i < n; i++ {
			if cur < l[i] {
				cur = l[i]
			}
			if cur > r[i] {
				fmt.Fprint(writer, 0)
			} else {
				fmt.Fprint(writer, cur)
				cur++
			}
			if i+1 < n {
				writer.WriteByte(' ')
			}
		}
		writer.WriteByte('\n')
	}
}

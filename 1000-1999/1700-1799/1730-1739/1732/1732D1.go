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

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	set := make(map[int64]struct{})
	set[0] = struct{}{}
	next := make(map[int64]int64)

	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(reader, &op)
		if op == "+" {
			var x int64
			fmt.Fscan(reader, &x)
			set[x] = struct{}{}
		} else if op == "?" {
			var k int64
			fmt.Fscan(reader, &k)
			v := next[k]
			for {
				if _, found := set[v]; !found {
					fmt.Fprintln(writer, v)
					next[k] = v
					break
				}
				v += k
			}
		}
	}
}

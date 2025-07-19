package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(reader, &l, &r)
		found := false
		for i := int64(2); i*i <= r; i++ {
			k := (r / i) * i
			if k < l {
				continue
			}
			if min(i, k-i) != 1 {
				fmt.Fprintln(writer, i, k-i)
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintln(writer, -1)
		}
	}
}

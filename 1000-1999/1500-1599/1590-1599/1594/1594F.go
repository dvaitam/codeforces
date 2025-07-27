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
		var s, n, k int64
		fmt.Fscan(reader, &s, &n, &k)

		if k > s {
			fmt.Fprintln(writer, "NO")
			continue
		}
		if k == s {
			fmt.Fprintln(writer, "YES")
			continue
		}

		q := s / k
		rem := s % k
		var total int64
		if k > s {
			total = s + 1
		} else {
			highCnt := rem + 1
			lowCnt := k - highCnt
			total = highCnt*((q+2)/2) + lowCnt*((q+1)/2)
		}

		if n >= total {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

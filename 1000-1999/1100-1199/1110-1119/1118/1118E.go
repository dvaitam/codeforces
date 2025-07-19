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

	var n, k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	if n > k*(k-1) {
		fmt.Fprintln(writer, "NO")
		return
	}
	fmt.Fprintln(writer, "YES")
	var cnt, add, fi, se int64
	for cnt < n {
		if cnt%k == 0 {
			add++
		}
		cnt++
		fi = fi%k + 1
		se = (fi+add-1)%k + 1
		fmt.Fprintln(writer, fi, se)
	}
}

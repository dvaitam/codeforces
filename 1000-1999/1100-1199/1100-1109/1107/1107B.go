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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var k, x int64
		fmt.Fscan(reader, &k, &x)
		ans := x + (k-1)*9
		fmt.Fprintln(writer, ans)
	}
}

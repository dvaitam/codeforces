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

	var n, m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var ans int64
	if m == 1 {
		ans = n - 1
	} else {
		ans = n * (m - 1)
	}
	fmt.Fprintln(writer, ans)
}

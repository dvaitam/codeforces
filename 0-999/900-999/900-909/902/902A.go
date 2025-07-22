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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	maxReach := 0
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		if a <= maxReach && b > maxReach {
			maxReach = b
		}
	}

	if maxReach >= m {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}

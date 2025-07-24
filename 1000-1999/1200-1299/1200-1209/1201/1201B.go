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

	sum := int64(0)
	maxVal := int64(0)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		sum += x
		if x > maxVal {
			maxVal = x
		}
	}

	if sum%2 == 0 && maxVal <= sum/2 {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}

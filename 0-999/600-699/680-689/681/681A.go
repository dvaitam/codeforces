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
	good := false
	for i := 0; i < n; i++ {
		var name string
		var before, after int
		fmt.Fscan(reader, &name, &before, &after)
		if before >= 2400 && after > before {
			good = true
		}
	}
	if good {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}

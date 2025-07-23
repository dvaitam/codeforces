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

	var a, b, c int64
	if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
		return
	}

	if a == b {
		fmt.Fprintln(writer, "YES")
		return
	}
	if c == 0 {
		fmt.Fprintln(writer, "NO")
		return
	}
	if (b-a)%c == 0 && (b-a)/c >= 0 {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}

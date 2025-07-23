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

	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	if n%2 == 1 {
		fmt.Fprintln(writer, 0)
		return
	}
	half := n / 2
	ans := (half - 1) / 2
	fmt.Fprintln(writer, ans)
}

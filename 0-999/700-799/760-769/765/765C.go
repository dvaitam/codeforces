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

	var k, a, b int64
	if _, err := fmt.Fscan(reader, &k, &a, &b); err != nil {
		return
	}
	if (a < k && b < k) || (a%k != 0 && b%k != 0) || (a%k != 0 && b < k) || (b%k != 0 && a < k) {
		fmt.Fprintln(writer, -1)
		return
	}
	fmt.Fprintln(writer, a/k+b/k)
}

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

	q := n / 3
	r := n % 3
	ans := 2 * q
	if r > 0 {
		ans++
	}

	fmt.Fprintln(writer, ans)
}

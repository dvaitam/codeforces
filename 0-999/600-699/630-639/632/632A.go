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
	var p int64
	if _, err := fmt.Fscan(reader, &n, &p); err != nil {
		return
	}
	ops := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &ops[i])
	}

	var money, cur int64
	for i := n - 1; i >= 0; i-- {
		if ops[i] == "half" {
			money += cur * p
			cur *= 2
		} else { // halfplus
			money += cur*p + p/2
			cur = cur*2 + 1
		}
	}

	fmt.Fprintln(writer, money)
}

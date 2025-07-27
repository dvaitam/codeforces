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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		res := solve(s)
		fmt.Fprintln(writer, res)
	}
}

func solve(s string) int {
	bal := 0
	minBal := 0
	res := 0
	for i, c := range s {
		if c == '+' {
			bal++
		} else {
			bal--
		}
		if bal < minBal {
			res += i + 1
			minBal = bal
		}
	}
	res += len(s)
	return res
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

func solve(s string) string {
	n := len(s)
	digits := []byte(s)
	minRight := byte('9')
	for i := n - 1; i >= 0; i-- {
		orig := digits[i]
		if orig > minRight && orig < '9' {
			digits[i] = orig + 1
		}
		if orig < minRight {
			minRight = orig
		}
	}
	sort.Slice(digits, func(i, j int) bool { return digits[i] < digits[j] })
	return string(digits)
}

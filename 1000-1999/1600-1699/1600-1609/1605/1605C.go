package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)
		ans := solve(s)
		fmt.Fprintln(writer, ans)
	}
}

func solve(s string) int {
	if strings.Contains(s, "aa") {
		return 2
	}
	if strings.Contains(s, "aba") || strings.Contains(s, "aca") {
		return 3
	}
	if strings.Contains(s, "abca") || strings.Contains(s, "acba") {
		return 4
	}
	if strings.Contains(s, "abbacca") || strings.Contains(s, "accabba") {
		return 7
	}
	return -1
}

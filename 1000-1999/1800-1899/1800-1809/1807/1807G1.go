package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func canConstruct(c []int) bool {
	sort.Ints(c)
	if len(c) == 0 || c[0] != 1 {
		return false
	}
	sum := 1
	for i := 1; i < len(c); i++ {
		if c[i] > sum {
			return false
		}
		sum += c[i]
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &c[i])
		}
		if canConstruct(c) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

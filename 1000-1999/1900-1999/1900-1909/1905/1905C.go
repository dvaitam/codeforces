package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func minOperations(s string) int {
	t := []byte(s)
	sorted := make([]byte, len(t))
	copy(sorted, t)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	if string(sorted) == s {
		return 0
	}
	prefixMax := t[0]
	ops := 0
	for i := 1; i < len(t); i++ {
		if t[i] < prefixMax {
			ops++
		}
		if t[i] > prefixMax {
			prefixMax = t[i]
		}
	}
	return ops
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		fmt.Fprintln(writer, minOperations(s))
	}
}

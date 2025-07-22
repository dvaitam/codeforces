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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var s string
	fmt.Fscan(reader, &s)

	i := 0
	for i < n-1 {
		if s[i] == '1' {
			start := i
			for i < n-1 && s[i] == '1' {
				i++
			}
			end := i
			sort.Ints(a[start : end+1])
		} else {
			i++
		}
	}

	for i := 0; i < n; i++ {
		if a[i] != i+1 {
			fmt.Fprintln(writer, "NO")
			return
		}
	}
	fmt.Fprintln(writer, "YES")
}

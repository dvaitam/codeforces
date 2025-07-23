package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	names := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &names[i])
	}

	seen := make(map[string]bool, n)
	result := make([]string, 0, n)
	for i := n - 1; i >= 0; i-- {
		name := names[i]
		if !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for _, name := range result {
		fmt.Fprintln(writer, name)
	}
}

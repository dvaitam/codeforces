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

	numbers := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &numbers[i])
	}

	// Count in how many numbers each substring appears
	counts := make(map[string]int)
	for _, num := range numbers {
		seen := make(map[string]bool)
		for l := 1; l <= 9; l++ {
			for i := 0; i+l <= 9; i++ {
				sub := num[i : i+l]
				if !seen[sub] {
					counts[sub]++
					seen[sub] = true
				}
			}
		}
	}

	for _, num := range numbers {
		found := false
		for l := 1; l <= 9 && !found; l++ {
			for i := 0; i+l <= 9; i++ {
				sub := num[i : i+l]
				if counts[sub] == 1 {
					fmt.Fprintln(writer, sub)
					found = true
					break
				}
			}
		}
	}
}

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
	hashtags := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &hashtags[i])
	}

	for i := n - 2; i >= 0; i-- {
		a := hashtags[i]
		b := hashtags[i+1]
		la, lb := len(a), len(b)
		j := 1
		for j < la && j < lb && a[j] == b[j] {
			j++
		}
		if j == la {
			continue
		}
		if j == lb {
			hashtags[i] = a[:lb]
			continue
		}
		if a[j] > b[j] {
			hashtags[i] = a[:j]
		}
	}

	for _, s := range hashtags {
		fmt.Fprintln(writer, s)
	}
}

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
		freq := make([]int, 26)
		for _, ch := range s {
			freq[ch-'a']++
		}
		total := 0
		for _, f := range freq {
			if f >= 2 {
				total += 2
			} else if f == 1 {
				total += 1
			}
		}
		k := total / 2
		fmt.Fprintln(writer, k)
	}
}

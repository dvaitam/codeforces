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
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		ans := 0
		for i := 0; i < n; i++ {
			freq := [10]int{}
			distinct := 0
			maxFreq := 0
			for j := i; j < n && j < i+100; j++ {
				d := s[j] - '0'
				if freq[d] == 0 {
					distinct++
				}
				freq[d]++
				if freq[d] > maxFreq {
					maxFreq = freq[d]
				}
				if maxFreq <= distinct {
					ans++
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

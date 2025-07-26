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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)
		if n%2 == 1 {
			fmt.Fprintln(writer, -1)
			continue
		}
		freq := make([]int, 26)
		for i := 0; i < n; i++ {
			freq[s[i]-'a']++
		}
		maxFreq := 0
		for _, v := range freq {
			if v > maxFreq {
				maxFreq = v
			}
		}
		if maxFreq > n/2 {
			fmt.Fprintln(writer, -1)
			continue
		}
		pairCount := make([]int, 26)
		samePairs := 0
		for i := 0; i < n/2; i++ {
			if s[i] == s[n-1-i] {
				samePairs++
				pairCount[s[i]-'a']++
			}
		}
		maxPair := 0
		for _, v := range pairCount {
			if v > maxPair {
				maxPair = v
			}
		}
		ans := (samePairs + 1) / 2
		if maxPair > ans {
			ans = maxPair
		}
		fmt.Fprintln(writer, ans)
	}
}

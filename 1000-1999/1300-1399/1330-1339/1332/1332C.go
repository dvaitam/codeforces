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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)

		freq := make([][26]int, k)
		for i := 0; i < n; i++ {
			c := s[i] - 'a'
			idx := i % k
			freq[idx][c]++
		}
		m := n / k
		ans := 0
		for i := 0; i <= (k-1)/2; i++ {
			j := k - 1 - i
			counts := [26]int{}
			groupSize := m
			if i != j {
				groupSize = 2 * m
				for ch := 0; ch < 26; ch++ {
					counts[ch] = freq[i][ch] + freq[j][ch]
				}
			} else {
				for ch := 0; ch < 26; ch++ {
					counts[ch] = freq[i][ch]
				}
			}
			maxFreq := 0
			for ch := 0; ch < 26; ch++ {
				if counts[ch] > maxFreq {
					maxFreq = counts[ch]
				}
			}
			ans += groupSize - maxFreq
		}
		fmt.Fprintln(writer, ans)
	}
}

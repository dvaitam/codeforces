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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make(map[int]int, n)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
		}
		freqCount := make(map[int]int)
		freqs := make([]int, 0, len(freq))
		for _, f := range freq {
			freqs = append(freqs, f)
			freqCount[f]++
		}
		sort.Ints(freqs)
		uniqueFreq := make([]int, 0, len(freqCount))
		for i, f := range freqs {
			if i == 0 || f != freqs[i-1] {
				uniqueFreq = append(uniqueFreq, f)
			}
		}
		sort.Ints(uniqueFreq)
		ans := n
		suffix := 0
		for i := len(uniqueFreq) - 1; i >= 0; i-- {
			f := uniqueFreq[i]
			suffix += freqCount[f]
			keep := suffix * f
			if n-keep < ans {
				ans = n - keep
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

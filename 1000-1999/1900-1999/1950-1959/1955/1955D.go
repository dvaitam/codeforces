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
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		freqB := make(map[int]int, m)
		for i := 0; i < m; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freqB[x]++
		}

		freq := make(map[int]int)
		match := 0
		for i := 0; i < m; i++ {
			v := a[i]
			freq[v]++
			if freq[v] <= freqB[v] {
				match++
			}
		}
		count := 0
		if match >= k {
			count++
		}
		for i := m; i < n; i++ {
			outVal := a[i-m]
			if freq[outVal] <= freqB[outVal] {
				match--
			}
			freq[outVal]--
			if freq[outVal] == 0 {
				delete(freq, outVal)
			}

			v := a[i]
			freq[v]++
			if freq[v] <= freqB[v] {
				match++
			}
			if match >= k {
				count++
			}
		}
		fmt.Fprintln(writer, count)
	}
}

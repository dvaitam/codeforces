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
	var s string
	fmt.Fscan(reader, &s)
	var q int
	fmt.Fscan(reader, &q)

	// Precompute answers for each character and each possible number of changes
	best := make([][]int, 26)
	for i := range best {
		best[i] = make([]int, n+1)
	}

	bytes := []byte(s)
	for ch := 0; ch < 26; ch++ {
		prefix := make([]int, n+1)
		c := byte('a' + ch)
		for i := 0; i < n; i++ {
			if bytes[i] != c {
				prefix[i+1] = prefix[i] + 1
			} else {
				prefix[i+1] = prefix[i]
			}
		}
		minChanges := make([]int, n+1)
		for i := 1; i <= n; i++ {
			minChanges[i] = n + 1
		}
		for l := 0; l < n; l++ {
			for r := l; r < n; r++ {
				length := r - l + 1
				mism := prefix[r+1] - prefix[l]
				if mism < minChanges[length] {
					minChanges[length] = mism
				}
			}
		}
		for length := 1; length <= n; length++ {
			m := minChanges[length]
			if m <= n {
				for k := m; k <= n; k++ {
					if best[ch][k] < length {
						best[ch][k] = length
					}
				}
			}
		}
		for k := 1; k <= n; k++ {
			if best[ch][k] < best[ch][k-1] {
				best[ch][k] = best[ch][k-1]
			}
		}
	}

	for i := 0; i < q; i++ {
		var m int
		var str string
		fmt.Fscan(reader, &m, &str)
		if m > n {
			m = n
		}
		ch := str[0] - 'a'
		fmt.Fprintln(writer, best[ch][m])
	}
}

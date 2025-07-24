package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	n := len(s)
	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i]
		if s[i] == '1' {
			prefix[i+1]++
		}
	}

	B := int(math.Sqrt(float64(n))) + 1
	var ans int64

	for k := 1; k <= B; k++ {
		freq := make(map[int]int)
		for j := 0; j <= n; j++ {
			val := j - k*prefix[j]
			if c, ok := freq[val]; ok {
				ans += int64(c)
			}
			freq[val]++
		}
	}

	maxd := n / B
	for d := 1; d <= maxd; d++ {
		limit := d * (B + 1)
		for r := 0; r < d; r++ {
			freq := make(map[int]int)
			i := r
			for j := r; j <= n; j += d {
				for i <= j-limit {
					freq[prefix[i]]++
					i += d
				}
				ans += int64(freq[prefix[j]-d])
			}
		}
	}

	fmt.Fprintln(writer, ans)
}

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
		n := len(s)
		best := 0
		// single digit subsequences
		freq := make([]int, 10)
		for i := 0; i < n; i++ {
			d := int(s[i] - '0')
			freq[d]++
		}
		for i := 0; i < 10; i++ {
			if freq[i] > best {
				best = freq[i]
			}
		}
		// alternating subsequences for every pair of digits
		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				if i == j {
					continue
				}
				expect := i
				length := 0
				for k := 0; k < n; k++ {
					d := int(s[k] - '0')
					if d == expect {
						length++
						if expect == i {
							expect = j
						} else {
							expect = i
						}
					}
				}
				if length%2 == 1 {
					length--
				}
				if length > best {
					best = length
				}
			}
		}
		fmt.Fprintln(writer, n-best)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problemC from contest 1847. Given the strengths of the
// initial Stand users, DIO can repeatedly append the XOR of any suffix of the
// current list. Observing that after each operation the prefix XOR of the list
// becomes one of the previously seen prefix XORs, any appended value can be
// represented as the XOR of two prefix XORs from the original array. Therefore
// the maximum achievable strength is simply the maximum XOR of any two prefix
// XORs of the initial sequence.
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
		prefix := 0
		seen := make([]bool, 256)
		seen[0] = true
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			prefix ^= x
			seen[prefix] = true
		}
		ans := 0
		for i := 0; i < 256; i++ {
			if !seen[i] {
				continue
			}
			for j := 0; j < 256; j++ {
				if seen[j] {
					v := i ^ j
					if v > ans {
						ans = v
					}
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

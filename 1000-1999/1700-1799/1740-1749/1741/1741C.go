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
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(reader, &a[i])
		}

		prefix := make([]int, n+1)
		pos := make(map[int]int, n+1)
		pos[0] = 0
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + a[i]
			pos[prefix[i+1]] = i + 1
		}

		ans := n
		for i := 1; i <= n; i++ {
			target := prefix[i]
			last := i
			maxLen := i
			valid := true
			for last < n {
				needed := prefix[last] + target
				idx, ok := pos[needed]
				if !ok {
					valid = false
					break
				}
				segLen := idx - last
				if segLen > maxLen {
					maxLen = segLen
				}
				last = idx
			}
			if valid && last == n {
				if maxLen < ans {
					ans = maxLen
				}
			}
		}

		fmt.Fprintln(writer, ans)
	}
}

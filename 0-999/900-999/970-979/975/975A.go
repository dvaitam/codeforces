package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	roots := make(map[string]bool)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		seen := make([]bool, 26)
		for _, ch := range s {
			seen[ch-'a'] = true
		}
		var root []byte
		for j := 0; j < 26; j++ {
			if seen[j] {
				root = append(root, byte('a'+j))
			}
		}
		roots[string(root)] = true
	}
	fmt.Fprintln(out, len(roots))
}

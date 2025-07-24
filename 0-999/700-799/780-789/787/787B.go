package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	for g := 0; g < m; g++ {
		var k int
		fmt.Fscan(in, &k)
		seen := make(map[int]bool, k)
		ok := true
		for i := 0; i < k; i++ {
			var v int
			fmt.Fscan(in, &v)
			if seen[-v] {
				ok = false
			}
			seen[v] = true
		}
		if ok {
			fmt.Println("YES")
			return
		}
	}
	fmt.Println("NO")
}

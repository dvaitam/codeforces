package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	longest := 0
	for i := 0; i < k; i++ {
		var m int
		fmt.Fscan(in, &m)
		prefix := 0
		for j := 0; j < m; j++ {
			var x int
			fmt.Fscan(in, &x)
			if j == prefix && x == j+1 {
				prefix++
			}
		}
		if prefix > longest {
			longest = prefix
		}
	}
	result := 2*(n-longest) + 1 - k
	fmt.Println(result)
}

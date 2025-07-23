package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var dirs string
	fmt.Fscan(reader, &dirs)
	pos := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &pos[i])
	}
	// any possible collision time fits well below 1e9
	const inf int = int(1e9 + 7)
	ans := inf
	for i := 0; i < n-1; i++ {
		if dirs[i] == 'R' && dirs[i+1] == 'L' {
			d := pos[i+1] - pos[i]
			if d/2 < ans {
				ans = d / 2
			}
		}
	}
	if ans == inf {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}

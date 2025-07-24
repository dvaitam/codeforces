package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(reader, &n, &k)
	var s string
	fmt.Fscan(reader, &s)
	posG, posT := -1, -1
	for i, ch := range s {
		if ch == 'G' {
			posG = i
		} else if ch == 'T' {
			posT = i
		}
	}
	// Check if reachable by jumps of length k
	if (posT-posG)%k != 0 {
		fmt.Println("NO")
		return
	}
	// Determine jump direction
	step := k
	if posT < posG {
		step = -k
	}
	// Simulate jumps
	for cur := posG; ; cur += step {
		if s[cur] == '#' {
			fmt.Println("NO")
			return
		}
		if cur == posT {
			fmt.Println("YES")
			return
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	grundy := make([]int, 61)
	for s := 1; s <= 60; s++ {
		g := 1
		for (g+1)*(g+2)/2 <= s {
			g++
		}
		grundy[s] = g
	}
	xor := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		xor ^= grundy[x]
	}
	if xor == 0 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

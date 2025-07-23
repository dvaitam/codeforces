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
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(reader, &p)
		children[p] = append(children[p], i)
	}
	for v := 1; v <= n; v++ {
		if len(children[v]) == 0 {
			continue
		}
		leafChildren := 0
		for _, ch := range children[v] {
			if len(children[ch]) == 0 {
				leafChildren++
			}
		}
		if leafChildren < 3 {
			fmt.Println("No")
			return
		}
	}
	fmt.Println("Yes")
}

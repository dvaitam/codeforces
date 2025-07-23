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
	jumps := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &jumps[i])
	}
	visited := make([]bool, n)
	pos := 0
	for {
		if pos < 0 || pos >= n {
			fmt.Println("FINITE")
			return
		}
		if visited[pos] {
			fmt.Println("INFINITE")
			return
		}
		visited[pos] = true
		if dirs[pos] == '>' {
			pos += jumps[pos]
		} else {
			pos -= jumps[pos]
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b, c int
	if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
		return
	}
	remainder := a % b
	visited := make(map[int]bool)
	for pos := 1; ; pos++ {
		remainder *= 10
		digit := remainder / b
		if digit == c {
			fmt.Println(pos)
			return
		}
		remainder %= b
		if remainder == 0 {
			if c == 0 {
				fmt.Println(pos + 1)
			} else {
				fmt.Println(-1)
			}
			return
		}
		if visited[remainder] {
			fmt.Println(-1)
			return
		}
		visited[remainder] = true
	}
}

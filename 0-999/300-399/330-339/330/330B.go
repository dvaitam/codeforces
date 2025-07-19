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

	var M, n int
	for {
		if _, err := fmt.Fscan(reader, &M, &n); err != nil {
			break
		}
		vis := make([]int, M+1)
		for i := 0; i < n; i++ {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			vis[a]++
			vis[b]++
		}
		var izero int
		for i := 1; i <= M; i++ {
			if vis[i] == 0 {
				izero = i
				break
			}
		}
		fmt.Fprintln(writer, M-1)
		for j := 1; j <= M; j++ {
			if j != izero {
				fmt.Fprintln(writer, izero, j)
			}
		}
	}
}

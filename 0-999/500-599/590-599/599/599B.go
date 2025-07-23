package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	pos := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		pos[v] = append(pos[v], i)
	}
	ans := make([]int, m)
	ambiguous := false
	for i := 0; i < m; i++ {
		var b int
		fmt.Fscan(reader, &b)
		if len(pos[b]) == 0 {
			fmt.Println("Impossible")
			return
		}
		if len(pos[b]) > 1 {
			ambiguous = true
		}
		ans[i] = pos[b][0]
	}
	if ambiguous {
		fmt.Println("Ambiguity")
		return
	}
	fmt.Println("Possible")
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i, v := range ans {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}

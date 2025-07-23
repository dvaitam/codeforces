package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var adj = [][]int{
	{1, 2}, // index 0
	{0, 3}, // index 1
	{0, 3}, // index 2
	{1, 2}, // index 3
}

func neighbors(state string) []string {
	idx := strings.IndexByte(state, 'X')
	res := make([]string, 0, len(adj[idx]))
	b := []byte(state)
	for _, j := range adj[idx] {
		b[idx], b[j] = b[j], b[idx]
		res = append(res, string(b))
		b[idx], b[j] = b[j], b[idx]
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a1, a2, b1, b2 string
	if _, err := fmt.Fscan(reader, &a1); err != nil {
		return
	}
	fmt.Fscan(reader, &a2)
	fmt.Fscan(reader, &b1)
	fmt.Fscan(reader, &b2)
	start := a1 + a2
	target := b1 + b2

	if start == target {
		fmt.Println("YES")
		return
	}

	visited := map[string]bool{start: true}
	queue := []string{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, nxt := range neighbors(cur) {
			if !visited[nxt] {
				if nxt == target {
					fmt.Println("YES")
					return
				}
				visited[nxt] = true
				queue = append(queue, nxt)
			}
		}
	}
	fmt.Println("NO")
}

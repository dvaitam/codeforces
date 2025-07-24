package main

import (
	"bufio"
	"fmt"
	"os"
)

func match(s string, pos int, t string) int {
	j := pos
	for i := 0; i < len(t) && j < len(s); i++ {
		if t[i] == s[j] {
			j++
		}
	}
	return j
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]string, n)
	for i := range arr {
		fmt.Fscan(in, &arr[i])
	}
	var s string
	fmt.Fscan(in, &s)
	m := len(s)
	dist := make([]int, m+1)
	for i := range dist {
		dist[i] = -1
	}
	queue := make([]int, 0, m+1)
	dist[0] = 0
	queue = append(queue, 0)
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		if pos == m {
			fmt.Println(dist[pos])
			return
		}
		for _, t := range arr {
			np := match(s, pos, t)
			if np > pos && dist[np] == -1 {
				dist[np] = dist[pos] + 1
				queue = append(queue, np)
			}
		}
	}
	if dist[m] == -1 {
		fmt.Println(-1)
	} else {
		fmt.Println(dist[m])
	}
}

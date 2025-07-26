package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(reader, &t)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; t > 0; t-- {
		var n, a, b int
		fmt.Fscan(reader, &n, &a, &b)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		dist := bfs(n, a, b, p)
		fmt.Fprintln(writer, dist)
	}
}

func bfs(n, a, b int, p []int) int {
	next := make([]int, n+2)
	prev := make([]int, n+2)
	for i := 0; i <= n+1; i++ {
		next[i] = i
		prev[i] = i
	}

	var findNext func(int) int
	findNext = func(x int) int {
		if x > n+1 {
			return n + 1
		}
		if next[x] != x {
			next[x] = findNext(next[x])
		}
		return next[x]
	}

	var findPrev func(int) int
	findPrev = func(x int) int {
		if x < 0 {
			return 0
		}
		if prev[x] != x {
			prev[x] = findPrev(prev[x])
		}
		return prev[x]
	}

	removeNext := func(x int) { next[x] = findNext(x + 1) }
	removePrev := func(x int) { prev[x] = findPrev(x - 1) }

	visited := make([]bool, n+2)
	distance := make([]int, n+2)
	queue := []int{a}
	visited[a] = true
	removeNext(a)
	removePrev(a)
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		if v == b {
			return distance[v]
		}
		// explore to the right
		r := v + p[v-1]
		if r > n {
			r = n
		}
		for j := findNext(v + 1); j <= r; j = findNext(j) {
			if visited[j] {
				removeNext(j)
				continue
			}
			if p[j-1] >= j-v {
				visited[j] = true
				distance[j] = distance[v] + 1
				queue = append(queue, j)
				removePrev(j)
			}
			removeNext(j)
		}
		// explore to the left
		l := v - p[v-1]
		if l < 1 {
			l = 1
		}
		for j := findPrev(v - 1); j >= l; j = findPrev(j) {
			if visited[j] {
				removePrev(j)
				continue
			}
			if p[j-1] >= v-j {
				visited[j] = true
				distance[j] = distance[v] + 1
				queue = append(queue, j)
				removeNext(j)
			}
			removePrev(j)
		}
	}
	return distance[b]
}

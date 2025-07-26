package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var x int64
	if _, err := fmt.Fscan(in, &n, &x); err != nil {
		return
	}
	fmt.Fprintln(out, bfs(n, x))
}

type node struct {
	val   int64
	steps int
}

func bfs(n int, x int64) int {
	limit := int64(1)
	for i := 0; i < n; i++ {
		limit *= 10
	}

	queue := []node{{x, 0}}
	visited := map[int64]bool{x: true}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if digitsLen(cur.val) == n {
			return cur.steps
		}
		digits := uniqueDigits(cur.val)
		for _, d := range digits {
			if d <= 1 {
				continue
			}
			next := cur.val * int64(d)
			if next >= limit {
				continue
			}
			if !visited[next] {
				visited[next] = true
				queue = append(queue, node{next, cur.steps + 1})
			}
		}
	}
	return -1
}

func digitsLen(x int64) int {
	if x == 0 {
		return 1
	}
	l := 0
	for x > 0 {
		x /= 10
		l++
	}
	return l
}

func uniqueDigits(x int64) []int {
	seen := [10]bool{}
	res := make([]int, 0, 10)
	if x == 0 {
		return []int{0}
	}
	for x > 0 {
		d := int(x % 10)
		if !seen[d] {
			seen[d] = true
			res = append(res, d)
		}
		x /= 10
	}
	return res
}

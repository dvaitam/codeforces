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

	var n int64
	fmt.Fscan(in, &n)
	if n == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	visited := make(map[int64]int)
	type state struct {
		val  int64
		dist int
	}
	queue := []state{{n, 0}}
	visited[n] = 0
	head := 0
	for head < len(queue) {
		cur := queue[head]
		head++
		if cur.val == 0 {
			fmt.Fprintln(out, cur.dist)
			return
		}
		digits := digitsOf(cur.val)
		for _, digit := range digits {
			next := cur.val - int64(digit)
			if _, ok := visited[next]; !ok {
				visited[next] = cur.dist + 1
				queue = append(queue, state{next, cur.dist + 1})
			}
		}
	}
}

func digitsOf(x int64) []int {
	if x == 0 {
		return []int{0}
	}
	digits := make([]int, 0, 19)
	for v := x; v > 0; v /= 10 {
		d := int(v % 10)
		if d > 0 {
			digits = append(digits, d)
		}
	}
	return digits
}

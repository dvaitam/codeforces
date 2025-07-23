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

	var n int
	fmt.Fscan(reader, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	dist := make([]int, n)
	res := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}

	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		if i == 0 || i == n-1 || (i > 0 && arr[i] == arr[i-1]) || (i+1 < n && arr[i] == arr[i+1]) {
			dist[i] = 0
			res[i] = arr[i]
			queue = append(queue, i)
		}
	}

	for head := 0; head < len(queue); head++ {
		i := queue[head]
		for _, j := range []int{i - 1, i + 1} {
			if j >= 0 && j < n && dist[j] == -1 {
				dist[j] = dist[i] + 1
				res[j] = res[i]
				queue = append(queue, j)
			}
		}
	}

	maxStep := 0
	for _, d := range dist {
		if d > maxStep {
			maxStep = d
		}
	}

	fmt.Fprintln(writer, maxStep)
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}

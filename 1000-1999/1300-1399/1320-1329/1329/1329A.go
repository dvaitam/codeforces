package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N, M int
	if _, err := fmt.Fscan(reader, &N, &M); err != nil {
		return
	}
	A := make([]int, M)
	for i := 0; i < M; i++ {
		fmt.Fscan(reader, &A[i])
	}

	sum := 0
	for _, v := range A {
		sum += v
	}
	if sum < N {
		fmt.Fprint(writer, -1)
		return
	}
	for i, v := range A {
		if i+v > N {
			fmt.Fprint(writer, -1)
			return
		}
	}

	out := make([]int, M+1)
	for i := 0; i <= M; i++ {
		out[i] = i + 1
	}
	out[M] = N + 1

	for i := M - 1; i >= 0; i-- {
		pos := out[i] + A[i] - 1
		move := max(0, out[i+1]-pos-1)
		out[i] += move
	}

	for i := 0; i < M; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, out[i])
	}
	writer.WriteByte('\n')
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N, M int
	if _, err := fmt.Fscan(reader, &N, &M); err != nil {
		return
	}

	R := make([]int, 0, N)
	T := make([]int, 0, M)
	for i := 0; i < N+M; i++ {
		var x, t int
		fmt.Fscan(reader, &x, &t)
		if t == 0 {
			R = append(R, x)
		} else {
			T = append(T, x)
		}
	}

	ans := make([]int, M)
	j := 0
	for i := 0; i < M; i++ {
		for j < N && (i == M-1 || abs(R[j]-T[i]) <= abs(R[j]-T[i+1])) {
			ans[i]++
			j++
		}
	}

	for i := 0; i < M; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}

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

	var N, K int
	if _, err := fmt.Fscan(reader, &N, &K); err != nil {
		return
	}
	A := make([]int, N+2)
	s := 0
	n := 0
	for i := 1; i <= N; i++ {
		for j := i * 2; j <= N; j += i {
			A[j]++
		}
		s += A[i]
		if s >= K {
			n = i
			break
		}
	}
	if s < K {
		fmt.Fprintln(writer, "No")
		return
	}
	s -= K
	k := n
	removed := make([]bool, n+2)
	for i := 2; i <= n && s > 0; i++ {
		cost := A[i] + n/i - 1
		if s >= cost {
			s -= cost
			removed[i] = true
			k--
		}
	}
	fmt.Fprintln(writer, "Yes")
	fmt.Fprintln(writer, k)
	for i := 1; i <= n; i++ {
		if !removed[i] {
			fmt.Fprintf(writer, "%d ", i)
		}
	}
	fmt.Fprintln(writer)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var N, M int
	fmt.Fscan(reader, &N, &M)
	a := make([][]int, N)
	for i := 0; i < N; i++ {
		a[i] = make([]int, M)
		for j := 0; j < M; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}
	b := make([][]int, N)
	for i := 0; i < N; i++ {
		b[i] = make([]int, M)
		for j := 0; j < M; j++ {
			fmt.Fscan(reader, &b[i][j])
		}
	}
	for sum := 0; sum <= N+M-2; sum++ {
		var diagA []int
		var diagB []int
		for i := 0; i < N; i++ {
			j := sum - i
			if j >= 0 && j < M {
				diagA = append(diagA, a[i][j])
				diagB = append(diagB, b[i][j])
			}
		}
		sort.Ints(diagA)
		sort.Ints(diagB)
		if len(diagA) != len(diagB) {
			fmt.Println("NO")
			return
		}
		for k := range diagA {
			if diagA[k] != diagB[k] {
				fmt.Println("NO")
				return
			}
		}
	}
	fmt.Println("YES")
}

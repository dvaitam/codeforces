package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([][]int, 2)
	b := make([][]int, 2)
	for i := 0; i < 2; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}
	for i := 0; i < 2; i++ {
		b[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &b[i][j])
		}
	}

	sumA, sumB := 0, 0
	for i := 0; i < 2; i++ {
		for j := 0; j < n; j++ {
			sumA += a[i][j]
			sumB += b[i][j]
		}
	}
	if sumA != sumB {
		fmt.Println(-1)
		return
	}

	c0, c1 := 0, 0
	cost := 0
	for i := 0; i < n; i++ {
		u := c0 + a[0][i] - b[0][i]
		v := c1 + a[1][i] - b[1][i]
		arr := []int{0, u, -v}
		sort.Ints(arr)
		x := arr[1]
		cost += abs(x) + abs(u-x) + abs(v+x)
		c0 = u - x
		c1 = v + x
	}
	if c0 != 0 || c1 != 0 {
		fmt.Println(-1)
		return
	}
	fmt.Println(cost)
}

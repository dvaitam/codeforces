package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	mat := make([][]int64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &mat[i][j])
		}
	}
	a := make([]int64, n)
	b := make([]int64, m)
	for j := 0; j < m; j++ {
		b[j] = mat[0][j]
	}
	for i := 0; i < n; i++ {
		a[i] = mat[i][0] - b[0]
	}
	var g int64
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x := abs(mat[i][j] - a[i] - b[j])
			g = gcd(g, x)
		}
	}
	if g == 0 {
		g = 1000000001
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if g <= mat[i][j] {
				fmt.Println("NO")
				return
			}
		}
	}
	fmt.Println("YES")
	for i := 0; i < n; i++ {
		if a[i] < 0 {
			a[i] += g
		}
	}
	fmt.Println(g)
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(a[i])
	}
	fmt.Println()
	for j := 0; j < m; j++ {
		if j > 0 {
			fmt.Print(" ")
		}
		fmt.Print(b[j])
	}
	fmt.Println()
}

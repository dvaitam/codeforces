package main

import (
	"bufio"
	"fmt"
	"os"
)

func generate(n int) [][]int {
	mat := make([][]int, n)
	for i := range mat {
		mat[i] = make([]int, n)
	}
	half := n / 2
	for i := 0; i < half; i++ {
		for j := 0; j < half; j++ {
			d := (i + j) % half
			mat[i][j] = d
			mat[i][n-1-j] = half + (half-i+d)%half
			mat[n-1-i][j] = half + (half-j+d)%half
			mat[n-1-i][n-1-j] = n - 1 - d
		}
	}
	return mat
}

func shiftPattern(mat [][]int, k int) [][]int {
	n := len(mat)
	res := make([][]int, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		copy(res[i], mat[i])
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			res[i][j] = (res[i][j] + k*i) % n
		}
	}
	return res
}

func solveCase(n int) [][]int {
	if n == 1 {
		return [][]int{{0}}
	}
	if n%2 == 1 {
		return nil
	}
	base := generate(n)
	return shiftPattern(base, 1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		ans := solveCase(n)
		if ans == nil {
			fmt.Fprintln(out, "NO")
			continue
		}
		fmt.Fprintln(out, "YES")
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, ans[i][j])
			}
			fmt.Fprintln(out)
		}
	}
}

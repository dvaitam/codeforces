package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	n, m int
	A    [][]int
	S    int

	x, y []int64
	X, Y []int64
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}
	scanInt64 := func() int64 {
		scanner.Scan()
		val, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		return val
	}

	if !scanner.Scan() {
		return
	}
	n, _ = strconv.Atoi(scanner.Text())
	m = scanInt()

	x = make([]int64, n)
	y = make([]int64, n)
	for i := 0; i < n; i++ {
		x[i] = scanInt64()
		y[i] = scanInt64()
	}

	X = make([]int64, m)
	Y = make([]int64, m)
	for i := 0; i < m; i++ {
		X[i] = scanInt64()
		Y[i] = scanInt64()
	}

	A = make([][]int, n)
	for i := range A {
		A[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if x[i] < x[j] {
				a := y[j] - y[i]
				b := x[i] - x[j]
				c_val := -x[i]*a - y[i]*b

				for k := 0; k < m; k++ {
					if X[k] > x[i] && X[k] <= x[j] && a*X[k]+b*Y[k]+c_val > 0 {
						A[i][j]++
					}
				}
				A[j][i] = -A[i][j]
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if A[i][j]+A[j][k]+A[k][i] == 0 {
					S++
				}
			}
		}
	}

	fmt.Printf("%d\n", S)
}

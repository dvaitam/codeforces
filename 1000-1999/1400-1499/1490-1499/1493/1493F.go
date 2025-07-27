package main

import (
	"bufio"
	"fmt"
	"os"
)

func primeFactorsDistinct(x int) []int {
	res := []int{}
	for p := 2; p*p <= x; p++ {
		if x%p == 0 {
			res = append(res, p)
			for x%p == 0 {
				x /= p
			}
		}
	}
	if x > 1 {
		res = append(res, x)
	}
	return res
}

func checkRows(a [][]int, period int) bool {
	n := len(a)
	m := len(a[0])
	for i := period; i < n; i++ {
		for j := 0; j < m; j++ {
			if a[i][j] != a[i-period][j] {
				return false
			}
		}
	}
	return true
}

func checkCols(a [][]int, period int) bool {
	n := len(a)
	m := len(a[0])
	for j := period; j < m; j++ {
		for i := 0; i < n; i++ {
			if a[i][j] != a[i][j-period] {
				return false
			}
		}
	}
	return true
}

func numDivisors(x int) int {
	if x == 0 {
		return 0
	}
	res := 1
	for p := 2; p*p <= x; p++ {
		if x%p == 0 {
			e := 0
			for x%p == 0 {
				x /= p
				e++
			}
			res *= e + 1
		}
	}
	if x > 1 {
		res *= 2
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}

	rowPeriod := n
	for _, p := range primeFactorsDistinct(rowPeriod) {
		for rowPeriod%p == 0 && checkRows(a, rowPeriod/p) {
			rowPeriod /= p
		}
	}

	colPeriod := m
	for _, p := range primeFactorsDistinct(colPeriod) {
		for colPeriod%p == 0 && checkCols(a, colPeriod/p) {
			colPeriod /= p
		}
	}

	cntRows := numDivisors(n / rowPeriod)
	cntCols := numDivisors(m / colPeriod)

	fmt.Fprintln(out, cntRows*cntCols)
}

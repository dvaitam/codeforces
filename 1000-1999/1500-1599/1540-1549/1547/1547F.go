package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func buildSparse(arr []int) ([][]int, []int) {
	n := len(arr)
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}
	K := log[n] + 1
	st := make([][]int, K)
	st[0] = make([]int, n)
	copy(st[0], arr)
	for k := 1; k < K; k++ {
		size := n - (1 << k) + 1
		st[k] = make([]int, size)
		step := 1 << (k - 1)
		for i := 0; i < size; i++ {
			st[k][i] = gcd(st[k-1][i], st[k-1][i+step])
		}
	}
	return st, log
}

func query(st [][]int, log []int, l, r int) int {
	length := r - l
	k := log[length]
	return gcd(st[k][l], st[k][r-(1<<k)])
}

func allEqual(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[0] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if allEqual(a) {
			fmt.Fprintln(writer, 0)
			continue
		}
		g := 0
		for _, v := range a {
			g = gcd(g, v)
		}
		b := make([]int, 2*n)
		copy(b, a)
		copy(b[n:], a)
		st, log := buildSparse(b)
		low, high := 1, n
		for low < high {
			mid := (low + high) / 2
			ok := true
			for i := 0; i < n; i++ {
				if query(st, log, i, i+mid) != g {
					ok = false
					break
				}
			}
			if ok {
				high = mid
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintln(writer, low-1)
	}
}

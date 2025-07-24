package main

import (
	"bufio"
	"fmt"
	"os"
)

func equalRotation(a, b []int) bool {
	n := len(a)
	for shift := 0; shift < n; shift++ {
		ok := true
		for i := 0; i < n; i++ {
			if a[i] != b[(i+shift)%n] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, L int
	if _, err := fmt.Fscan(reader, &n, &L); err != nil {
		return
	}
	A := make([]int, n)
	B := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &A[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &B[i])
	}
	da := make([]int, n)
	db := make([]int, n)
	for i := 0; i < n-1; i++ {
		da[i] = A[i+1] - A[i]
		db[i] = B[i+1] - B[i]
	}
	da[n-1] = L - A[n-1] + A[0]
	db[n-1] = L - B[n-1] + B[0]

	if equalRotation(da, db) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

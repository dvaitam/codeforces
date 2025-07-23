package main

import (
	"bufio"
	"fmt"
	"os"
)

func compute(arr []int, rooms, d, b int) int {
	n := len(arr)
	p := 0
	var avail int64
	miss := 0
	for i := 1; i <= rooms; i++ {
		limit := i * (d + 1)
		if limit > n {
			limit = n
		}
		for p < limit {
			avail += int64(arr[p])
			p++
		}
		if avail >= int64(b) {
			avail -= int64(b)
		} else {
			miss++
		}
	}
	return miss
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, d, b int
	if _, err := fmt.Fscan(reader, &n, &d, &b); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	leftRooms := (n + 1) / 2
	rightRooms := n / 2

	left := compute(a, leftRooms, d, b)

	// reverse array for right side
	rev := make([]int, n)
	for i := 0; i < n; i++ {
		rev[i] = a[n-1-i]
	}
	right := compute(rev, rightRooms, d, b)

	if left > right {
		fmt.Println(left)
	} else {
		fmt.Println(right)
	}
}

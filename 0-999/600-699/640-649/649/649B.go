package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	var a, b int
	fmt.Fscan(in, &a, &b)

	perEntrance := m * k
	entranceA := (a-1)/perEntrance + 1
	entranceB := (b-1)/perEntrance + 1
	floorA := ((a-1)%perEntrance)/k + 1
	floorB := ((b-1)%perEntrance)/k + 1

	downA := min((floorA-1)*5, 10+(floorA-1))
	upB := min((floorB-1)*5, 10+(floorB-1))

	diff := abs(entranceA - entranceB)
	walk := min(diff, n-diff) * 15

	fmt.Println(downA + walk + upB)
}

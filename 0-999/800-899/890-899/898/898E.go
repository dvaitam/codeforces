package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func isSquare(x int) bool {
	r := int(math.Sqrt(float64(x)))
	return r*r == x
}

func costToSquare(x int) int {
	r := int(math.Sqrt(float64(x)))
	c1 := x - r*r
	c2 := (r+1)*(r+1) - x
	if c1 < c2 {
		return c1
	}
	return c2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	squareCosts := []int{}
	nonsquareCosts := []int{}
	squares := 0

	for _, v := range a {
		if isSquare(v) {
			squares++
			if v == 0 {
				squareCosts = append(squareCosts, 2)
			} else {
				squareCosts = append(squareCosts, 1)
			}
		} else {
			nonsquareCosts = append(nonsquareCosts, costToSquare(v))
		}
	}

	target := n / 2
	if squares == target {
		fmt.Fprintln(writer, 0)
		return
	}
	if squares > target {
		sort.Ints(squareCosts)
		moves := 0
		need := squares - target
		for i := 0; i < need; i++ {
			moves += squareCosts[i]
		}
		fmt.Fprintln(writer, moves)
	} else {
		sort.Ints(nonsquareCosts)
		moves := 0
		need := target - squares
		for i := 0; i < need; i++ {
			moves += nonsquareCosts[i]
		}
		fmt.Fprintln(writer, moves)
	}
}

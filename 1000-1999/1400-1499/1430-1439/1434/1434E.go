package main

import (
	"bufio"
	"fmt"
	"os"
)

var arr []int
var m int
var memo map[[2]int]int

func sg(i int, diff int) int {
	key := [2]int{i, diff}
	if v, ok := memo[key]; ok {
		return v
	}
	moves := make(map[int]bool)
	for j := i + 1; j < m; j++ {
		d := arr[j] - arr[i]
		if d > diff {
			g := sg(j, d)
			moves[g] = true
		}
	}
	g := 0
	for moves[g] {
		g++
	}
	memo[key] = g
	return g
}

func sgSingle(i int) int {
	moves := make(map[int]bool)
	for j := i + 1; j < m; j++ {
		d := arr[j] - arr[i]
		g := sg(j, d)
		moves[g] = true
	}
	g := 0
	for moves[g] {
		g++
	}
	return g
}

func grundy(sequence []int) int {
	arr = sequence
	m = len(arr)
	memo = make(map[[2]int]int)
	moves := make(map[int]bool)
	for i := 0; i < m; i++ {
		g := sgSingle(i)
		moves[g] = true
	}
	g := 0
	for moves[g] {
		g++
	}
	return g
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	xor := 0
	for i := 0; i < n; i++ {
		var length int
		fmt.Fscan(in, &length)
		seq := make([]int, length)
		for j := 0; j < length; j++ {
			fmt.Fscan(in, &seq[j])
		}
		g := grundy(seq)
		xor ^= g
	}
	if xor != 0 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

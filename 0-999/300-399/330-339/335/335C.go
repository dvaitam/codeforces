package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	leftBit  = 1
	rightBit = 2
)

var memo = map[string]int{}

func encode(seg []int) string {
	b := make([]byte, len(seg))
	for i, v := range seg {
		b[i] = byte(v)
	}
	return string(b)
}

func connected(a, b int) bool {
	if a == 0 || b == 0 {
		return false
	}
	return (a&leftBit != 0 && b&rightBit != 0) || (a&rightBit != 0 && b&leftBit != 0)
}

func splitSegments(arr []int) [][]int {
	res := [][]int{}
	n := len(arr)
	i := 0
	for i < n {
		for i < n && arr[i] == 0 {
			i++
		}
		if i >= n {
			break
		}
		j := i
		for j+1 < n {
			if arr[j+1] == 0 {
				break
			}
			if !connected(arr[j], arr[j+1]) {
				break
			}
			j++
		}
		seg := make([]int, j-i+1)
		copy(seg, arr[i:j+1])
		res = append(res, seg)
		i = j + 1
	}
	return res
}

func grundy(seg []int) int {
	if len(seg) == 0 {
		return 0
	}
	key := encode(seg)
	if g, ok := memo[key]; ok {
		return g
	}
	mexUsed := map[int]struct{}{}
	for i, state := range seg {
		if state&leftBit != 0 {
			g := evaluateMove(seg, i, true)
			mexUsed[g] = struct{}{}
		}
		if state&rightBit != 0 {
			g := evaluateMove(seg, i, false)
			mexUsed[g] = struct{}{}
		}
	}
	g := 0
	for {
		if _, ok := mexUsed[g]; !ok {
			break
		}
		g++
	}
	memo[key] = g
	return g
}

func evaluateMove(seg []int, idx int, pickLeft bool) int {
	newSeg := make([]int, len(seg))
	copy(newSeg, seg)
	newSeg[idx] = 0
	var oppBit int
	if pickLeft {
		oppBit = rightBit
	} else {
		oppBit = leftBit
	}
	for delta := -1; delta <= 1; delta++ {
		if delta == 0 {
			continue
		}
		nr := idx + delta
		if nr >= 0 && nr < len(seg) {
			newSeg[nr] &^= oppBit
		}
	}
	subSegs := splitSegments(newSeg)
	res := 0
	for _, s := range subSegs {
		res ^= grundy(s)
	}
	return res
}

func initialReclaim(states []int, row int, col int) {
	row--
	if row < 0 || row >= len(states) {
		return
	}
	// remove entire row
	states[row] = 0
	var oppBit int
	if col == 1 {
		oppBit = rightBit
	} else {
		oppBit = leftBit
	}
	for delta := -1; delta <= 1; delta++ {
		nr := row + delta
		if nr >= 0 && nr < len(states) {
			if nr == row {
				continue
			}
			states[nr] &^= oppBit
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var r, n int
	if _, err := fmt.Fscan(in, &r, &n); err != nil {
		return
	}
	states := make([]int, r)
	for i := range states {
		states[i] = leftBit | rightBit
	}
	for i := 0; i < n; i++ {
		var row, col int
		fmt.Fscan(in, &row, &col)
		initialReclaim(states, row, col)
	}

	segments := splitSegments(states)
	total := 0
	for _, seg := range segments {
		total ^= grundy(seg)
	}
	if total != 0 {
		fmt.Println("WIN")
	} else {
		fmt.Println("LOSE")
	}
}

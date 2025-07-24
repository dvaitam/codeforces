package main

import (
	"bufio"
	"fmt"
	"os"
)

type Op struct {
	t int
	x int
	y int
}

var n, k, q int
var a [][]int
var ops []Op
var memo map[int64]int

func key(node, idx int) int64 {
	return int64(node)<<32 | int64(idx)
}

func eval(node, idx int) int {
	keyVal := key(node, idx)
	if val, ok := memo[keyVal]; ok {
		return val
	}
	stack := [][2]int{{node, idx}}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		nid, j := top[0], top[1]
		kkey := key(nid, j)
		if _, ok := memo[kkey]; ok {
			stack = stack[:len(stack)-1]
			continue
		}
		if nid <= k {
			memo[kkey] = a[nid-1][j-1]
			stack = stack[:len(stack)-1]
			continue
		}
		op := ops[nid-1]
		keyX := key(op.x, j)
		keyY := key(op.y, j)
		_, okX := memo[keyX]
		_, okY := memo[keyY]
		if !okX {
			stack = append(stack, [2]int{op.x, j})
			continue
		}
		if !okY {
			stack = append(stack, [2]int{op.y, j})
			continue
		}
		valX := memo[keyX]
		valY := memo[keyY]
		if op.t == 1 { // max
			if valX >= valY {
				memo[kkey] = valX
			} else {
				memo[kkey] = valY
			}
		} else { // min
			if valX <= valY {
				memo[kkey] = valX
			} else {
				memo[kkey] = valY
			}
		}
		stack = stack[:len(stack)-1]
	}
	return memo[keyVal]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &k, &q); err != nil {
		return
	}

	a = make([][]int, k)
	for i := 0; i < k; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}

	// Prepare slice for all creatures (initial k and up to q new ones).
	ops = make([]Op, k+q+1)
	for i := 1; i <= k; i++ {
		ops[i-1] = Op{t: 0}
	}

	memo = make(map[int64]int)
	nextID := k + 1
	for i := 0; i < q; i++ {
		var t, x, y int
		fmt.Fscan(reader, &t, &x, &y)
		if t == 3 {
			res := eval(x, y)
			fmt.Fprintln(writer, res)
		} else {
			ops[nextID-1] = Op{t: t, x: x, y: y}
			nextID++
		}
	}
}

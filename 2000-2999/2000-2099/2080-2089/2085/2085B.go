package main

import (
	"bufio"
	"fmt"
	"os"
)

type operation struct {
	l, r int
}

func mexSegment(seg []int) int {
	seen := make([]bool, len(seg)+2)
	for _, v := range seg {
		if v >= 0 && v < len(seen) {
			seen[v] = true
		}
	}
	for i, ok := range seen {
		if !ok {
			return i
		}
	}
	return len(seen)
}

func applyOperation(arr []int, l, r int) []int {
	mex := mexSegment(arr[l : r+1])
	newArr := make([]int, 0, len(arr)-(r-l))
	newArr = append(newArr, arr[:l]...)
	newArr = append(newArr, mex)
	newArr = append(newArr, arr[r+1:]...)
	return newArr
}

func findZero(arr []int) int {
	for i, v := range arr {
		if v == 0 {
			return i
		}
	}
	return -1
}

func solve(arr []int) []operation {
	cur := append([]int(nil), arr...)
	ops := make([]operation, 0, len(cur))

	for len(cur) > 1 {
		idx := findZero(cur)
		if idx == -1 {
			break
		}

		l, r := idx, idx
		if idx+1 < len(cur) && cur[idx+1] == 0 {
			r = idx + 1
		} else if idx > 0 && cur[idx-1] == 0 {
			l = idx - 1
		} else if idx+1 < len(cur) {
			r = idx + 1
		} else {
			l = idx - 1
		}

		ops = append(ops, operation{l + 1, r + 1})
		cur = applyOperation(cur, l, r)
	}

	if len(cur) > 1 {
		ops = append(ops, operation{1, len(cur)})
		mex := mexSegment(cur)
		cur = []int{mex}
	}

	return ops
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}
		ops := solve(arr)
		fmt.Fprintln(out, len(ops))
		for _, op := range ops {
			fmt.Fprintf(out, "%d %d\n", op.l, op.r)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const offset = 200000
const base = 400001

func key(x, y int) int64 {
	return int64(x+offset)*base + int64(y+offset)
}

func visitedOutside(arr []int, l, r int) bool {
	if len(arr) == 0 {
		return false
	}
	i := sort.SearchInts(arr, l)
	if i > 0 {
		return true
	}
	j := sort.SearchInts(arr, r)
	return j < len(arr)
}

func visitedInRange(arr []int, L, R int) bool {
	if len(arr) == 0 {
		return false
	}
	i := sort.SearchInts(arr, L)
	return i < len(arr) && arr[i] <= R
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	preX := make([]int, n+1)
	preY := make([]int, n+1)
	posMap := make(map[int64][]int)
	posMap[key(0, 0)] = []int{0}

	for i := 1; i <= n; i++ {
		preX[i] = preX[i-1]
		preY[i] = preY[i-1]
		switch s[i-1] {
		case 'U':
			preY[i]++
		case 'D':
			preY[i]--
		case 'L':
			preX[i]--
		case 'R':
			preX[i]++
		}
		k := key(preX[i], preY[i])
		posMap[k] = append(posMap[k], i)
	}

	for ; q > 0; q-- {
		var x, y, l, r int
		fmt.Fscan(in, &x, &y, &l, &r)
		arr := posMap[key(x, y)]
		if visitedOutside(arr, l, r) {
			fmt.Fprintln(out, "YES")
			continue
		}
		tx := preX[l-1] + preX[r] - x
		ty := preY[l-1] + preY[r] - y
		arr2 := posMap[key(tx, ty)]
		if visitedInRange(arr2, l-1, r-1) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

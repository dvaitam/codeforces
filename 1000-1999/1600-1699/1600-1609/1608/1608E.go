package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Point struct {
	x int
	y int
}

func countBetween(arr []int, l, r int) int {
	// count elements strictly between l and r
	li := sort.Search(len(arr), func(i int) bool { return arr[i] > l })
	ri := sort.Search(len(arr), func(i int) bool { return arr[i] >= r })
	return ri - li
}

func possible(t int, xs, ys [3][]int) bool {
	// try vertical and horizontal orientation
	colors := []int{0, 1, 2}
	perm := make([]int, 3)
	copy(perm, colors)
	for {
		// vertical stripes (by x)
		for _, orient := range []bool{true, false} {
			var arr [3][]int
			if orient {
				arr = xs
			} else {
				arr = ys
			}
			c1, c2, c3 := perm[0], perm[1], perm[2]
			if len(arr[c1]) < t || len(arr[c2]) < t || len(arr[c3]) < t {
				// not enough cells of some color
			} else {
				x1 := arr[c1][t-1]
				x2 := arr[c3][len(arr[c3])-t]
				if x1 < x2 {
					midCount := countBetween(arr[c2], x1, x2)
					if midCount >= t {
						return true
					}
				}
			}
		}
		if !nextPermutation(perm) {
			break
		}
	}
	return false
}

func nextPermutation(a []int) bool {
	// standard next permutation
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for k, l := i+1, n-1; k < l; k, l = k+1, l-1 {
		a[k], a[l] = a[l], a[k]
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var xs [3][]int
	var ys [3][]int
	xs[0] = make([]int, 0)
	xs[1] = make([]int, 0)
	xs[2] = make([]int, 0)
	ys[0] = make([]int, 0)
	ys[1] = make([]int, 0)
	ys[2] = make([]int, 0)
	for i := 0; i < n; i++ {
		var x, y, c int
		fmt.Fscan(in, &x, &y, &c)
		c--
		xs[c] = append(xs[c], x)
		ys[c] = append(ys[c], y)
	}
	for i := 0; i < 3; i++ {
		sort.Ints(xs[i])
		sort.Ints(ys[i])
	}
	low, high := 0, n/3
	for low < high {
		mid := (low + high + 1) / 2
		if possible(mid, xs, ys) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	fmt.Println(low * 3)
}

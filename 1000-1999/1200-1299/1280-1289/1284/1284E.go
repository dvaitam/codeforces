package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Point struct{ x, y int64 }

func cross(p, a, b Point) int64 {
	return (a.x-p.x)*(b.y-p.y) - (a.y-p.y)*(b.x-p.x)
}

func C3(n int64) int64 {
	if n < 3 {
		return 0
	}
	return n * (n - 1) * (n - 2) / 6
}

func C5(n int64) int64 {
	if n < 5 {
		return 0
	}
	return n * (n - 1) * (n - 2) * (n - 3) * (n - 4) / 120
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	pts := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}

	L := make([][]int, n)
	for i := 0; i < n; i++ {
		L[i] = make([]int, n)
	}

	idx := make([]int, n-1)
	for i := 0; i < n; i++ {
		m := 0
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			idx[m] = j
			m++
		}
		arr := idx[:m]
		sort.Slice(arr, func(a, b int) bool {
			angA := math.Atan2(float64(pts[arr[a]].y-pts[i].y), float64(pts[arr[a]].x-pts[i].x))
			angB := math.Atan2(float64(pts[arr[b]].y-pts[i].y), float64(pts[arr[b]].x-pts[i].x))
			return angA < angB
		})
		arr2 := make([]int, 2*m)
		copy(arr2, arr)
		copy(arr2[m:], arr)
		k := 0
		for idxj := 0; idxj < m; idxj++ {
			if k < idxj+1 {
				k = idxj + 1
			}
			for k < idxj+m {
				cp := cross(pts[i], pts[arr[idxj]], pts[arr2[k]])
				if cp > 0 {
					k++
				} else {
					break
				}
			}
			L[i][arr[idxj]] = k - idxj - 1
		}
	}

	var sumPairs int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			l := int64(L[i][j])
			r := int64(n-2) - l
			sumPairs += C3(l) + C3(r)
		}
	}
	total5 := C5(int64(n))
	ans := 5*total5 - sumPairs
	fmt.Fprintln(out, ans)
}

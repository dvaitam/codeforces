package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func solve(r *bufio.Reader) (int, int) {
	var n int
	fmt.Fscan(r, &n)
	s := make([][]int, n)
	var xArr []int
	for i := 0; i < n; i++ {
		var k int
		fmt.Fscan(r, &k)
		s[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(r, &s[i][j])
			xArr = append(xArr, s[i][j])
		}
		sort.Ints(s[i])
	}
	sort.Ints(xArr)
	m := 0
	for i, v := range xArr {
		if i == 0 || v != xArr[i-1] {
			xArr[m] = v
			m++
		}
	}
	xArr = xArr[:m]
	for i := 0; i < n; i++ {
		for j := 0; j < len(s[i]); j++ {
			s[i][j] = sort.SearchInts(xArr, s[i][j])
		}
	}
	ord := make([]int, n)
	for i := 0; i < n; i++ {
		ord[i] = i
	}
	sort.Slice(ord, func(i, j int) bool {
		return len(s[ord[i]]) < len(s[ord[j]])
	})
	const B = 200
	got := make([]bool, m)
	for _, idxI := range ord {
		if len(s[idxI]) >= B {
			for _, y := range s[idxI] {
				got[y] = true
			}
			for _, idxJ := range ord {
				if idxJ == idxI {
					break
				}
				cnt := 0
				for _, y := range s[idxJ] {
					if got[y] {
						cnt++
						if cnt >= 2 {
							return idxI + 1, idxJ + 1
						}
					}
				}
			}
			for _, y := range s[idxI] {
				got[y] = false
			}
		}
	}
	has := make([][]int, m)
	for i := 0; i < n; i++ {
		if len(s[i]) < B {
			for _, y := range s[i] {
				has[y] = append(has[y], i)
			}
		}
	}
	got2 := make([]int, m)
	for i := range got2 {
		got2[i] = -1
	}
	for x := 0; x < m; x++ {
		if len(has[x]) > 1 {
			for _, i := range has[x] {
				for _, y := range s[i] {
					if y == x {
						continue
					}
					if got2[y] != -1 {
						return i + 1, got2[y] + 1
					}
					got2[y] = i
				}
			}
			for _, i := range has[x] {
				for _, y := range s[i] {
					got2[y] = -1
				}
			}
		}
	}
	return -1, -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		a, b := solve(reader)
		if a == -1 {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, a, b)
		}
	}
}

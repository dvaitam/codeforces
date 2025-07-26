package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Row struct {
	id      int
	vals    []int
	prefMin []int
	prefMax []int
	sufMin  []int
	sufMax  []int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		rows := make([]*Row, n)
		for i := 0; i < n; i++ {
			vals := make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &vals[j])
			}
			r := &Row{id: i, vals: vals}
			r.prefMin = make([]int, m)
			r.prefMax = make([]int, m)
			r.sufMin = make([]int, m)
			r.sufMax = make([]int, m)
			for j := 0; j < m; j++ {
				if j == 0 {
					r.prefMin[j] = vals[j]
					r.prefMax[j] = vals[j]
				} else {
					if vals[j] < r.prefMin[j-1] {
						r.prefMin[j] = vals[j]
					} else {
						r.prefMin[j] = r.prefMin[j-1]
					}
					if vals[j] > r.prefMax[j-1] {
						r.prefMax[j] = vals[j]
					} else {
						r.prefMax[j] = r.prefMax[j-1]
					}
				}
			}
			for j := m - 1; j >= 0; j-- {
				if j == m-1 {
					r.sufMin[j] = vals[j]
					r.sufMax[j] = vals[j]
				} else {
					if vals[j] < r.sufMin[j+1] {
						r.sufMin[j] = vals[j]
					} else {
						r.sufMin[j] = r.sufMin[j+1]
					}
					if vals[j] > r.sufMax[j+1] {
						r.sufMax[j] = vals[j]
					} else {
						r.sufMax[j] = r.sufMax[j+1]
					}
				}
			}
			rows[i] = r
		}

		sort.Slice(rows, func(i, j int) bool {
			return rows[i].vals[0] > rows[j].vals[0]
		})

		// build prefix and suffix aggregates
		preMinLeft := make([][]int, n)
		preMaxRight := make([][]int, n)
		sufMaxLeft := make([][]int, n)
		sufMinRight := make([][]int, n)
		for i := 0; i < n; i++ {
			preMinLeft[i] = make([]int, m)
			preMaxRight[i] = make([]int, m)
			sufMaxLeft[i] = make([]int, m)
			sufMinRight[i] = make([]int, m)
		}

		for j := 0; j < m; j++ {
			preMinLeft[0][j] = rows[0].prefMin[j]
			preMaxRight[0][j] = rows[0].sufMax[j]
		}
		for i := 1; i < n; i++ {
			for j := 0; j < m; j++ {
				if rows[i].prefMin[j] < preMinLeft[i-1][j] {
					preMinLeft[i][j] = rows[i].prefMin[j]
				} else {
					preMinLeft[i][j] = preMinLeft[i-1][j]
				}
				if rows[i].sufMax[j] > preMaxRight[i-1][j] {
					preMaxRight[i][j] = rows[i].sufMax[j]
				} else {
					preMaxRight[i][j] = preMaxRight[i-1][j]
				}
			}
		}

		for j := 0; j < m; j++ {
			sufMaxLeft[n-1][j] = rows[n-1].prefMax[j]
			sufMinRight[n-1][j] = rows[n-1].sufMin[j]
		}
		for i := n - 2; i >= 0; i-- {
			for j := 0; j < m; j++ {
				if rows[i].prefMax[j] > sufMaxLeft[i+1][j] {
					sufMaxLeft[i][j] = rows[i].prefMax[j]
				} else {
					sufMaxLeft[i][j] = sufMaxLeft[i+1][j]
				}
				if rows[i].sufMin[j] < sufMinRight[i+1][j] {
					sufMinRight[i][j] = rows[i].sufMin[j]
				} else {
					sufMinRight[i][j] = sufMinRight[i+1][j]
				}
			}
		}

		ansI, ansK := -1, -1
		for i := 0; i < n-1 && ansI == -1; i++ {
			for k := 0; k < m-1; k++ {
				if preMinLeft[i][k] > sufMaxLeft[i+1][k] && sufMinRight[i+1][k+1] > preMaxRight[i][k+1] {
					ansI = i
					ansK = k
					break
				}
			}
		}

		if ansI == -1 {
			fmt.Fprintln(out, "NO")
			continue
		}

		color := make([]byte, n)
		for idx := 0; idx <= ansI; idx++ {
			color[rows[idx].id] = 'R'
		}
		for idx := ansI + 1; idx < n; idx++ {
			color[rows[idx].id] = 'B'
		}
		fmt.Fprintln(out, "YES")
		fmt.Fprintf(out, "%s %d\n", string(color), ansK+1)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscanf(reader, "%d\n", &n); err != nil {
		return
	}

	mat := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscanf(reader, "%s\n", &mat[i])
	}

	w := make([][]int, n)
	for i := 0; i < n; i++ {
		w[i] = make([]int, n)
	}
	S := make([]int, n)

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			c := 0
			for r := 0; r < n; r++ {
				if mat[r][i] == '1' && mat[r][j] == '1' {
					c++
				}
			}
			w[i][j] = c
			w[j][i] = c
		}
	}
	
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				S[i] += w[i][j]
			}
		}
	}

	visited := make([]bool, n)
	var components [][]int

	for i := 0; i < n; i++ {
		if !visited[i] {
			comp := []int{}
			q := []int{i}
			visited[i] = true
			for len(q) > 0 {
				u := q[0]
				q = q[1:]
				comp = append(comp, u)
				for v := 0; v < n; v++ {
					if !visited[v] && w[u][v] > 0 {
						visited[v] = true
						q = append(q, v)
					}
				}
			}
			components = append(components, comp)
		}
	}

	finalOrder := []int{}

	for _, comp := range components {
		if len(comp) <= 2 {
			finalOrder = append(finalOrder, comp...)
			continue
		}

		found := false
	search:
		for _, E1 := range comp {
			minW := int(1e9)
			for _, v := range comp {
				if v != E1 && w[E1][v] < minW {
					minW = w[E1][v]
				}
			}

			candidates := []int{}
			for _, v := range comp {
				if v != E1 && w[E1][v] == minW {
					candidates = append(candidates, v)
				}
			}

			sort.Slice(candidates, func(i, j int) bool {
				return S[candidates[i]] < S[candidates[j]]
			})

			limit := 5
			if len(candidates) < limit {
				limit = len(candidates)
			}

			for _, E2 := range candidates[:limit] {
				cCopy := make([]int, len(comp))
				copy(cCopy, comp)

				sort.Slice(cCopy, func(i, j int) bool {
					u, v := cCopy[i], cCopy[j]
					d1 := w[E1][u] - w[E2][u]
					d2 := w[E1][v] - w[E2][v]
					if d1 != d2 {
						return d1 > d2
					}
					if S[u] != S[v] {
						return S[u] < S[v]
					}
					for k := 0; k < n; k++ {
						if w[u][k] != w[v][k] {
							return w[u][k] < w[v][k]
						}
					}
					return u < v
				})

				valid := true
				for i := 0; i < n; i++ {
					started := false
					ended := false
					for _, c := range cCopy {
						if mat[i][c] == '1' {
							if ended {
								valid = false
								break
							}
							started = true
						} else {
							if started {
								ended = true
							}
						}
					}
					if !valid {
						break
					}
				}

				if valid {
					finalOrder = append(finalOrder, cCopy...)
					found = true
					break search
				}
			}
		}

		if !found {
			fmt.Println("NO")
			return
		}
	}

	for i := 0; i < n; i++ {
		started := false
		ended := false
		for _, c := range finalOrder {
			if mat[i][c] == '1' {
				if ended {
					fmt.Println("NO")
					return
				}
				started = true
			} else {
				if started {
					ended = true
				}
			}
		}
	}

	fmt.Println("YES")
	for i := 0; i < n; i++ {
		for _, c := range finalOrder {
			fmt.Print(string(mat[i][c]))
		}
		fmt.Println()
	}
}
package main

import (
	"bufio"
	"fmt"
	"os"
)

// Simplex solver for problems of the form:
// maximize c^T x subject to A x <= b and x >= 0.
// Returns optimum value and solution vector.
func simplex(A [][]float64, b []float64, c []float64) ([]float64, float64) {
	m := len(A)
	n := len(c)
	// Tableau dimensions: (m constraints + 1 objective) x (n vars + m slacks + 1 rhs).
	cols := n + m
	tableau := make([][]float64, m+1)
	for i := 0; i < m+1; i++ {
		tableau[i] = make([]float64, cols+1)
	}
	// Fill constraint rows.
	for i := 0; i < m; i++ {
		copy(tableau[i], A[i])
		tableau[i][n+i] = 1 // slack
		tableau[i][cols] = b[i]
	}
	// Objective row: negative coefficients because we maximize.
	for j := 0; j < n; j++ {
		tableau[m][j] = -c[j]
	}

	// Basis: track which variable is basic in each constraint row.
	basis := make([]int, m)
	for i := 0; i < m; i++ {
		basis[i] = n + i
	}

	pivot := func(row, col int) {
		// Normalize pivot row.
		pv := tableau[row][col]
		for j := col; j <= cols; j++ {
			tableau[row][j] /= pv
		}
		// Eliminate column in other rows.
		for i := 0; i < m+1; i++ {
			if i == row {
				continue
			}
			factor := tableau[i][col]
			if factor == 0 {
				continue
			}
			for j := col; j <= cols; j++ {
				tableau[i][j] -= factor * tableau[row][j]
			}
		}
		basis[row] = col
	}

	for {
		// Choose entering variable (most positive coefficient in objective row).
		enter := -1
		for j := 0; j < cols; j++ {
			if tableau[m][j] > 1e-9 {
				enter = j
				break // Bland's rule: smallest index positive coeff.
			}
		}
		if enter == -1 {
			break // optimal
		}
		// Choose leaving row.
		bestRow := -1
		bestVal := 0.0
		for i := 0; i < m; i++ {
			if tableau[i][enter] > 1e-9 {
				ratio := tableau[i][cols] / tableau[i][enter]
				if bestRow == -1 || ratio < bestVal-1e-12 || (abs(ratio-bestVal) <= 1e-12 && basis[i] > basis[bestRow]) {
					bestRow = i
					bestVal = ratio
				}
			}
		}
		if bestRow == -1 {
			// Unbounded; should not happen, return zero solution.
			return make([]float64, n), 0
		}
		pivot(bestRow, enter)
	}

	sol := make([]float64, n)
	for i := 0; i < m; i++ {
		if basis[i] < n {
			sol[basis[i]] = tableau[i][cols]
		}
	}
	val := tableau[m][cols]
	// Since we maximized c^T x, objective value is tableau[m][cols].
	return sol, val
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	var s int
	fmt.Fscan(in, &s)
	s--

	// Precompute parent for paths using BFS from every node for neighbor-on-path queries.
	// With n <= 100, we can compute next step towards target by BFS for each source.
	nextStep := make([][]int, n)
	for i := 0; i < n; i++ {
		nextStep[i] = make([]int, n)
		for j := 0; j < n; j++ {
			nextStep[i][j] = -1
		}
		q := []int{i}
		prev := make([]int, n)
		for j := 0; j < n; j++ {
			prev[j] = -1
		}
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			for _, to := range g[v] {
				if prev[to] == -1 && to != i {
					prev[to] = v
					q = append(q, to)
				}
			}
		}
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			// walk back from j to find neighbor of i on path.
			cur := j
			for prev[cur] != i {
				cur = prev[cur]
			}
			nextStep[i][j] = cur
		}
	}

	// State indexing: prev in [0..n] where prev==n represents "no previous" for start.
	nonePrev := n
	// Values V[prev][p][b].
	V := make([][][]float64, n+1)
	for prev := 0; prev <= n; prev++ {
		V[prev] = make([][]float64, n)
		for p := 0; p < n; p++ {
			V[prev][p] = make([]float64, n)
		}
	}

	// Iterate value improvements.
	for iter := 0; iter < 2000; iter++ {
		change := 0.0
		// Precompute max over components for next states of form (prevNext, curr, neighbor).
		maxNext := make([][][]float64, n)
		for prev := 0; prev < n; prev++ {
			maxNext[prev] = make([][]float64, n)
			for curr := 0; curr < n; curr++ {
				maxNext[prev][curr] = make([]float64, n)
			}
		}
		for prev := 0; prev < n; prev++ {
			for curr := 0; curr < n; curr++ {
				if prev == curr {
					continue
				}
				for _, nb := range g[curr] {
					mx := 0.0
					for b := 0; b < n; b++ {
						if b == curr {
							continue
						}
						if nextStep[curr][b] == nb {
							val := V[prev][curr][b]
							if val > mx {
								mx = val
							}
						}
					}
					maxNext[prev][curr][nb] = mx
				}
			}
		}

		for prev := 0; prev <= n; prev++ {
			for p := 0; p < n; p++ {
				// Build LP for actions from p; if prev==p (invalid) skip.
				if prev != nonePrev && prev == p {
					continue
				}
				// actions are neighbors of p.
				acts := g[p]
				m := len(acts)
				if m == 0 {
					continue
				}
				// Prepare F matrix: actions x possible b.
				bList := make([]int, 0, n-1)
				for bi := 0; bi < n; bi++ {
					if bi == p {
						continue
					}
					bList = append(bList, bi)
				}
				F := make([][]float64, m)
				for i := 0; i < m; i++ {
					F[i] = make([]float64, len(bList))
				}
				K := 0.0
				for bj, bi := range bList {
					for ai, a := range acts {
						val := 1.0
						if a != bi {
							ns := nextStep[a][bi]
							vnext := maxNext[p][a][ns]
							val += vnext
						}
						F[ai][bj] = val
						if val > K {
							K = val
						}
					}
				}
				// Solve min max via simplex with substitution removing last action.
				var solAlpha []float64
				if m == 1 {
					solAlpha = []float64{1}
				} else {
					rows := len(bList) + 1
					cols := m - 1 + 1 // variables = alpha0..alpha_{m-2}, z'
					A := make([][]float64, rows)
					bvec := make([]float64, rows)
					// Column constraints.
					for bj := range bList {
						row := make([]float64, cols)
						for ai := 0; ai < m-1; ai++ {
							row[ai] = F[ai][bj] - F[m-1][bj]
						}
						row[m-1] = -1 // coefficient for z'
						A[bj] = row
						bvec[bj] = K - F[m-1][bj]
						if bvec[bj] < 0 {
							bvec[bj] = 0
						}
					}
					// Sum alpha <= 1 constraint.
					row := make([]float64, cols)
					for ai := 0; ai < m-1; ai++ {
						row[ai] = 1
					}
					A[len(bList)] = row
					bvec[len(bList)] = 1
					cvec := make([]float64, cols)
					cvec[m-1] = -1 // maximize -z'
					sol, _ := simplex(A, bvec, cvec)
					solAlpha = make([]float64, m)
					sumAlpha := 0.0
					for ai := 0; ai < m-1; ai++ {
						if sol[ai] < 0 {
							sol[ai] = 0
						}
						solAlpha[ai] = sol[ai]
						sumAlpha += sol[ai]
					}
					if sumAlpha >= 1 {
						sumAlpha = 1
					}
					solAlpha[m-1] = 1 - sumAlpha
				}
				// Update V for each b using computed alpha.
				for _, bi := range bList {
					val := 0.0
					for ai, a := range acts {
						prob := solAlpha[ai]
						if prob == 0 {
							continue
						}
						v := 1.0
						if a != bi {
							ns := nextStep[a][bi]
							v += maxNext[p][a][ns]
						}
						val += prob * v
					}
					diff := abs(val - V[prev][p][bi])
					if diff > change {
						change = diff
					}
					V[prev][p][bi] = val
				}
			}
		}
		if change < 1e-10 {
			break
		}
	}

	ans := 0.0
	for b := 0; b < n; b++ {
		if b == s {
			continue
		}
		if V[nonePrev][s][b] > ans {
			ans = V[nonePrev][s][b]
		}
	}
	fmt.Printf("%.10f\n", ans)
}

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var tc int
	fmt.Fscan(in, &tc)
	for ; tc > 0; tc-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}

		fmt.Fprintln(out, solveCase(a))
	}
}

func solveCase(a [][]int64) int64 {
	n := len(a)
	m := len(a[0])

	rowXor := make([]int64, n)
	colXor := make([]int64, m)

	for i := 0; i < n; i++ {
		var x int64
		for j := 0; j < m; j++ {
			x ^= a[i][j]
			colXor[j] ^= a[i][j]
		}
		rowXor[i] = x
	}

	// Build extended matrix B of size (n+1) x (m+1).
	rN := n + 1
	cM := m + 1
	B := make([][]int64, rN)
	for i := 0; i < rN; i++ {
		B[i] = make([]int64, cM)
	}
	for i := 0; i < n; i++ {
		copy(B[i], a[i])
		B[i][m] = rowXor[i]
	}
	for j := 0; j < m; j++ {
		B[n][j] = colXor[j]
	}
	var corner int64
	for _, v := range colXor {
		corner ^= v
	}
	B[n][m] = corner

	// After observing that each operation just swaps a row/column with the extra XOR row/column,
	// any reachable matrix is obtained by permuting the (n+1) rows and (m+1) columns of B
	// and then dropping one row and one column. We choose which row/column to drop and then
	// find the cheapest ordering of the remaining ones (a Hamiltonian path on rows/columns).

	// Precompute pairwise sums across columns/rows.
	rowTot := make([][]int64, rN)
	for i := 0; i < rN; i++ {
		rowTot[i] = make([]int64, rN)
		for j := i + 1; j < rN; j++ {
			var sum int64
			for c := 0; c < cM; c++ {
				diff := B[i][c] - B[j][c]
				if diff < 0 {
					diff = -diff
				}
				sum += diff
			}
			rowTot[i][j] = sum
			rowTot[j][i] = sum
		}
	}

	colTot := make([][]int64, cM)
	for i := 0; i < cM; i++ {
		colTot[i] = make([]int64, cM)
		for j := i + 1; j < cM; j++ {
			var sum int64
			for r := 0; r < rN; r++ {
				diff := B[r][i] - B[r][j]
				if diff < 0 {
					diff = -diff
				}
				sum += diff
			}
			colTot[i][j] = sum
			colTot[j][i] = sum
		}
	}

	const inf int64 = 1 << 60
	dpBuf := make([]int64, (1<<16)*16)

	// rowBest[c][r] = minimal vertical contribution when column c is excluded and row r is excluded.
	rowBest := make([][]int64, cM)
	for i := range rowBest {
		rowBest[i] = make([]int64, rN)
	}

	for exclCol := 0; exclCol < cM; exclCol++ {
		weights := buildWeightsRows(B, rowTot, exclCol)
		runTSP(weights, dpBuf, rN, inf)
		allMask := (1 << rN) - 1
		for exclRow := 0; exclRow < rN; exclRow++ {
			mask := allMask ^ (1 << exclRow)
			best := inf
			for bs := mask; bs > 0; bs &= bs - 1 {
				last := bits.TrailingZeros(uint(bs))
				val := dpBuf[mask*rN+last]
				if val < best {
					best = val
				}
			}
			rowBest[exclCol][exclRow] = best
		}
	}

	// colBest[r][c] = minimal horizontal contribution when row r is excluded and column c is excluded.
	colBest := make([][]int64, rN)
	for i := range colBest {
		colBest[i] = make([]int64, cM)
	}

	for exclRow := 0; exclRow < rN; exclRow++ {
		weights := buildWeightsCols(B, colTot, exclRow)
		runTSP(weights, dpBuf, cM, inf)
		allMask := (1 << cM) - 1
		for exclCol := 0; exclCol < cM; exclCol++ {
			mask := allMask ^ (1 << exclCol)
			best := inf
			for bs := mask; bs > 0; bs &= bs - 1 {
				last := bits.TrailingZeros(uint(bs))
				val := dpBuf[mask*cM+last]
				if val < best {
					best = val
				}
			}
			colBest[exclRow][exclCol] = best
		}
	}

	ans := inf
	for exclRow := 0; exclRow < rN; exclRow++ {
		for exclCol := 0; exclCol < cM; exclCol++ {
			total := rowBest[exclCol][exclRow] + colBest[exclRow][exclCol]
			if total < ans {
				ans = total
			}
		}
	}
	return ans
}

func buildWeightsRows(B [][]int64, rowTot [][]int64, exclCol int) [][]int64 {
	n := len(B)
	w := make([][]int64, n)
	for i := 0; i < n; i++ {
		w[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			diff := B[i][exclCol] - B[j][exclCol]
			if diff < 0 {
				diff = -diff
			}
			w[i][j] = rowTot[i][j] - diff
		}
	}
	return w
}

func buildWeightsCols(B [][]int64, colTot [][]int64, exclRow int) [][]int64 {
	m := len(B[0])
	w := make([][]int64, m)
	for i := 0; i < m; i++ {
		w[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			if i == j {
				continue
			}
			diff := B[exclRow][i] - B[exclRow][j]
			if diff < 0 {
				diff = -diff
			}
			w[i][j] = colTot[i][j] - diff
		}
	}
	return w
}

func runTSP(weights [][]int64, dp []int64, n int, inf int64) {
	maskCnt := 1 << n
	need := maskCnt * n
	for i := 0; i < need; i++ {
		dp[i] = inf
	}
	for i := 0; i < n; i++ {
		dp[(1<<i)*n+i] = 0
	}
	full := maskCnt - 1
	for mask := 1; mask <= full; mask++ {
		base := mask * n
		remaining := full ^ mask
		for bs := mask; bs > 0; bs &= bs - 1 {
			last := bits.TrailingZeros(uint(bs))
			cur := dp[base+last]
			if cur == inf {
				continue
			}
			for nb := remaining; nb > 0; nb &= nb - 1 {
				nxt := bits.TrailingZeros(uint(nb))
				nmask := mask | (1 << nxt)
				idx := nmask*n + nxt
				val := cur + weights[last][nxt]
				if val < dp[idx] {
					dp[idx] = val
				}
			}
		}
	}
}

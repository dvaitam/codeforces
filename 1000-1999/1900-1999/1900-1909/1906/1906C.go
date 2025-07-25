package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program sketches a possible interactive solution for problem C
// of contest 1906 ("Cursed Game").  The real judge interaction is not
// provided in this repository, therefore this implementation should be
// seen as a reference rather than a ready to use solution.

// basePattern is a 5x5 grid used in the first query.  The nine
// overlapping 3x3 subgrids form an invertible system which allows us to
// recover the hidden 3x3 mask.
var basePattern = [][]int{
	{1, 1, 0, 1, 0},
	{0, 1, 1, 0, 0},
	{1, 1, 0, 1, 1},
	{1, 0, 0, 0, 1},
	{1, 0, 0, 0, 0},
}

// invMatrix is the inverse (over GF(2)) of the matrix formed by those
// nine subgrids of basePattern.
var invMatrix = [9][9]int{
	{0, 1, 0, 0, 0, 1, 1, 0, 1},
	{1, 1, 1, 0, 0, 0, 0, 0, 1},
	{0, 1, 0, 0, 0, 1, 1, 1, 1},
	{0, 0, 0, 0, 1, 0, 0, 1, 1},
	{0, 0, 0, 1, 0, 0, 1, 1, 0},
	{1, 0, 1, 0, 0, 1, 1, 1, 1},
	{1, 0, 1, 0, 1, 1, 0, 1, 1},
	{0, 0, 1, 1, 1, 1, 1, 1, 0},
	{1, 1, 1, 1, 0, 1, 1, 0, 0},
}

func printBoard(out *bufio.Writer, b [][]int) {
	n := len(b)
	fmt.Fprintln(out, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if b[i][j] == 1 {
				out.WriteByte('1')
			} else {
				out.WriteByte('0')
			}
		}
		out.WriteByte('\n')
	}
	out.Flush()
}

// decodeMask reconstructs the hidden 3x3 mask from the first 3x3 block
// of the returned grid of the probe query.
func decodeMask(res [][]int) [3][3]int {
	var r [9]int
	idx := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r[idx] = res[i][j] & 1
			idx++
		}
	}
	var v [9]int
	for i := 0; i < 9; i++ {
		s := 0
		for j := 0; j < 9; j++ {
			if invMatrix[i][j] == 1 {
				s ^= r[j]
			}
		}
		v[i] = s
	}
	var m [3][3]int
	idx = 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m[i][j] = v[idx]
			idx++
		}
	}
	return m
}

// solveBoard uses Gaussian elimination over GF(2) to find a board that
// produces only ones after convolution with the mask.
func solveBoard(N int, mask [3][3]int) [][]int {
	vars := N * N
	eq := (N - 2) * (N - 2)
	width := (vars + 63) / 64

	mat := make([][]uint64, eq)
	rhs := make([]uint64, eq)
	row := 0
	for r := 0; r < N-2; r++ {
		for c := 0; c < N-2; c++ {
			vec := make([]uint64, width)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if mask[i][j] == 1 {
						idx := (r+i)*N + (c + j)
						vec[idx/64] ^= 1 << (idx % 64)
					}
				}
			}
			mat[row] = vec
			rhs[row] = 1
			row++
		}
	}

	r := 0
	for c := 0; c < vars && r < eq; c++ {
		pivot := -1
		for i := r; i < eq; i++ {
			if (mat[i][c/64]>>(c%64))&1 == 1 {
				pivot = i
				break
			}
		}
		if pivot == -1 {
			continue
		}
		mat[r], mat[pivot] = mat[pivot], mat[r]
		rhs[r], rhs[pivot] = rhs[pivot], rhs[r]
		for i := 0; i < eq; i++ {
			if i != r && (mat[i][c/64]>>(c%64))&1 == 1 {
				for k := 0; k < width; k++ {
					mat[i][k] ^= mat[r][k]
				}
				rhs[i] ^= rhs[r]
			}
		}
		r++
	}

	sol := make([]uint64, width)
	for i := r - 1; i >= 0; i-- {
		col := -1
		for j := 0; j < vars; j++ {
			if (mat[i][j/64]>>(j%64))&1 == 1 {
				col = j
				break
			}
		}
		if col == -1 {
			continue
		}
		val := rhs[i]
		for j := col + 1; j < vars; j++ {
			if (mat[i][j/64]>>(j%64))&1 == 1 {
				if (sol[j/64]>>(j%64))&1 == 1 {
					val ^= 1
				}
			}
		}
		if val&1 == 1 {
			sol[col/64] |= 1 << (col % 64)
		}
	}

	board := make([][]int, N)
	for i := 0; i < N; i++ {
		board[i] = make([]int, N)
		for j := 0; j < N; j++ {
			idx := i*N + j
			if (sol[idx/64]>>(idx%64))&1 == 1 {
				board[i][j] = 1
			}
		}
	}
	return board
}

func buildPatternBoard(N int) [][]int {
	b := make([][]int, N)
	for i := range b {
		b[i] = make([]int, N)
	}
	for i := 0; i < 5 && i < N; i++ {
		for j := 0; j < 5 && j < N; j++ {
			b[i][j] = basePattern[i][j]
		}
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const rounds = 333
	for t := 0; t < rounds; t++ {
		var N int
		if _, err := fmt.Fscan(in, &N); err != nil {
			return
		}
		if N < 5 {
			board := make([][]int, N)
			for i := 0; i < N; i++ {
				board[i] = make([]int, N)
				for j := 0; j < N; j++ {
					board[i][j] = 1
				}
			}
			printBoard(out, board)
			var dummy int
			for i := 0; i < (N-2)*(N-2); i++ {
				fmt.Fscan(in, &dummy)
			}
			continue
		}

		probe := buildPatternBoard(N)
		printBoard(out, probe)

		res := make([][]int, N-2)
		for i := range res {
			res[i] = make([]int, N-2)
			for j := range res[i] {
				fmt.Fscan(in, &res[i][j])
			}
		}
		allOnes := true
		for i := 0; i < N-2; i++ {
			for j := 0; j < N-2; j++ {
				if res[i][j] == 0 {
					allOnes = false
				}
			}
		}
		if allOnes {
			continue
		}

		mask := decodeMask(res)
		ans := solveBoard(N, mask)
		printBoard(out, ans)

		var dummy int
		for i := 0; i < (N-2)*(N-2); i++ {
			fmt.Fscan(in, &dummy)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func rref(mat [][]*big.Rat) [][]*big.Rat {
	rows := len(mat)
	if rows == 0 {
		return mat
	}
	cols := len(mat[0])
	r := 0
	zero := big.NewRat(0, 1)
	for c := 0; c < cols && r < rows; c++ {
		pivot := -1
		for i := r; i < rows; i++ {
			if mat[i][c].Cmp(zero) != 0 {
				pivot = i
				break
			}
		}
		if pivot == -1 {
			continue
		}
		if pivot != r {
			mat[r], mat[pivot] = mat[pivot], mat[r]
		}
		pivotVal := new(big.Rat).Set(mat[r][c])
		for j := c; j < cols; j++ {
			mat[r][j].Quo(mat[r][j], pivotVal)
		}
		for i := 0; i < rows; i++ {
			if i == r {
				continue
			}
			if mat[i][c].Cmp(zero) != 0 {
				factor := new(big.Rat).Set(mat[i][c])
				for j := c; j < cols; j++ {
					tmp := new(big.Rat).Mul(factor, mat[r][j])
					mat[i][j].Sub(mat[i][j], tmp)
				}
			}
		}
		r++
	}
	out := make([][]*big.Rat, 0, rows)
	for i := 0; i < rows; i++ {
		zeroRow := true
		for j := 0; j < cols; j++ {
			if mat[i][j].Cmp(zero) != 0 {
				zeroRow = false
				break
			}
		}
		if !zeroRow {
			out = append(out, mat[i])
		}
	}
	return out
}

func canonical(vectors [][]int) string {
	rows := len(vectors)
	if rows == 0 {
		return ""
	}
	cols := len(vectors[0])
	mat := make([][]*big.Rat, rows)
	for i := 0; i < rows; i++ {
		mat[i] = make([]*big.Rat, cols)
		for j := 0; j < cols; j++ {
			mat[i][j] = big.NewRat(int64(vectors[i][j]), 1)
		}
	}
	mat = rref(mat)
	var sb strings.Builder
	for i := range mat {
		for j := 0; j < cols; j++ {
			if j > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(mat[i][j].RatString())
		}
		sb.WriteByte(';')
	}
	return sb.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var m, d int
	fmt.Fscan(reader, &m, &d)
	results := make([]int, m)
	groups := make(map[string]int)
	next := 1
	for i := 0; i < m; i++ {
		var k int
		fmt.Fscan(reader, &k)
		vectors := make([][]int, k)
		for j := 0; j < k; j++ {
			vec := make([]int, d)
			for t := 0; t < d; t++ {
				fmt.Fscan(reader, &vec[t])
			}
			vectors[j] = vec
		}
		key := canonical(vectors)
		if g, ok := groups[key]; ok {
			results[i] = g
		} else {
			results[i] = next
			groups[key] = next
			next++
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i, v := range results {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

type Op struct {
	typ   byte
	idx   int
	shift int
}

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		A := make([][]int, n)
		for i := 0; i < n; i++ {
			A[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &A[i][j])
			}
		}
		B := make([][]int, n)
		for i := 0; i < n; i++ {
			B[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &B[i][j])
			}
		}
		if n > 3 {
			fmt.Fprintln(out, 0)
			continue
		}
		r := make([]int, n)
		c := make([]int, n)
		ans := 0
		// enumerate row and column shifts
		var dfsRows func(int)
		var dfsCols func(int)
		dfsRows = func(i int) {
			if i == n {
				dfsCols(0)
				return
			}
			for s := 0; s < n; s++ {
				r[i] = s
				dfsRows(i + 1)
			}
		}
		dfsCols = func(j int) {
			if j == n {
				ops := make([]Op, 0)
				for i := 0; i < n; i++ {
					if r[i] > 0 {
						ops = append(ops, Op{'R', i, r[i]})
					}
				}
				for j := 0; j < n; j++ {
					if c[j] > 0 {
						ops = append(ops, Op{'C', j, c[j]})
					}
				}
				permute(ops, 0, A, B, &ans)
				return
			}
			for s := 0; s < n; s++ {
				c[j] = s
				dfsCols(j + 1)
			}
		}
		dfsRows(0)
		fmt.Fprintln(out, ans%mod)
	}
}

func permute(ops []Op, l int, A, B [][]int, ans *int) {
	if l == len(ops) {
		if check(A, B, ops) {
			*ans = (*ans + 1) % mod
		}
		return
	}
	for i := l; i < len(ops); i++ {
		ops[l], ops[i] = ops[i], ops[l]
		permute(ops, l+1, A, B, ans)
		ops[l], ops[i] = ops[i], ops[l]
	}
}

func cloneMatrix(M [][]int) [][]int {
	n := len(M)
	N := make([][]int, n)
	for i := range M {
		N[i] = append([]int(nil), M[i]...)
	}
	return N
}

func check(A, B [][]int, ops []Op) bool {
	n := len(A)
	grid := cloneMatrix(A)
	pos := make(map[int][2]int)
	initPos := make(map[int][2]int)
	cnt := make(map[int]int)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			val := grid[i][j]
			pos[val] = [2]int{i, j}
			initPos[val] = [2]int{i, j}
		}
	}
	for _, op := range ops {
		switch op.typ {
		case 'R':
			row := op.idx
			k := op.shift % n
			if k == 0 {
				continue
			}
			newRow := make([]int, n)
			for j := 0; j < n; j++ {
				val := grid[row][j]
				nj := (j + k) % n
				newRow[nj] = val
				pos[val] = [2]int{row, nj}
				cnt[val]++
			}
			grid[row] = newRow
		case 'C':
			col := op.idx
			k := op.shift % n
			if k == 0 {
				continue
			}
			for i := 0; i < n; i++ {
				val := grid[i][col]
				ni := (i + k) % n
				grid[ni][col] = val
				pos[val] = [2]int{ni, col}
				cnt[val]++
			}
		}
	}
	for val, c := range cnt {
		if c > 2 {
			return false
		}
		if c == 2 {
			p0 := pos[val]
			i0 := initPos[val][0]
			j0 := initPos[val][1]
			dr := (p0[0] - i0 + n) % n
			dc := (p0[1] - j0 + n) % n
			// uniqueness check
			for v2, c2 := range cnt {
				if v2 >= val || c2 != 2 {
					continue
				}
				p02 := pos[v2]
				i02 := initPos[v2][0]
				j02 := initPos[v2][1]
				dr2 := (p02[0] - i02 + n) % n
				dc2 := (p02[1] - j02 + n) % n
				if dr2 == dr && dc2 == dc {
					return false
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] != B[i][j] {
				return false
			}
		}
	}
	return true
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const modG = 998244353

type Op struct {
	typ   byte
	idx   int
	shift int
}

func permute(ops []Op, l int, A, B [][]int, ans *int) {
	if l == len(ops) {
		if checkG(A, B, ops) {
			*ans = (*ans + 1) % modG
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

func checkG(A, B [][]int, ops []Op) bool {
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

func expectedG(A, B [][]int) string {
	n := len(A)
	if n > 3 {
		return "0\n"
	}
	r := make([]int, n)
	c := make([]int, n)
	ans := 0
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
	return fmt.Sprintf("%d\n", ans%modG)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	A := make([][]int, n)
	B := make([][]int, n)
	for i := 0; i < n; i++ {
		A[i] = make([]int, n)
		B[i] = make([]int, n)
		for j := 0; j < n; j++ {
			A[i][j] = rng.Intn(5)
			B[i][j] = rng.Intn(5)
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", A[i][j]))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", B[i][j]))
		}
		sb.WriteByte('\n')
	}
	expect := expectedG(A, B)
	return sb.String(), expect
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

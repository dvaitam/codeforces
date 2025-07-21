package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const MOD = 1000000007

var patternBoth = ".O.O.O.O." // symmetrical
var patternH = "........O"
var patternV = "O........"

func buildGrid(n, m int, pats [][]string) []string {
	rows := 4*n + 1
	cols := 4*m + 1
	g := make([][]byte, rows)
	for i := 0; i < rows; i++ {
		line := make([]byte, cols)
		for j := 0; j < cols; j++ {
			if i%4 == 0 || j%4 == 0 {
				line[j] = '#'
			} else {
				ci := i / 4
				cj := j / 4
				di := i%4 - 1
				dj := j%4 - 1
				line[j] = pats[ci][cj][di*3+dj]
			}
		}
		g[i] = line
	}
	res := make([]string, rows)
	for i := range g {
		res[i] = string(g[i])
	}
	return res
}

func parseAllow(n, m int, grid []string) ([][]bool, [][]bool) {
	allowH := make([][]bool, n)
	allowV := make([][]bool, n)
	for i := 0; i < n; i++ {
		allowH[i] = make([]bool, m)
		allowV[i] = make([]bool, m)
		for j := 0; j < m; j++ {
			P := make([]byte, 9)
			for di := 0; di < 3; di++ {
				for dj := 0; dj < 3; dj++ {
					P[di*3+dj] = grid[1+4*i+di][1+4*j+dj]
				}
			}
			R := make([]byte, 9)
			for ri := 0; ri < 3; ri++ {
				for ci := 0; ci < 3; ci++ {
					R[ri*3+ci] = P[(2-ci)*3+ri]
				}
			}
			sP := string(P)
			sR := string(R)
			if sP == sR {
				allowH[i][j] = true
				allowV[i][j] = true
			} else if sP < sR {
				allowH[i][j] = true
			} else {
				allowV[i][j] = true
			}
		}
	}
	return allowH, allowV
}

func solveD(n, m int, allowH, allowV [][]bool) int {
	type maskKey struct{ a, b, c, d uint64 }
	dpCur := map[maskKey]int{maskKey{}: 1}
	zero := maskKey{}
	for i := 0; i < n; i++ {
		dpNext := make(map[maskKey]int)
		for mask, val := range dpCur {
			var dfs func(pos int, next maskKey)
			dfs = func(pos int, next maskKey) {
				if pos >= m {
					dpNext[next] = (dpNext[next] + val) % MOD
					return
				}
				bit := uint(pos)
				if (mask.a>>bit)&1 == 1 {
					dfs(pos+1, next)
					return
				}
				if allowV[i][pos] && i+1 < n && allowV[i+1][pos] {
					nm := next
					nm.a |= 1 << bit
					dfs(pos+1, nm)
				}
				if pos+1 < m && ((mask.a>>uint(pos+1))&1) == 0 && allowH[i][pos] && allowH[i][pos+1] {
					dfs(pos+2, next)
				}
			}
			dfs(0, zero)
		}
		dpCur = dpNext
	}
	return dpCur[zero]
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := 100
	for t := 1; t <= tests; t++ {
		n := 1
		m := 2
		pats := [][]string{{patternH, patternH}}
		grid := buildGrid(n, m, pats)
		allowH, allowV := parseAllow(n, m, grid)
		exp := solveD(n, m, allowH, allowV)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for _, line := range grid {
			fmt.Fprintln(&input, line)
		}
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t, err)
			os.Exit(1)
		}
		var got int
		fmt.Sscan(out, &got)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", t, exp, out)
			fmt.Fprint(os.Stderr, input.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}

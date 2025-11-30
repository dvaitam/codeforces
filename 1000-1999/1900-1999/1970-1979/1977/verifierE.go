package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	infHop = 1 << 30
)

func solveCase(mat []string) string {
	n := len(mat)
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if mat[i][j] == '1' {
				adj[j][i] = true
			}
		}
	}

	reach := make([][]bool, n)
	for i := range reach {
		reach[i] = make([]bool, n)
		copy(reach[i], adj[i])
		reach[i][i] = true
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if reach[i][k] {
				for j := 0; j < n; j++ {
					if reach[k][j] {
						reach[i][j] = true
					}
				}
			}
		}
	}

	g := make([][]int, n)
	for j := 0; j < n; j++ {
		for i := 0; i < j; i++ {
			if reach[j][i] {
				g[j] = append(g[j], i)
			}
		}
	}

	matchL := make([]int, n)
	matchR := make([]int, n)
	for i := 0; i < n; i++ {
		matchL[i] = -1
		matchR[i] = -1
	}
	dist := make([]int, n)
	bfs := func() {
		q := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if matchL[i] == -1 {
				dist[i] = 0
				q = append(q, i)
			} else {
				dist[i] = infHop
			}
		}
		head := 0
		for head < len(q) {
			v := q[head]
			head++
			for _, u := range g[v] {
				w := matchR[u]
				if w != -1 && dist[w] == infHop {
					dist[w] = dist[v] + 1
					q = append(q, w)
				}
			}
		}
	}
	var dfs func(int) bool
	dfs = func(v int) bool {
		for _, u := range g[v] {
			w := matchR[u]
			if w == -1 || (dist[w] == dist[v]+1 && dfs(w)) {
				matchL[v] = u
				matchR[u] = v
				return true
			}
		}
		dist[v] = infHop
		return false
	}

	for {
		bfs()
		flow := 0
		for i := 0; i < n; i++ {
			if matchL[i] == -1 && dfs(i) {
				flow++
			}
		}
		if flow == 0 {
			break
		}
	}

	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}
	chain := 0
	for v := 0; v < n; v++ {
		if matchR[v] == -1 {
			cur := v
			for cur != -1 {
				color[cur] = chain
				cur = matchL[cur]
			}
			chain++
		}
	}

	if chain == 1 {
		return strings.TrimSpace(strings.Repeat("0 ", n))
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if color[i] == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

type testCase struct {
	mat []string
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(testcasesE, "\n")
	var tests []testCase
	for i := 0; i < len(lines); {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		if i+n >= len(lines) {
			return nil, fmt.Errorf("unexpected EOF parsing matrix")
		}
		mat := make([]string, n)
		for r := 0; r < n; r++ {
			mat[r] = strings.TrimSpace(lines[i+1+r])
		}
		tests = append(tests, testCase{mat: mat})
		i += n + 1
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.mat)))
		sb.WriteByte('\n')
		for _, row := range tc.mat {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

const testcasesE = `4
0011
0001
0000
0000

4
0110
0000
0001
0000

4
0010
0011
0000
0000

5
01001
00110
00010
00000
00000

4
0011
0000
0001
0000

5
00001
00100
00010
00000
00000

5
00100
00011
00001
00000
00000

3
001
000
000

3
010
000
000

3
000
001
000

5
00011
00110
00010
00000
00000

5
00100
00111
00000
00001
00000

3
001
000
000

4
0001
0011
0000
0000

3
000
001
000

5
00101
00000
00010
00001
00000

5
00000
00110
00010
00001
00000

4
0011
0001
0000
0000

3
010
000
000

4
0001
0010
0000
0000

5
01111
00001
00010
00000
00000

5
01101
00010
00001
00000
00000

6
010100
000100
000110
000001
000001
000000

6
001110
001101
000110
000000
000001
000000

3
010
000
000

4
0010
0001
0001
0000

5
01000
00101
00001
00000
00000

3
000
001
000

3
001
000
000

5
01001
00101
00000
00001
00000

6
001101
001001
000000
000010
000001
000000

3
001
000
000

3
001
000
000

6
010111
000010
000111
000011
000000
000000

3
000
001
000

5
00100
00000
00011
00001
00000

4
0011
0010
0000
0000

6
011010
001100
000100
000001
000000
000000

3
001
000
000

4
0100
0010
0000
0000

5
00010
00101
00000
00001
00000

4
0001
0010
0000
0000

3
001
000
000

3
001
001
000

3
000
001
000

4
0010
0001
0000
0000

4
0010
0001
0001
0000

3
010
000
000

3
010
000
000

6
000100
001000
000110
000001
000000
000000

5
01000
00001
00011
00000
00000

4
0100
0010
0000
0000

4
0010
0000
0001
0000

5
00010
00101
00001
00001
00000

4
0100
0010
0000
0000

6
001001
001010
000100
000000
000001
000000

4
0011
0001
0000
0000

5
01010
00010
00001
00000
00000

3
000
001
000

4
0001
0010
0001
0000

6
011100
000011
000100
000000
000001
000000

4
0101
0001
0000
0000

4
0100
0000
0001
0000

6
011010
001010
000011
000000
000001
000000

5
00110
00101
00011
00001
00000

6
000100
001001
000101
000010
000000
000000

4
0100
0010
0001
0000

3
001
000
000

6
001001
001100
000011
000000
000001
000000

6
010100
000100
000000
000011
000001
000000

3
001
000
000

4
0010
0000
0001
0000

6
010101
000100
000011
000000
000001
000000

3
010
000
000

3
001
000
000

5
01100
00100
00001
00000
00000

4
0010
0011
0000
0000

6
011000
000011
000100
000001
000001
000000

3
001
000
000

6
000111
001100
000010
000001
000000
000000

3
001
000
000

3
000
001
000

4
0001
0010
0000
0000

6
010010
001001
000000
000011
000001
000000

3
001
000
000

5
00110
00000
00010
00001
00000

6
001110
001000
000100
000001
000000
000000

3
010
001
000

3
011
001
000

6
000011
001000
000101
000001
000000
000000

3
011
000
000

4
0001
0010
0000
0000

5
00110
00101
00000
00001
00000

6
000111
001101
000000
000010
000001
000000

3
000
001
000

5
01000
00011
00000
00001
00000

6
010100
000000
000100
000010
000001
000000

6
010100
001101
000101
000000
000001
000000

5
00011
00111
00000
00001
00000

6
001000
001100
000111
000001
000000
000000`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}

	allInput := buildAllInput(tests)
	allOutput, err := runCandidate(bin, allInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}
	outLines := strings.Split(strings.TrimSpace(allOutput), "\n")
	if len(outLines) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d output lines, got %d\n", len(tests), len(outLines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveCase(tc.mat)
		if strings.TrimSpace(outLines[i]) != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput matrix size: %d\nexpected: %s\ngot: %s\n", i+1, len(tc.mat), want, outLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

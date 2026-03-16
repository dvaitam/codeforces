package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCaseB struct {
	n    int
	perm []int
	mat  []string
}

const testcasesRaw = `100
2
2 1
01
10
4
1 3 2 4
0011
0000
1000
1000
1
1
0
6
6 1 5 2 3 4
001100
000101
100100
111011
000101
010110
4
4 2 1 3
0000
0010
0100
0000
1
1
0
6
2 1 6 4 5 3
000000
000110
000010
010000
011000
000000
4
3 2 4 1
0000
0000
0000
0000
5
5 3 1 2 4
00000
00100
01000
00001
00010
5
3 4 2 1 5
00001
00001
00000
00000
11000
3
3 2 1
001
001
110
2
2 1
00
00
5
2 4 5 1 3
01110
10110
11001
11001
00110
5
1 5 4 3 2
00001
00001
00010
00100
11000
2
1 2
00
00
1
1
0
2
2 1
00
00
1
1
0
6
6 1 3 5 4 2
000000
001000
010001
000011
000101
001110
1
1
0
3
3 1 2
001
001
110
1
1
0
5
5 3 1 4 2
00010
00011
00001
11000
01100
3
1 2 3
001
001
110
3
2 1 3
001
000
100
4
1 2 4 3
0001
0011
0101
1110
5
1 4 5 3 2
00010
00000
00010
10101
00010
1
1
0
4
3 4 2 1
0000
0000
0000
0000
2
2 1
00
00
3
3 1 2
001
000
100
3
2 3 1
011
100
100
1
1
0
2
2 1
00
00
5
3 5 1 2 4
01101
10000
10001
00001
10110
4
2 4 1 3
0011
0001
1000
1100
4
1 3 2 4
0100
1001
0001
0110
5
1 2 3 4 5
01001
10110
01000
01000
10000
2
1 2
00
00
1
1
0
3
3 2 1
000
000
000
2
1 2
01
10
5
2 3 4 1 5
01110
10110
11010
11100
00000
2
2 1
00
00
3
2 3 1
011
101
110
3
3 1 2
000
000
000
2
2 1
01
10
6
5 4 6 2 3 1
001001
001000
110000
000000
000001
100010
1
1
0
6
3 2 6 1 5 4
000101
000110
000100
111000
010000
100000
2
2 1
00
00
6
1 3 5 6 2 4
001000
000001
100110
001010
001100
010000
1
1
0
1
1
0
5
2 3 1 4 5
01100
10100
11000
00001
00010
6
3 4 2 6 5 1
010010
100000
000000
000000
100001
000010
3
3 1 2
000
001
010
4
3 2 1 4
0100
1000
0000
0000
1
1
0
6
6 4 1 3 2 5
010101
100000
000000
100000
000000
100000
3
2 3 1
010
100
000
6
4 5 3 1 2 6
000000
001000
010000
000000
000000
000000
2
2 1
00
00
6
1 3 4 2 6 5
000000
000100
000000
010000
000000
000000
6
6 5 3 2 4 1
000000
000000
000001
000000
000000
001000
2
1 2
00
00
4
3 1 4 2
0110
1011
1100
0100
2
2 1
00
00
4
1 2 4 3
0110
1011
1100
0100
5
2 1 5 4 3
01000
10101
01001
00000
01100
5
5 1 4 2 3
00100
00000
10001
00000
00100
2
1 2
00
00
6
4 1 6 3 5 2
011000
101000
110100
001011
000101
000110
5
4 5 3 2 1
00000
00001
00001
00001
01110
2
1 2
00
00
2
2 1
01
10
3
1 3 2
010
100
000
5
4 3 5 2 1
01100
10011
10010
01100
01000
1
1
0
1
1
0
4
4 1 3 2
0011
0000
1001
1010
5
3 4 5 2 1
00010
00110
01010
11100
00000
3
3 1 2
010
100
000
2
1 2
00
00
6
4 3 2 1 6 5
000001
001001
010111
001010
001101
111010
6
6 2 3 4 5 1
010110
100000
000011
100000
101000
001000
1
1
0
6
1 5 4 2 6 3
000001
000000
000001
000000
000000
101000
6
3 4 2 1 6 5
010001
101000
010001
000000
000001
101010
1
1
0
2
1 2
00
00
3
3 1 2
001
000
100
3
1 2 3
000
000
000
1
1
0
2
2 1
00
00
5
1 4 5 2 3
00000
00000
00011
00100
00100
5
3 5 4 1 2
00000
00001
00010
00100
01000
1
1
0
2
2 1
00
00
4
4 2 3 1
0000
0011
0100
0100`

func parseTestcases() ([]testCaseB, error) {
	in := bufio.NewReader(strings.NewReader(testcasesRaw))
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseB, T)
	for i := 0; i < T; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		perm := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &perm[j])
		}
		mat := make([]string, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &mat[j])
		}
		cases[i] = testCaseB{n: n, perm: perm, mat: mat}
	}
	return cases, nil
}

func solveCase(tc testCaseB) string {
	n := tc.n
	p := append([]int(nil), tc.perm...)
	adj := tc.mat
	visited := make([]bool, n)
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		queue := []int{i}
		visited[i] = true
		comp := []int{i}
		for q := 0; q < len(queue); q++ {
			u := queue[q]
			for v := 0; v < n; v++ {
				if adj[u][v] == '1' && !visited[v] {
					visited[v] = true
					queue = append(queue, v)
					comp = append(comp, v)
				}
			}
		}
		vals := make([]int, len(comp))
		for j, idx := range comp {
			vals[j] = p[idx]
		}
		sort.Ints(comp)
		sort.Ints(vals)
		for j, idx := range comp {
			p[idx] = vals[j]
		}
	}
	var sb strings.Builder
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.perm {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, row := range tc.mat {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		expected := solveCase(tc)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

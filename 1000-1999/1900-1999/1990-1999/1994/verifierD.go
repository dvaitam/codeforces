package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesD = `3 19 18 5
4 20 16 19 3
6 1 16 9 18 8 7
5 18 18 16 13 5
3 5 17 13
2 3 6
6 2 10 1 9 16 20
5 14 13 19 15 5
4 4 2 5 16
3 9 14 10
5 17 13 19 12 18
6 14 19 8 11 1 9
6 6 11 18 19 19 4
3 19 9 10
2 3 16
5 3 12 3 14 5
2 10 14
5 4 2 20 20 2
5 19 11 18 9 17
3 2 10 1
2 4 20
6 2 7 14 10 20 9
3 2 11 11
4 5 13 13 15
6 13 20 18 4 20 17
4 14 8 10 14
4 17 10 18 11
2 14 19
4 1 13 20 19
3 2 11 15
4 12 20 9 16
2 19 2
2 12 9
5 10 19 20 11 6
4 6 11 12 20
4 10 13 4 1
6 5 10 17 8 9 8
4 6 14 4 4
6 11 11 8 15 6 3
4 7 19 15 9
3 4 2 17
3 11 19 6
4 11 3 20 12
6 5 14 10 17 9 15
4 14 10 14 19
5 2 14 5 7 1
5 20 17 14 18 8
2 15 17
4 18 11 8 3
6 10 4 8 2 2 17
3 14 19 2
2 16 4
3 17 10 8
2 17 18
5 2 20 4 11 5
4 18 16 2 12
3 7 4 18
2 6 8
4 5 1 16 19
5 2 9 8 9 20
6 17 14 2 16 11 1
2 5 2
2 2 3
5 2 3 17 17 16
4 6 11 3 12
5 13 19 10 12 9
3 11 14 4
3 18 1 13
2 19 6
2 12 15
6 18 13 2 20 14 2
4 16 11 14 14
5 1 8 7 18 9
6 3 14 8 14 5 1
4 12 18 9 4
5 4 17 13 4 11
6 18 4 19 1 16 5
3 13 2 17
2 19 4
5 6 1 11 4 1
2 16 10
6 10 3 2 19 17 17
3 4 18 4
6 2 18 11 19 6 3
3 6 8 15
6 13 9 12 20 13 12
6 14 3 13 17 8 14
3 14 19 19
6 16 5 13 5 6 4
5 16 17 15 19 6
3 9 7 5
6 17 11 8 18 10 14
6 19 19 9 7 10 1
4 16 13 7 6
6 12 8 11 16 5 14
5 20 7 15 19 18
2 16 3
5 2 15 8 8 3
3 9 8 7
4 5 6 20 2
`

type dsu struct {
	parent []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &dsu{parent: p}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra != rb {
		d.parent[rb] = ra
	}
}

func solveCase(arr []int) [][2]int {
	n := len(arr)
	uf := newDSU(n)
	edges := make([][2]int, 0, n-1)
	for t := n - 1; t >= 1; t-- {
		bucket := make([]int, t)
		for i := range bucket {
			bucket[i] = -1
		}
		for i := 0; i < n; i++ {
			if uf.find(i) == i {
				r := arr[i] % t
				if bucket[r] == -1 {
					bucket[r] = i
				} else {
					u := bucket[r]
					v := i
					edges = append(edges, [2]int{u, v})
					uf.union(u, v)
					break
				}
			}
		}
	}
	for i, j := 0, len(edges)-1; i < j; i, j = i+1, j-1 {
		edges[i], edges[j] = edges[j], edges[i]
	}
	return edges
}

type testCase struct {
	arr []int
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesD), "\n")
	tests := make([]testCase, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			return nil, fmt.Errorf("bad test line %d", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d length mismatch", i+1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[1+j])
			if err != nil {
				return nil, err
			}
			arr[j] = val
		}
		tests[i] = testCase{arr: arr}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		for _, v := range tc.arr {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(output)
	pos := 0
	for i, tc := range tests {
		if pos >= len(outFields) || strings.ToLower(outFields[pos]) != "yes" {
			fmt.Printf("case %d missing YES\n", i+1)
			os.Exit(1)
		}
		pos++
		n := len(tc.arr)
		uf := newDSU(n)
		for x := 1; x <= n-1; x++ {
			if pos+1 >= len(outFields) {
				fmt.Printf("case %d missing edges\n", i+1)
				os.Exit(1)
			}
			u, _ := strconv.Atoi(outFields[pos])
			v, _ := strconv.Atoi(outFields[pos+1])
			pos += 2
			if u < 1 || u > n || v < 1 || v > n || u == v {
				fmt.Printf("case %d invalid edge %d %d\n", i+1, u, v)
				os.Exit(1)
			}
			diff := tc.arr[u-1] - tc.arr[v-1]
			if diff < 0 {
				diff = -diff
			}
			if diff%x != 0 {
				fmt.Printf("case %d condition failed: diff %d not divisible by %d\n", i+1, diff, x)
				os.Exit(1)
			}
			uf.union(u-1, v-1)
		}
		comps := 0
		for j := 0; j < n; j++ {
			if uf.find(j) == j {
				comps++
			}
		}
		if comps > 1 {
			fmt.Printf("case %d not connected\n", i+1)
			os.Exit(1)
		}
	}
	if pos != len(outFields) {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

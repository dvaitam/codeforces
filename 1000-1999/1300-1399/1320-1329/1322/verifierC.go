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

// Embedded testcases from testcasesC.txt to avoid external dependency.
const testcasesRaw = `4 3 17 15 16 6 2 4 3 4 4 3
1 1 8 1 1
3 7 14 1 15 1 2 2 1 3 1 1 1 3 3 2 2 3 2
3 8 4 2 5 1 2 2 1 3 1 1 1 2 3 3 3 3 2 1 3
4 10 18 7 1 2 4 4 2 4 1 2 2 1 3 4 4 3 2 3 3 3 2 2 4 1
1 0 13
2 0 14 19
5 25 14 18 7 1 5 3 4 4 3 3 1 5 4 5 1 2 2 2 5 1 3 4 2 4 5 3 3 5 3 2 4 1 2 2 1 1 5 3 2 4 1 3 5 5 2 4 4 5 5 1 1 1 4 2 3
4 5 13 17 13 7 4 4 2 4 1 2 4 3 1 1
4 11 13 9 3 6 4 4 2 4 1 2 3 4 4 3 1 1 1 4 2 3 3 3 2 2 3 2
3 6 10 15 17 1 2 2 1 1 1 2 3 3 3 1 3
1 0 20
2 2 8 13 2 1 2 2
5 13 4 1 6 1 20 2 4 5 5 2 1 5 1 4 2 1 4 4 5 3 3 5 3 3 2 2 5 1 3 5 2
2 2 9 16 2 1 2 2
4 5 4 17 9 18 4 4 2 4 3 1 1 1 4 2
4 8 9 4 12 17 4 4 2 4 2 1 1 1 1 4 4 2 3 2 4 1
5 23 10 14 13 18 2 3 4 4 3 5 4 5 1 2 2 2 5 1 3 4 2 4 5 3 3 5 3 2 4 1 2 2 1 1 5 3 2 4 1 3 5 5 2 4 4 5 5 1 1 2 3
1 1 20 1 1
3 9 11 14 19 1 2 2 1 3 1 1 1 2 3 3 3 2 2 3 2 1 3
2 4 20 4 1 1 1 2 2 1 2 2
2 2 13 9 1 2 2 2
1 0 7
2 4 6 18 1 1 1 2 2 1 2 2
2 2 20 18 2 1 2 2
1 0 3
5 0 19 17 11 20 4
5 7 1 2 4 10 3 5 5 3 4 5 4 2 3 2 2 3 2 2 5
3 2 20 9 7 1 1 3 2
2 1 19 1 2 2
2 2 3 20 1 1 2 1
2 4 1 14 1 1 1 2 2 1 2 2
2 2 7 4 2 1 2 2
3 8 13 3 6 2 1 3 1 1 1 2 3 3 3 2 2 3 2 1 3
2 3 7 13 1 1 1 2 2 1
4 2 18 7 17 7 2 4 1 2
1 0 20
2 4 19 17 1 1 1 2 2 1 2 2
4 6 17 7 1 5 2 4 4 3 3 1 1 4 3 3 4 1
3 9 20 3 8 1 2 2 1 3 1 1 1 2 3 3 3 2 2 3 2 1 3
2 2 17 14 1 1 2 2
4 2 10 5 17 7 2 2 4 3
5 13 9 7 4 6 5 2 4 5 5 2 1 1 5 1 1 5 1 4 2 1 4 4 5 2 2 1 3 3 5 5 2
2 4 20 6 1 1 1 2 2 1 2 2
1 1 11 1 1
4 0 10 8 7 3
5 23 5 8 17 7 15 3 4 4 3 3 1 5 4 5 1 2 2 1 3 4 2 4 5 3 3 2 4 1 2 2 1 1 5 3 2 4 1 3 5 5 2 4 4 5 5 1 1 1 4 2 3
5 21 11 1 5 1 3 4 3 3 1 5 4 5 1 2 2 1 3 4 2 4 5 3 3 5 3 2 4 2 1 3 2 4 1 3 5 5 2 4 4 5 5 1 1 1 4 2 3
3 5 9 6 6 1 2 3 1 1 1 2 3 3 2
4 16 6 7 13 20 4 4 1 3 2 4 1 2 2 1 3 4 4 3 3 1 1 1 1 4 4 2 2 3 3 3 2 2 3 2 4 1
2 3 10 18 1 1 1 2 2 1
1 0 20
1 0 20
2 3 3 19 1 1 1 2 2 1
1 1 15 1 1
2 2 14 17 1 1 1 2
3 8 11 18 19 1 2 2 1 1 1 2 3 3 3 2 2 3 2 1 3
4 4 9 9 5 2 2 3 3 2 1 2 1 4
3 6 4 9 13 2 1 3 1 2 3 2 2 3 2 1 3
5 20 4 6 16 7 19 4 3 3 1 5 4 2 2 2 5 1 3 4 2 4 5 3 3 2 4 1 2 2 1 1 5 3 2 3 5 5 2 4 4 5 5 1 4 2 3
5 23 14 16 15 4 9 3 4 4 3 3 1 5 4 5 1 2 5 1 3 4 2 4 5 3 3 5 3 2 4 1 2 2 1 1 5 3 2 4 1 5 2 4 4 5 5 1 1 1 4 2 3
3 2 6 18 6 1 3 2 1
5 8 4 2 13 3 13 4 4 2 1 1 5 1 1 5 1 4 5 2 2 1 3
2 2 13 6 2 1 2 2
4 8 6 6 3 7 2 4 1 2 3 4 4 1 4 3 2 3 2 2 1 3
4 13 4 12 8 4 4 4 1 2 3 4 2 1 4 3 4 1 1 1 4 2 2 3 3 3 2 2 3 2 1 3
5 23 8 19 1 15 1 3 4 3 1 5 4 2 2 2 5 1 3 4 2 4 5 3 3 5 3 2 4 1 2 2 1 1 5 3 2 4 1 3 5 5 2 4 4 5 5 1 1 1 4 2 3
4 1 12 11 6 3 3 4
1 0 5
5 18 11 15 3 13 9 4 4 1 3 2 4 1 2 2 1 4 3 1 5 1 1 5 4 1 4 4 5 3 3 5 3 3 2 2 5 4 1 3 5 5 2
3 3 19 12 10 1 1 3 3 2 1
1 1 13 1 1
3 9 6 13 1 1 2 2 1 3 1 1 1 2 3 3 3 2 2 3 2 1 3
4 9 1 6 19 15 4 4 2 4 2 1 3 4 1 1 4 2 2 3 3 2 4 1
5 7 2 14 16 17 2 2 1 5 4 4 2 2 3 3 3 2 2 5 2
3 0 13 4 19
4 5 13 8 14 13 4 4 2 4 2 1 3 1 2 2
5 3 5 6 10 1 11 5 2 1 2 5 5
5 13 9 1 12 5 2 1 2 3 4 4 3 3 1 4 1 5 1 2 3 4 5 3 3 2 2 5 3 1 3 5 2
5 13 4 9 7 17 10 2 4 4 3 3 1 1 1 5 4 1 4 2 3 2 2 5 3 3 2 2 5 1 3 5 2
3 8 5 3 15 1 2 2 1 3 1 1 1 2 3 3 3 2 2 3 2
1 0 18
5 4 10 12 16 17 8 2 3 5 5 2 1 4 2
3 2 4 17 11 3 1 1 3
5 3 16 17 20 16 10 1 4 1 2 3 5
3 8 11 8 19 1 2 2 1 3 1 2 3 3 3 2 2 3 2 1 3
3 3 12 12 12 1 1 3 2 2 3
5 17 10 13 20 1 18 1 2 5 5 2 1 3 4 4 3 3 1 4 1 1 1 1 5 5 1 1 4 2 3 2 2 5 3 2 5 1 3 3 5
5 23 20 14 15 15 13 3 4 4 3 3 1 5 4 5 1 2 2 2 5 1 3 4 2 3 3 5 3 2 4 1 2 1 5 3 2 4 1 3 5 5 2 4 4 5 5 1 1 1 4 2 3
2 0 2 13
4 12 20 2 11 1 2 4 1 2 3 4 4 1 3 1 4 3 1 1 1 4 3 3 2 2 3 2 1 3
5 15 18 16 14 6 15 4 4 2 1 4 1 1 5 1 1 5 1 1 4 2 3 3 3 2 2 3 2 2 5 1 3 3 5 5 2
5 8 5 7 17 6 20 2 4 5 5 3 4 1 5 4 3 2 3 2 5 1 3
1 1 5 1 1
3 7 3 12 12 1 2 2 1 3 1 1 1 2 3 3 3 3 2
2 0 15 14
5 25 4 2 13 10 20 3 4 4 3 3 1 5 4 5 1 2 2 2 5 1 3 4 2 4 5 3 3 5 3 2 4 1 2 2 1 1 5 3 2 4 1 3 5 5 2 4 4 5 5 1 1 1 4 2 3
4 3 2 4 11 11 3 1 1 2 1 3
1 1 16 1 1
2 1 8 5 1 1`

type testCase struct {
	n     int
	edges [][2]int
	vals  []int64
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

// referenceSolution embeds the logic from 1322C.go so no oracle binary is needed.
func referenceSolution(tc testCase) int64 {
	neigh := make([][]int, tc.n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		neigh[v] = append(neigh[v], u)
	}
	mp := make(map[string]int64)
	for v := 1; v <= tc.n; v++ {
		if len(neigh[v]) == 0 {
			continue
		}
		sort.Ints(neigh[v])
		var sb strings.Builder
		for _, u := range neigh[v] {
			sb.WriteString(strconv.Itoa(u))
			sb.WriteByte(',')
		}
		key := sb.String()
		mp[key] += tc.vals[v-1]
	}
	var g int64
	for _, val := range mp {
		if g == 0 {
			g = val
		} else {
			g = gcd(g, val)
		}
	}
	return g
}

func parseTestcases() []testCase {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Buffer(make([]byte, 1024), 1<<20)
	cases := make([]testCase, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			panic("invalid testcase line")
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			panic("invalid n")
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			panic("invalid m")
		}
		need := 2 + n + 2*m
		if len(parts) != need {
			panic("testcase length mismatch")
		}
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[2+i])
			if err != nil {
				panic("invalid value")
			}
			vals[i] = int64(v)
		}
		edges := make([][2]int, m)
		idx := 2 + n
		for i := 0; i < m; i++ {
			u, err1 := strconv.Atoi(parts[idx])
			v, err2 := strconv.Atoi(parts[idx+1])
			if err1 != nil || err2 != nil {
				panic("invalid edge")
			}
			edges[i] = [2]int{u, v}
			idx += 2
		}
		cases = append(cases, testCase{n: n, edges: edges, vals: vals})
	}
	return cases
}

func run(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.vals[i]))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := parseTestcases()
	for i, tc := range tests {
		want := referenceSolution(tc)
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output %q\n", i+1, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

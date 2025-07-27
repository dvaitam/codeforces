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

func expected(tc testCase) int64 {
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
			sb.WriteString(fmt.Sprintf("%d,", u))
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
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) < 2+n*1+m*2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[2+i])
			vals[i] = int64(v)
		}
		edges := make([][2]int, m)
		pos := 2 + n
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(parts[pos])
			v, _ := strconv.Atoi(parts[pos+1])
			pos += 2
			edges[i] = [2]int{u, v}
		}
		tc := testCase{n: n, edges: edges, vals: vals}
		want := expected(tc)
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(gotStr, &got)
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

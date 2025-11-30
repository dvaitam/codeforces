package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesE = `2
1 0
1 1
3
0 1 0
1 1 0
1 0 1
3
1 0 0
1 1 1
1 0 0
2
1 0
0 0
3
0 1 0
1 1 0
1 1 0
2
1 0
1 0
4
1 0 1 0
1 0 1 0
0 1 1 0
0 0 0 0
2
1 0
0 1
4
1 1 1 0
0 1 1 1
1 1 1 1
0 0 0 0
2
1 1
0 1
4
0 0 1 0
0 1 1 0
0 0 1 1
1 0 0 1
2
0 1
1 1
2
0 0
0 1
2
1 1
0 1
3
1 1 1
1 0 1
0 0 0
3
1 0 1
1 1 1
0 0 1
4
1 0 0 0
1 0 0 1
1 1 0 1
1 1 0 0
3
0 0 1
0 1 1
0 1 0
4
1 0 1 0
0 0 1 0
0 1 0 0
0 0 1 0
4
1 1 0 1
1 0 0 1
1 0 1 0
0 1 1 1
2
1 0
1 0
2
1 0
1 1
2
1 1
0 0
4
1 1 0 1
0 1 1 0
0 0 1 0
1 1 1 0
4
0 0 1 0
1 1 1 0
1 1 1 1
0 1 1 0
4
1 0 1 1
0 0 0 1
0 0 0 1
0 1 1 0
4
0 1 1 1
1 1 1 0
0 1 0 0
0 0 1 0
2
1 0
0 0
4
1 0 0 0
0 0 1 0
1 0 0 1
0 1 1 0
3
1 0 0
0 0 0
0 0 1
3
1 1 1
1 1 0
0 0 0
2
1 0
1 0
3
0 1 0
0 1 1
1 1 1
2
0 1
0 1
3
1 1 1
0 1 1
1 1 1
3
1 1 0
0 1 0
1 0 0
4
1 1 1 1
1 1 0 0
1 1 0 0
0 0 1 0
2
0 0
1 1
2
0 0
0 1
2
1 1
1 1
2
1 0
0 0
2
1 1
0 1
2
1 0
0 1
4
0 1 0 0
0 1 0 0
0 1 1 1
0 0 1 0
3
1 1 1
1 1 0
1 0 0
4
1 1 1 1
0 1 0 0
0 0 1 1
0 1 1 1
2
1 0
0 0
4
0 0 0 1
1 1 0 1
0 0 1 0
0 0 1 0
2
1 0
1 1
3
0 1 0
0 0 0
0 1 1
3
1 0 1
0 0 0
0 1 0
2
1 0
1 0
4
1 0 1 0
0 1 0 1
1 1 1 1
0 0 1 1
3
1 1 1
0 0 0
1 1 0
2
0 0
0 1
4
0 0 1 1
1 1 1 1
1 1 0 0
1 1 1 0
2
0 1
0 1
3
0 0 0
0 1 0
1 0 1
3
1 1 1
0 0 1
1 0 0
4
1 0 1 1
1 0 0 1
0 1 1 1
1 0 1 1
2
1 1
1 0
2
1 0
1 1
2
1 0
1 0
2
1 0
0 1
3
0 0 1
0 0 1
0 1 1
3
1 1 1
1 0 0
0 0 0
4
1 0 1 1
1 1 1 1
0 0 0 1
0 1 1 1
4
1 1 1 0
0 0 0 0
1 0 0 0
0 0 1 0
4
1 0 0 1
1 0 0 0
0 0 0 1
1 1 0 1
2
1 0
1 0
2
1 1
1 0
4
1 0 1 0
0 0 0 0
1 0 0 1
0 1 0 0
3
1 1 0
1 0 0
1 1 1
3
1 0 1
1 0 0
1 0 0
4
0 1 0 1
1 0 0 0
1 1 1 1
0 1 1 1
4
0 0 0 1
1 1 0 1
0 1 0 0
1 0 1 0
2
1 0
0 0
4
1 1 0 0
0 0 1 0
0 0 0 0
1 1 1 1
4
0 1 1 0
0 0 1 0
0 0 0 0
1 1 0 1
4
1 1 1 1
0 0 1 1
0 0 0 1
0 0 1 1
2
1 0
0 1
3
0 1 0
0 1 1
1 1 0
2
1 0
1 0
4
0 1 1 0
0 1 1 0
0 0 0 1
1 0 0 1
3
0 1 0
0 1 0
0 0 0
3
1 1 1
0 0 0
0 1 0
2
1 0
0 0
4
0 0 0 1
0 1 1 1
1 0 0 0
1 0 1 1
3
1 1 0
1 0 0
0 0 0
4
0 1 1 0
1 1 0 1
1 0 1 0
1 1 0 1
3
0 1 0
1 0 1
1 0 1
4
1 1 0 0
0 1 0 1
1 0 1 1
1 1 0 1
3
0 0 0
1 1 0
0 1 1
3
0 1 1
0 0 0
1 1 1
2
1 1
1 1
3
1 0 1
0 0 0
1 1 0
3
0 1 0
1 0 1
1 0 1
2
0 0
0 1
3
1 0 0
1 1 1
0 1 1
2
0 0
0 1`

func solve(n int, mat [][]int) string {
	vis := make([]bool, n)
	q := make([]int, 0, n)
	vis[0] = true
	q = append(q, 0)
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for v := 0; v < n; v++ {
			if mat[u][v] == 1 && !vis[v] {
				vis[v] = true
				q = append(q, v)
			}
		}
	}
	for i := 0; i < n; i++ {
		if !vis[i] {
			return "NO"
		}
	}
	for i := range vis {
		vis[i] = false
	}
	q = q[:0]
	vis[0] = true
	q = append(q, 0)
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for v := 0; v < n; v++ {
			if mat[v][u] == 1 && !vis[v] {
				vis[v] = true
				q = append(q, v)
			}
		}
	}
	for i := 0; i < n; i++ {
		if !vis[i] {
			return "NO"
		}
	}
	return "YES"
}

type testCase struct {
	input    string
	expected string
}

func parseCases(data string) ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(data))
	sc.Buffer(make([]byte, 0, 1024), 1<<20)
	sc.Split(bufio.ScanWords)
	cases := []testCase{}
	for {
		if !sc.Scan() {
			break
		}
		nVal := sc.Text()
		if strings.TrimSpace(nVal) == "" {
			continue
		}
		n, err := strconv.Atoi(nVal)
		if err != nil {
			return nil, err
		}
		mat := make([][]int, n)
		for i := 0; i < n; i++ {
			mat[i] = make([]int, n)
			for j := 0; j < n; j++ {
				if !sc.Scan() {
					return nil, fmt.Errorf("unexpected EOF")
				}
				val, err := strconv.Atoi(sc.Text())
				if err != nil {
					return nil, err
				}
				mat[i][j] = val
			}
		}
		expected := solve(n, mat)
		var input strings.Builder
		input.WriteString(strconv.Itoa(n))
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(strconv.Itoa(mat[i][j]))
			}
			input.WriteByte('\n')
		}
		cases = append(cases, testCase{input: input.String(), expected: expected})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases(testcasesE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

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

const testcasesData = `
100
2
2 1
3
1 2 3
6
4 3 2 6 1 5
5
1 4 3 5 2
5
3 1 2 4 5
2
2 1
7
2 3 6 7 1 4 5
3
3 2 1
4
1 2 3 4
4
3 2 1 4
1
1
1
1
2
1 2
2
1 2
6
1 2 5 6 3 4
5
2 1 5 4 3
3
2 1 3
7
6 1 4 7 2 5 3
6
4 1 5 3 2 6
1
1
4
3 2 1 4
7
6 3 5 7 2 4 1
7
7 3 2 5 6 1 4
5
5 4 1 2 3
1
1
1
1
1
1
5
3 2 4 1 5
3
1 2 3
2
2 1
7
1 5 4 2 6 7 3
4
2 1 3 4
5
4 1 2 3 5
3
1 3 2
5
4 2 1 5 3
5
5 4 2 1 3
1
1
6
1 5 2 4 3 6
5
2 5 1 4 3
6
6 2 4 5 3 1
3
1 2 3
2
1 2
2
1 2
7
4 1 2 6 7 5 3
7
3 6 4 2 5 1 7
6
4 5 1 6 2 3
6
4 2 3 5 1 6
3
2 1 3
4
3 4 1 2
6
3 1 5 4 2 6
7
3 6 4 2 5 7 1
7
2 5 4 6 1 7 3
4
1 2 4 3
4
1 2 4 3
3
1 3 2
1
1
4
4 3 1 2
4
2 1 3 4
1
1
6
1 2 6 3 5 4
2
2 1
7
4 6 7 2 1 3 5
7
7 1 3 4 2 5 6
1
1
4
3 2 4 1
2
2 1
5
2 3 1 4 5
3
3 2 1
7
2 6 7 3 1 4 5
1
1
5
3 4 5 2 1
7
3 7 6 5 4 1 2
2
1 2
5
3 2 1 4 5
3
2 3 1
7
3 7 4 5 6 1 2
1
1
7
3 7 2 4 5 1 6
3
3 2 1
4
1 2 3 4
3
3 1 2
3
3 1 2
2
2 1
6
2 5 3 1 4 6
1
1
3
1 3 2
6
3 2 6 4 5 1
6
1 2 5 6 3 4
4
2 3 4 1
5
5 4 2 1 3
4
3 4 1 2
3
1 2 3
1
1
4
4 2 3 1
6
2 4 5 6 3 1
4
3 4 1 2
1
1
5
4 2 3 5 1
2
2 1
3
2 3 1
`

type testCase struct {
	n    int
	perm []int
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

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	scanner.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !scanner.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	t, err := nextInt()
	if err != nil {
		return nil, err
	}
	if t < 0 {
		return nil, fmt.Errorf("negative test count")
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		perm := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := nextInt()
			if err != nil {
				return nil, err
			}
			perm[j] = val
		}
		cases = append(cases, testCase{n: n, perm: perm})
	}
	return cases, nil
}

// solve mirrors the logic of 1063E.go so the verifier is self contained.
func solve(tc testCase) string {
	n := tc.n
	goPos := make([]int, n)
	for i, a := range tc.perm {
		a--
		goPos[a] = i
	}
	dame := false
	for i := 0; i < n; i++ {
		if goPos[i] != i {
			dame = true
			break
		}
	}
	base := make([]byte, n)
	for i := range base {
		base[i] = '.'
	}

	var b strings.Builder
	if !dame {
		fmt.Fprintf(&b, "%d\n", n)
		line := string(base)
		for i := 0; i < n; i++ {
			b.WriteString(line)
			b.WriteByte('\n')
		}
		return strings.TrimSpace(b.String())
	}

	goPos[0] = -1
	fmt.Fprintf(&b, "%d\n", n-1)
	for h := 0; h < n; h++ {
		dL, dR := -1, -1
		for i := 0; i < n; i++ {
			if goPos[i] != i {
				dL = i
				break
			}
		}
		for i := n - 1; i >= 0; i-- {
			if goPos[i] != i {
				dR = i
				break
			}
		}
		gen := make([]byte, n)
		copy(gen, base)
		if dL == dR {
			// nothing to do
		} else if goPos[dL] == -1 {
			t := goPos[dR]
			gen[dL] = '/'
			gen[dR] = '/'
			gen[t] = '/'
			goPos[dL] = goPos[t]
			goPos[t] = goPos[dR]
			goPos[dR] = -1
		} else {
			t := goPos[dL]
			gen[dR] = '\\'
			gen[dL] = '\\'
			gen[t] = '\\'
			goPos[dR] = goPos[t]
			goPos[t] = goPos[dL]
			goPos[dL] = -1
		}
		b.Write(gen)
		b.WriteByte('\n')
	}
	return strings.TrimSpace(b.String())
}

func formatInput(tc testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", tc.n)
	for i, v := range tc.perm {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := formatInput(tc)
		exp := solve(tc)
		out, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

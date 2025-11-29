package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesM.txt so the verifier is self-contained.
const testcasesRaw = `9
14
6
10
20
18
17
8
10
11
8
9
18
13
9
3
11
11
17
19
11
19
14
15
13
5
2
6
13
5
4
8
10
15
17
3
15
2
13
4
16
11
8
12
16
14
17
2
10
13
5
17
12
11
8
20
9
4
8
12
10
2
9
17
15
15
12
18
6
17
6
4
18
20
10
8
13
12
11
4
15
19
4
18
15
9
10
7
6
9
7
5
9
11
10
3
19
5
15
11`

func parseTestcases() ([]int, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	res := make([]int, 0, len(lines))
	for idx, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		n, err := strconv.Atoi(ln)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", idx+1, err)
		}
		res = append(res, n)
	}
	return res, nil
}

func solve(n int) string {
	b := (n + 50 - 1) / 50
	if b <= 0 {
		b = 1
	}
	k := (n + b - 1) / b

	var rowsList [][]int
	var colsList [][]int
	for t := 0; t < k; t++ {
		start := t*b + 1
		end := (t + 1) * b
		if end > n {
			end = n
		}
		rows := make([]int, 0, end-start+1)
		for i := start; i <= end; i++ {
			rows = append(rows, i)
		}
		excluded := make([]bool, n+2)
		for _, i := range rows {
			excluded[i] = true
			if i >= 2 {
				excluded[i-1] = true
			}
		}
		cols := make([]int, 0, n)
		for j := 1; j <= n; j++ {
			if !excluded[j] {
				cols = append(cols, j)
			}
		}
		if len(cols) == 0 {
			continue
		}
		rowsList = append(rowsList, rows)
		colsList = append(colsList, cols)
	}

	var out strings.Builder
	fmt.Fprintln(&out, len(rowsList))
	for idx := range rowsList {
		rows := rowsList[idx]
		cols := colsList[idx]
		fmt.Fprint(&out, len(rows))
		for _, v := range rows {
			fmt.Fprint(&out, " ", v)
		}
		fmt.Fprintln(&out)
		fmt.Fprint(&out, len(cols))
		for _, v := range cols {
			fmt.Fprint(&out, " ", v)
		}
		fmt.Fprintln(&out)
	}
	return strings.TrimSpace(out.String())
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, n := range cases {
		input := fmt.Sprintf("%d\n", n)
		expect := solve(n)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `3
1 3
2 3

8
1 2
3 5
4 2
2 8
5 8
6 8
7 8

7
3 4
5 7
6 2
2 1
1 4
4 7

2
1 2

8
2 7
3 7
6 1
1 8
7 5
5 4
4 8

6
2 1
4 3
3 1
5 1
1 6

2
1 2

7
3 5
5 1
1 4
4 6
6 2
2 7

5
3 1
1 5
4 2
2 5

8
1 8
2 8
3 4
5 6
6 4
7 4
4 8

8
4 8
6 5
5 1
1 7
7 2
2 3
3 8

7
2 6
4 7
5 3
3 1
1 6
6 7

4
1 4
3 2
2 4

4
1 3
2 4
3 4

8
2 7
3 1
1 8
5 4
4 7
6 7
7 8

7
1 2
2 3
3 5
4 6
5 7
6 7

7
2 6
5 3
3 1
1 4
4 6
6 7

6
3 1
1 2
2 5
5 4
4 6

4
2 4
3 1
1 4

5
2 1
1 3
3 5
4 5

6
1 5
3 4
4 6
5 2
2 6

3
1 3
2 3

3
2 1
1 3

8
1 4
2 4
3 7
4 6
5 6
6 8
7 8

4
2 1
1 4
3 4

8
2 3
3 4
4 7
5 1
1 8
7 6
6 8

6
1 5
3 2
2 5
5 4
4 6

5
1 3
2 4
4 3
3 5

2
1 2

6
1 5
2 5
4 5
5 3
3 6

5
3 5
4 1
1 2
2 5

7
1 2
3 5
4 5
5 2
2 7
6 7

2
1 2

8
3 5
4 1
5 2
6 2
2 1
1 8
7 8

2
1 2

8
1 5
7 4
4 5
5 2
2 3
3 6
6 8

4
3 1
1 2
2 4

3
1 2
2 3

6
1 2
2 6
4 3
3 6
5 6

7
1 3
2 4
5 6
6 3
3 4
4 7

5
2 1
4 1
1 3
3 5

5
1 3
3 4
4 2
2 5

4
2 1
1 3
3 4

7
1 5
3 2
2 5
5 4
4 7
6 7

2
1 2

3
2 1
1 3

5
3 2
4 1
1 2
2 5

5
1 5
2 4
3 5
4 5

8
2 4
3 8
5 4
4 1
1 7
7 6
6 8

7
2 6
5 4
4 1
1 6
6 3
3 7

3
2 1
1 3

2
1 2

4
2 1
3 1
1 4

4
1 3
3 2
2 4

5
1 5
4 3
3 2
2 5

2
1 2

6
3 1
1 5
4 2
2 5
5 6

5
1 2
2 5
3 5
4 5

2
1 2

5
4 2
2 3
3 1
1 5

3
1 3
2 3

7
3 4
6 5
5 2
2 4
4 1
1 7

7
2 4
6 3
3 5
5 4
4 1
1 7

4
1 4
2 3
3 4

2
1 2

3
2 1
1 3

8
1 6
2 3
3 6
6 7
7 4
4 5
5 8

7
2 1
1 7
6 4
4 5
5 3
3 7

8
5 8
6 4
4 2
7 1
1 2
2 3
3 8

3
2 1
1 3

6
1 2
2 3
4 3
3 5
5 6

6
1 3
2 3
4 3
5 3
3 6

2
1 2

4
1 2
2 4
3 4

3
1 3
2 3

6
2 1
5 3
3 1
1 4
4 6

2
1 2

5
1 2
4 2
2 3
3 5

2
1 2

6
2 5
3 4
4 1
1 5
5 6

6
4 2
2 5
5 1
1 3
3 6

4
2 3
3 1
1 4

5
2 3
3 1
4 1
1 5

8
3 5
4 1
5 1
1 2
6 7
7 2
2 8

8
5 1
1 4
6 4
4 7
7 3
3 2
2 8

5
1 2
3 2
4 2
2 5

7
2 7
3 1
1 4
5 4
4 7
6 7

6
1 3
2 5
4 3
3 6
5 6

5
4 3
3 1
1 2
2 5

7
2 3
3 1
4 1
5 1
1 7
6 7

4
1 3
2 4
3 4

5
2 3
3 4
4 1
1 5

2
1 2

4
2 4
3 1
1 4

4
1 2
2 4
3 4

7
1 3
4 3
3 2
6 5
5 2
2 7

4
1 2
3 2
2 4

4
2 1
1 3
3 4

2
1 2

8
1 8
3 2
2 6
6 4
4 7
7 5
5 8`

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

func solve(n int, edges [][2]int) int {
	deg := make([]int, n+1)
	for _, e := range edges {
		deg[e[0]]++
		deg[e[1]]++
	}
	leaves := 0
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			leaves++
		}
	}
	return (leaves + 1) / 2
}

type testCase struct {
	n     int
	edges [][2]int
}

func parseCases(raw string) ([]testCase, error) {
	blocks := strings.Split(raw, "\n\n")
	var res []testCase
	for _, b := range blocks {
		fields := strings.Fields(b)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n: %w", err)
		}
		expected := 1 + 2*(n-1)
		if len(fields) != expected {
			return nil, fmt.Errorf("case with n=%d has %d numbers, expected %d", n, len(fields), expected)
		}
		edges := make([][2]int, 0, n-1)
		for i := 1; i+1 < len(fields); i += 2 {
			u, err1 := strconv.Atoi(fields[i])
			v, err2 := strconv.Atoi(fields[i+1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid edge numbers")
			}
			edges = append(edges, [2]int{u, v})
		}
		res = append(res, testCase{n: n, edges: edges})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases(testcases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		want := solve(tc.n, tc.edges)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		gotStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed to parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

type testCaseG struct {
	input string
}

const testcasesRaw = `100
6
1 2
2 3
2 4
1 5
1 6
2
6 5 4 6
3 3 1 3
5
1 2
2 3
3 4
1 5
2
5 5 3 5
1 4 3 1
4
1 2
2 3
2 4
1
2 3 1 1
6
1 2
2 3
2 4
2 5
5 6
5
6 6 5 1
6 3 2 5
4 3 3 5
1 1 5 6
5 2 1 5
4
1 2
1 3
2 4
4
2 2 2 1
2 2 1 3
3 3 4 4
3 2 4 1
3
1 2
2 3
1
2 1 1 3
5
1 2
2 3
1 4
4 5
4
3 4 1 2
3 4 4 2
1 1 5 1
3 5 3 2
4
1 2
2 3
3 4
3
4 4 4 3
4 4 4 2
4 4 3 4
4
1 2
1 3
3 4
4
3 3 4 4
2 3 1 4
1 2 1 1
3 3 4 1
3
1 2
2 3
5
1 2 3 2
1 2 1 1
2 2 1 3
3 1 3 2
1 1 3 3
2
1 2
1
2 2 1 1
6
1 2
1 3
1 4
3 5
4 6
4
6 1 6 2
2 4 4 4
4 5 3 4
6 5 1 3
4
1 2
2 3
1 4
1
3 3 3 4
5
1 2
2 3
1 4
1 5
2
5 2 5 1
2 3 3 1
2
1 2
4
2 1 1 1
2 1 2 2
1 1 1 2
1 1 1 1
4
1 2
2 3
1 4
2
2 4 2 3
2 1 4 2
4
1 2
2 3
1 4
4
2 2 3 3
2 2 2 2
4 1 1 1
1 4 1 1
6
1 2
2 3
3 4
1 5
1 6
5
2 4 2 5
6 3 1 4
1 6 5 1
2 6 4 1
1 6 2 1
2
1 2
3
2 2 1 1
1 2 1 2
2 1 1 2
4
1 2
1 3
1 4
3
3 4 1 3
2 1 1 3
1 4 1 4
4
1 2
2 3
1 4
2
2 4 2 2
4 4 4 3
4
1 2
2 3
1 4
2
4 1 2 1
1 2 1 2
6
1 2
2 3
2 4
4 5
2 6
4
2 4 1 3
4 3 3 2
3 3 5 6
3 2 1 3
3
1 2
2 3
5
3 1 3 1
1 2 2 1
1 3 1 3
3 2 2 2
3 3 2 2
2
1 2
2
1 1 1 1
1 1 1 1
4
1 2
1 3
3 4
2
2 2 2 2
4 4 1 2
5
1 2
2 3
3 4
2 5
3
3 1 3 1
1 1 4 4
3 1 1 3
2
1 2
5
2 1 2 1
1 2 2 2
1 1 1 2
2 1 1 2
2 1 2 1
5
1 2
1 3
2 4
3 5
3
5 3 3 5
3 4 5 1
5 2 5 1
6
1 2
2 3
1 4
1 5
3 6
1
2 5 3 2
3
1 2
1 3
3
3 2 3 2
1 3 3 1
3 1 1 2
6
1 2
2 3
1 4
4 5
2 6
4
4 5 5 3
1 6 6 1
4 1 5 5
2 5 4 3
3
1 2
2 3
5
3 2 3 2
3 2 2 2
3 2 1 3
2 3 1 2
1 3 3 2
6
1 2
1 3
3 4
1 5
1 6
4
4 5 1 2
6 3 3 1
1 3 6 3
3 4 4 5
2
1 2
3
2 2 2 2
1 1 1 2
1 1 2 1
2
1 2
2
1 1 2 2
2 1 1 2
3
1 2
1 3
1
1 1 3 2
5
1 2
2 3
2 4
3 5
5
4 2 2 5
1 2 2 4
2 5 3 4
5 3 5 5
2 5 2 1
5
1 2
2 3
3 4
1 5
1
3 3 4 3
5
1 2
2 3
2 4
3 5
5
1 4 5 2
3 3 5 4
2 1 4 5
2 3 5 4
4 3 1 4
5
1 2
1 3
1 4
1 5
4
4 1 1 5
5 3 3 2
2 4 5 1
4 5 5 1
2
1 2
2
2 2 2 2
1 2 2 2
3
1 2
2 3
1
1 3 1 1
2
1 2
4
1 2 2 2
2 1 1 2
1 1 1 1
1 2 2 1
4
1 2
2 3
1 4
3
3 2 1 2
2 4 1 4
1 1 1 1
6
1 2
1 3
3 4
4 5
3 6
2
4 2 1 2
3 2 5 5
4
1 2
1 3
1 4
3
1 4 2 2
4 4 1 4
2 2 2 3
5
1 2
2 3
3 4
1 5
4
2 2 3 2
4 3 1 3
3 3 4 4
3 5 4 3
3
1 2
1 3
5
2 3 1 3
3 1 1 1
1 3 2 2
3 1 1 3
2 1 3 1
5
1 2
1 3
2 4
3 5
2
3 5 1 5
3 2 3 3
2
1 2
1
1 1 2 1
3
1 2
2 3
1
3 2 2 1
6
1 2
2 3
3 4
1 5
1 6
2
6 5 5 2
5 2 3 2
4
1 2
2 3
2 4
3
2 2 1 2
4 1 1 1
1 3 3 4
5
1 2
1 3
2 4
4 5
3
4 4 5 5
5 2 1 1
4 5 4 1
6
1 2
2 3
3 4
2 5
1 6
1
1 2 2 4
3
1 2
2 3
2
1 1 2 1
2 2 3 3
6
1 2
1 3
2 4
2 5
1 6
1
4 6 6 6
4
1 2
1 3
2 4
4
4 2 3 2
1 4 4 1
2 3 3 1
3 2 3 3
4
1 2
2 3
2 4
3
4 3 2 1
4 3 1 2
3 1 1 1
3
1 2
2 3
1
2 1 1 1
2
1 2
3
1 1 1 1
1 1 2 2
1 1 2 1
5
1 2
2 3
2 4
1 5
5
5 4 3 5
4 5 5 5
1 3 3 2
1 2 3 1
2 1 5 3
4
1 2
1 3
2 4
4
3 1 2 3
1 4 4 1
3 3 2 1
4 3 1 2
2
1 2
1
1 2 1 2
2
1 2
3
2 2 1 1
1 2 2 1
1 1 1 2
2
1 2
5
2 2 2 1
1 1 2 1
1 1 2 1
2 2 2 1
2 1 1 2
3
1 2
1 3
2
1 1 2 3
1 2 3 3
2
1 2
4
1 1 2 1
2 2 2 1
1 2 2 1
2 1 1 2
2
1 2
1
1 1 1 1
3
1 2
2 3
2
1 1 3 3
3 3 3 3
4
1 2
1 3
3 4
4
1 3 2 2
1 4 3 2
3 1 2 4
1 3 1 1
2
1 2
1
1 2 1 2
6
1 2
2 3
3 4
4 5
2 6
1
2 1 2 1
3
1 2
2 3
5
3 3 2 1
2 1 2 3
1 2 2 3
2 1 3 3
3 2 2 2
2
1 2
5
2 1 2 2
1 1 2 2
2 1 2 2
1 2 1 1
2 1 1 2
3
1 2
1 3
4
3 2 2 2
1 1 3 1
1 1 3 3
3 2 3 1
3
1 2
2 3
3
1 3 3 2
1 1 1 1
1 2 2 2
3
1 2
2 3
1
3 3 2 3
4
1 2
2 3
1 4
1
2 1 1 4
5
1 2
2 3
3 4
2 5
5
2 1 5 4
4 2 4 5
4 5 1 5
5 4 3 2
3 2 4 2
4
1 2
2 3
3 4
5
2 4 2 3
1 3 4 3
3 2 2 3
4 1 1 1
1 4 3 4
2
1 2
3
1 1 1 2
1 1 2 2
2 1 2 2
2
1 2
5
1 2 1 1
2 1 1 2
1 2 1 1
2 2 1 1
1 1 1 2
3
1 2
1 3
5
1 2 1 3
3 1 2 1
3 1 2 3
1 3 3 1
3 2 1 3
4
1 2
2 3
1 4
3
4 2 2 1
3 3 1 4
4 2 4 4
3
1 2
1 3
5
3 2 3 1
2 1 2 1
1 3 3 2
3 2 2 2
2 1 3 1
6
1 2
2 3
2 4
3 5
3 6
2
4 5 5 4
5 4 2 3
3
1 2
2 3
5
2 1 1 1
1 2 1 1
2 1 3 2
2 3 2 1
3 2 3 1
4
1 2
2 3
3 4
1
3 2 4 1
3
1 2
2 3
2
1 3 1 3
3 1 1 2
6
1 2
2 3
2 4
3 5
1 6
3
1 2 1 5
4 4 1 2
2 5 4 1
6
1 2
1 3
1 4
1 5
5 6
4
2 1 3 1
3 6 4 6
2 6 2 2
1 5 4 2
2
1 2
4
1 1 1 1
2 1 2 2
1 2 2 1
2 1 2 1
4
1 2
2 3
3 4
3
1 4 1 4
3 4 2 4
2 1 2 4
3
1 2
2 3
5
3 1 1 2
1 3 1 1
3 2 3 2
3 3 1 3
1 1 1 3
6
1 2
2 3
3 4
1 5
1 6
5
2 3 2 5
1 1 6 5
6 2 4 3
6 6 1 3
2 2 6 1
5
1 2
2 3
3 4
2 5
1
3 1 5 1
4
1 2
1 3
2 4
5
3 4 3 3
1 2 2 3
2 1 4 4
3 2 3 1
1 3 3 2
5
1 2
2 3
1 4
1 5
5
2 5 1 3
5 2 5 5
4 4 2 2
4 4 2 5
5 2 4 5`

func parseTestcases() ([]testCaseG, error) {
	in := bufio.NewReader(strings.NewReader(testcasesRaw))
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseG, T)
	for i := 0; i < T; i++ {
		var n int
		fmt.Fscan(in, &n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n-1; j++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
		}
		var q int
		fmt.Fscan(in, &q)
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for j := 0; j < q; j++ {
			var u, v, x, y int
			fmt.Fscan(in, &u, &v, &x, &y)
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", u, v, x, y))
		}
		cases[i] = testCaseG{input: sb.String()}
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expected, err := run("500G.go", tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ref failed on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, strings.TrimSpace(expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

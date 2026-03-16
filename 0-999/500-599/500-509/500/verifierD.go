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

type testCaseD struct {
	input string
}

const testcasesRaw = `100
4
1 2 3
4 1 10
4 3 9
2
1 8
3 9
6
2 6 3
2 3 9
1 6 2
5 1 5
5 6 5
4
5 7
4 7
5 8
2 6
3
2 1 5
2 3 7
3
2 9
2 10
2 9
7
2 3 1
2 3 10
3 4 3
4 5 6
5 6 10
6 7 4
5
3 5
1 2
4 8
1 6
1 7
4
4 1 1
2 3 1
3 2 1
3
1 2
1 10
3 1
4
1 2 5
2 3 3
3 4 6
3
2 3
2 7
2 9
6
1 2 10
2 3 2
3 4 9
6 2 5
3 5 5
5
3 1
4 10
3 1
4 10
5 3
3
1 2 6
3 2 10
3
2 1
1 1
2 5
6
5 3 3
3 5 5
4 1 1
4 5 3
2 6 5
2
3 3
4 2
3
1 2 6
2 3 4
4
1 2
2 4
2 5
1 2
3
1 2 4
3 1 5
3
1 10
2 10
1 7
5
1 2 5
4 3 7
3 4 1
4 5 3
2
1 8
4 9
4
4 3 9
2 1 10
3 4 2
2
1 1
3 9
4
1 2 7
2 3 1
1 2 9
3
1 1
3 9
2 1
7
2 3 9
7 1 6
1 5 2
7 3 3
5 6 1
5 7 7
1
3 4
5
1 2 9
4 3 1
3 4 3
1 4 1
1
4 6
4
3 4 7
2 3 6
3 4 2
2
3 1
3 7
3
1 2 1
3 2 1
5
2 1
2 8
2 7
2 8
1 4
4
1 2 10
4 2 7
1 3 6
5
2 2
2 2
3 9
2 2
3 6
7
1 2 10
2 3 8
7 4 1
4 5 10
4 2 1
7 1 1
1
6 8
5
1 2 2
5 2 2
3 4 2
4 5 1
5
3 10
2 2
2 3
2 8
4 5
5
1 2 6
2 3 2
2 4 3
5 4 3
4
2 3
1 8
4 9
4 10
4
2 3 4
2 3 9
3 4 5
4
3 10
3 5
1 5
1 5
6
1 2 4
3 2 6
2 4 8
4 5 4
6 5 1
4
1 7
1 8
2 4
1 4
5
2 3 3
1 3 3
3 4 6
1 3 5
1
3 8
7
1 2 6
3 4 7
2 6 10
4 5 8
5 3 2
6 7 2
4
1 8
5 5
1 9
6 6
5
1 2 5
2 3 5
3 5 9
4 5 8
5
3 1
3 10
2 3
2 6
4 2
3
1 2 3
2 3 10
4
2 3
2 8
2 3
1 2
4
1 2 9
2 3 7
3 4 1
2
1 8
3 5
4
1 2 9
3 4 8
3 4 2
3
2 2
1 5
3 2
3
1 2 2
1 3 7
4
1 10
1 7
1 9
1 10
4
3 1 8
2 3 7
3 4 9
3
3 8
1 10
1 1
3
1 2 4
3 2 4
2
1 1
2 2
7
6 2 3
1 5 1
5 2 9
5 1 6
7 3 3
7 5 1
1
6 6
4
1 2 4
2 4 6
1 4 2
2
3 3
2 8
6
1 2 9
2 3 4
3 4 5
4 6 7
4 1 5
3
3 8
5 6
5 7
4
1 2 4
3 2 8
3 4 3
2
1 3
3 7
7
1 2 2
2 3 8
1 7 1
3 4 1
5 6 1
1 4 8
1
2 4
3
1 2 10
1 3 1
3
2 8
2 8
1 9
5
3 1 1
2 3 2
3 4 1
1 2 9
1
4 1
4
1 2 6
4 3 8
3 4 1
4
2 10
3 7
1 5
1 7
3
1 2 9
3 2 3
4
2 3
2 7
2 4
2 8
5
2 5 10
2 3 7
3 4 1
1 4 2
1
3 4
7
1 2 1
2 3 1
7 6 6
4 5 1
4 7 2
3 1 8
1
6 4
4
1 2 3
2 3 3
3 4 1
1
1 5
4
4 3 9
2 3 8
1 2 5
1
1 8
6
6 5 2
2 3 6
6 4 2
5 4 7
5 1 8
3
1 6
3 4
4 2
7
1 2 6
6 1 10
3 4 4
1 3 4
1 3 4
6 2 4
5
4 9
4 6
2 10
4 8
4 8
7
1 2 5
2 3 1
1 6 3
5 1 8
6 2 5
1 3 8
2
6 5
2 2
6
4 1 4
6 3 7
3 4 6
4 5 7
4 5 9
2
3 5
3 1
6
3 6 4
2 3 9
2 6 9
1 2 1
4 3 6
5
1 8
2 4
1 5
4 6
1 10
7
1 2 8
2 6 8
3 1 1
6 2 10
2 6 7
6 7 6
5
2 5
1 3
1 1
4 1
4 8
7
1 2 10
2 3 6
3 4 3
1 6 3
5 6 9
7 3 8
4
3 5
4 5
5 3
5 6
4
1 2 7
2 4 10
3 4 1
3
3 10
3 2
2 3
5
1 2 5
1 3 1
3 4 8
3 2 7
4
3 1
3 4
1 6
3 10
6
1 2 1
3 4 2
2 6 2
4 5 2
5 6 2
1
3 2
6
4 3 5
2 3 1
4 2 3
6 2 8
5 6 1
2
4 1
2 3
5
2 1 9
2 3 9
1 4 7
1 3 2
2
1 3
3 3
6
5 1 4
3 2 8
3 5 8
3 4 1
5 6 9
4
5 7
4 9
5 5
4 3
7
4 2 10
6 7 3
7 2 6
2 1 10
5 6 7
6 7 3
1
5 1
4
1 2 7
2 3 6
3 4 5
1
2 9
7
2 5 10
4 5 10
3 5 6
4 5 6
5 6 6
3 2 2
5
5 4
6 6
6 4
3 7
3 8
4
2 4 10
4 3 1
4 2 2
3
3 4
1 9
2 2
6
1 2 7
2 3 6
3 5 6
3 6 2
5 6 3
3
2 1
3 2
5 1
4
1 2 4
3 2 1
3 4 6
2
2 6
3 10
3
1 2 7
1 2 5
5
1 2
1 10
1 7
1 6
2 1
6
1 2 3
2 4 4
3 4 7
4 5 5
1 6 2
1
4 6
3
1 2 6
2 3 2
3
2 10
1 6
1 7
3
3 1 9
2 3 10
1
1 4
6
1 2 9
3 4 9
5 3 6
4 5 3
5 6 6
5
5 6
3 10
5 4
4 1
4 5
4
2 1 8
4 2 6
1 3 8
2
2 1
1 7
5
2 3 5
1 3 2
3 4 3
4 5 4
4
2 4
3 10
2 3
4 10
7
2 6 5
1 7 9
3 4 8
2 3 9
5 6 8
7 5 8
3
4 10
4 3
1 7
5
1 2 3
2 3 2
2 3 8
4 5 3
4
3 8
1 10
3 4
1 9
7
6 3 9
4 7 5
3 4 7
4 5 6
3 5 3
6 2 10
4
6 2
3 9
2 6
6 8
4
1 2 8
4 3 2
3 4 10
4
3 4
1 9
3 3
1 2
7
1 2 7
2 3 9
3 1 8
3 2 5
1 4 6
3 6 3
3
5 9
2 4
2 5
5
1 2 2
2 1 4
3 4 6
4 5 6
2
2 1
4 5
7
7 4 7
2 3 10
7 4 5
3 5 7
1 2 7
6 7 4
2
5 4
6 6
5
5 3 4
2 1 8
3 2 2
4 5 10
2
3 3
1 10
3
1 2 3
3 1 9
4
2 9
2 5
1 1
1 5
6
1 3 9
2 3 6
2 4 1
4 5 8
5 6 6
2
4 8
2 1
5
2 4 7
3 5 2
3 4 10
1 4 1
5
3 4
2 8
3 10
4 4
4 5
3
2 3 9
3 1 8
5
1 5
1 4
2 9
2 2
1 8
3
1 3 7
1 2 10
1
1 9
6
1 2 1
3 4 8
3 2 5
4 5 2
6 1 2
2
2 2
2 7
7
1 2 9
2 1 5
3 4 9
1 3 8
5 6 7
6 7 6
5
4 7
4 6
5 1
3 5
4 5
6
1 2 4
6 3 8
3 4 8
4 5 6
5 6 4
5
1 10
4 4
4 3
3 2
2 9
6
4 5 10
4 1 6
3 4 7
6 4 5
5 6 3
2
2 10
1 7
3
1 2 7
2 3 9
5
2 10
1 8
2 1
1 1
1 10
3
1 2 3
3 1 2
1
1 5
7
1 2 7
2 6 7
6 4 9
4 5 5
5 6 10
6 7 6
5
4 4
1 5
2 2
2 7
2 4
5
1 3 9
2 5 4
3 4 7
3 2 3
3
1 5
3 4
3 3
5
1 2 4
2 3 9
3 4 8
1 5 7
3
3 9
2 8
2 5
5
1 2 2
4 3 10
1 5 10
2 3 6
5
2 4
3 1
1 6
1 7
3 9
5
4 2 8
2 3 8
5 2 8
3 1 9
3
4 10
4 9
4 1
6
3 2 4
2 3 2
3 4 3
3 1 6
5 6 4
3
1 10
1 7
2 2
7
6 4 1
2 3 7
3 4 9
4 5 2
5 6 2
6 7 4
1
3 2
5
2 1 8
5 2 6
3 4 10
3 1 4
5
1 4
4 3
3 8
4 8
2 10
5
1 2 5
2 5 6
4 2 9
4 5 9
3
3 6
3 3
3 6`

func parseTestcases() ([]testCaseD, error) {
	in := bufio.NewReader(strings.NewReader(testcasesRaw))
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseD, T)
	for i := 0; i < T; i++ {
		var n int
		fmt.Fscan(in, &n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n-1; j++ {
			var x, y, w int
			fmt.Fscan(in, &x, &y, &w)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, w))
		}
		var q int
		fmt.Fscan(in, &q)
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for j := 0; j < q; j++ {
			var eid, nw int
			fmt.Fscan(in, &eid, &nw)
			sb.WriteString(fmt.Sprintf("%d %d\n", eid, nw))
		}
		cases[i] = testCaseD{input: sb.String()}
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expected, err := run("500D.go", tc.input)
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

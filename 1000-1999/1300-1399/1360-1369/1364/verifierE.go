package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `6
6
5
8
6
4
3
8
3
6
4
8
3
2
9
10
10
8
6
3
9
9
5
6
7
4
2
6
4
6
3
6
10
2
10
10
7
5
9
10
7
2
8
7
9
4
6
9
8
2
2
8
5
8
8
6
9
2
3
7
10
10
7
10
8
8
4
3
5
8
7
3
7
10
2
2
7
4
9
5
6
3
8
9
3
7
4
8
5
4
9
8
8
10
9
5
10
7
10
6
5
5
9
10
5
7
4
9`

type testcase struct {
	n int
}

// solve mirrors 1364E.go logic: output identity permutation 0..n-1.
func solve(tc testcase) string {
	var sb strings.Builder
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(i))
	}
	return sb.String()
}

func parseTestcases() ([]testcase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	cases := make([]testcase, 0, len(fields))
	for i, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", i+1, err)
		}
		cases = append(cases, testcase{n: n})
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		input := fmt.Sprintf("%d\n", tc.n)
		expect := solve(tc)
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

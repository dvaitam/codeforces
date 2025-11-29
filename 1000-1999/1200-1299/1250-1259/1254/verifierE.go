package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors the placeholder 1254E solution: always prints 0 after consuming input.
func solve() string {
	return "0"
}

// Tokens of all testcases from testcasesE.txt.
const testcaseData = `
6 1 2 2 3 3 4 1 5 4 6 2 0 1 6 0 0 4 1 2 2 3 1 4 0 0 1 0 3 1 2 1 3 2 3 1 3 1 2 1 3 0 2 3 4 1 2 2 3 2 4 0 1 3 2 3 1 2 2 3 0 3 0 4 1 2 1 3 3 4 4 0 0 0 3 1 2 1 3 2 3 0 5 1 2 2 3 2 4 1 5 4 1 0 0 5 2 1 2 2 1 3 1 2 2 3 0 1 3 2 1 2 2 1 5 1 2 1 3 3 4 1 5 4 3 0 1 5 3 1 2 2 3 1 0 3 6 1 2 2 3 2 4 3 5 4 6 0 3 6 0 4 2 3 1 2 2 3 0 2 3 6 1 2 2 3 2 4 3 5 4 6 3 0 6 0 5 4 6 1 2 2 3 2 4 3 5 4 6 3 2 0 4 5 6 2 1 2 1 2 6 1 2 2 3 3 4 3 5 1 6 0 2 6 0 4 0 2 1 2 1 0 3 1 2 2 3 2 3 1 3 1 2 1 3 2 3 1 6 1 2 1 3 3 4 3 5 1 6 3 5 6 4 2 1 3 1 2 2 3 3 2 1 4 1 2 2 3 1 4 3 4 1 2 3 1 2 2 3 2 3 0 6 1 2 2 3 2 4 3 5 4 6 0 0 1 3 0 0 3 1 2 1 3 0 2 0 5 1 2 1 3 3 4 3 5 3 4 5 2 0 3 1 2 2 3 3 2 0 4 1 2 1 3 1 4 1 0 2 0 5 1 2 2 3 2 4 4 5 4 0 2 5 3 4 1 2 2 3 1 4 2 4 3 1 2 1 2 0 0 2 1 2 2 1 5 1 2 1 3 1 4 1 5 2 1 3 0 0 3 1 2 1 3 1 3 0 6 1 2 2 3 3 4 2 5 1 6 4 1 3 0 6 0 6 1 2 1 3 2 4 1 5 5 6 2 3 5 4 0 0 4 1 2 1 3 1 4 0 1 3 2 6 1 2 2 3 1 4 2 5 4 6 6 0 2 1 0 4 6 1 2 2 3 2 4 2 5 3 6 6 3 4 2 0 0 3 1 2 1 3 2 3 1 4 1 2 2 3 2 4 0 0 4 1 5 1 2 2 3 1 4 3 5 0 2 3 4 1 2 1 2 2 1 2 1 2 2 0 3 1 2 1 3 3 1 2 4 1 2 1 3 2 4 3 2 0 0 3 1 2 2 3 0 0 3 5 1 2 1 3 1 4 3 5 0 1 5 4 2 2 1 2 2 1 2 1 2 0 1 5 1 2 1 3 1 4 3 5 1 5 0 0 3 5 1 2 2 3 3 4 2 5 1 3 4 0 2 4 1 2 2 3 1 4 2 3 1 4 3 1 2 1 3 0 0 0 4 1 2 1 3 2 4 3 2 1 4 5 1 2 1 3 3 4 1 5 4 3 1 0 2 4 1 2 1 3 1 4 0 0 1 4 6 1 2 1 3 3 4 4 5 2 6 0 3 0 2 0 6 3 1 2 2 3 2 1 0 6 1 2 1 3 2 4 1 5 2 6 0 5 4 0 6 0 6 1 2 2 3 3 4 3 5 4 6 4 2 5 6 1 3 2 1 2 2 0 3 1 2 2 3 0 0 1 3 1 2 2 3 3 2 1 5 1 2 1 3 3 4 2 5 4 1 2 5 3 5 1 2 1 3 2 4 3 5 0 0 0 5 1 5 1 2 1 3 2 4 2 5 1 5 4 2 3 3 1 2 1 3 1 3 2 2 1 2 0 2 3 1 2 2 3 2 1 3 4 1 2 1 3 1 4 3 0 2 0 5 1 2 1 3 3 4 1 5 5 4 2 1 0 6 1 2 1 3 2 4 1 5 2 6 1 0 3 2 0 4 3 1 2 2 3 1 2 0 6 1 2 1 3 1 4 2 5 2 6 1 0 6 3 5 2 6 1 2 2 3 1 4 1 5 2 6 2 3 5 0 1 6 2 1 2 0 0 5 1 2 1 3 2 4 3 5 1 0 5 0 2 2 1 2 2 1 2 1 2 0 1 4 1 2 1 3 2 4 1 0 3 2 3 1 2 2 3 2 0 1 4 1 2 1 3 2 4 0 0 4 2 5 1 2 2 3 3 4 4 5 5 1 4 3 2 6 1 2 1 3 1 4 4 5 4 6 0 4 6 1 5 0 5 1 2 2 3 1 4 4 5 3 5 0 1 0 5 1 2 1 3 3 4 1 5 0 0 4 3 0 3 1 2 1 3 2 1 3 3 1 2 1 3 0 2 1 5 1 2 1 3 3 4 2 5 3 0 1 5 0 3 1 2 1 3 3 0 1 2 1 2 2 1 5 1 2 1 3 3 4 1 5 3 5 0 4 0 5 1 2 1 3 3 4 4 5 1 0 3 4 2 3 1 2 1 3 0 1 3 4 1 2 2 3 2 4 0 2 3 1`

// testCase represents one tree input.
type testCase struct {
	n      int
	edges  [][2]int
	values []int
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	idx := 0
	readInt := func() (int, error) {
		if idx >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		val, err := strconv.Atoi(fields[idx])
		idx++
		return val, err
	}

	var cases []testCase
	for idx < len(fields) {
		n, err := readInt()
		if err != nil {
			return nil, fmt.Errorf("failed to read n: %v", err)
		}
		if n <= 0 {
			return nil, fmt.Errorf("invalid n %d", n)
		}
		tc := testCase{n: n}
		tc.edges = make([][2]int, 0, n-1)
		for i := 0; i < n-1; i++ {
			u, err := readInt()
			if err != nil {
				return nil, fmt.Errorf("failed to read edge u: %v", err)
			}
			v, err := readInt()
			if err != nil {
				return nil, fmt.Errorf("failed to read edge v: %v", err)
			}
			tc.edges = append(tc.edges, [2]int{u, v})
		}
		tc.values = make([]int, n)
		for i := 0; i < n; i++ {
			val, err := readInt()
			if err != nil {
				return nil, fmt.Errorf("failed to read value %d: %v", i+1, err)
			}
			tc.values[i] = val
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(strconv.Itoa(e[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e[1]))
		sb.WriteByte('\n')
	}
	for i, v := range tc.values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := formatInput(tc)
		want := solve()
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

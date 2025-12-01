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

const testcasesRaw = `3 3
1 2
2 3
2 1
1 1 0
2 3
4 1
1 2
2 3
2 4
1 1 2
3 1
1 2
2 3
1 1 2
4 3
1 2
2 3
3 4
2 2
1 2 3
2 1
4 1
1 2
2 3
3 4
1 4 3
4 4
1 2
1 3
1 4
2 3
1 1 0
2 3
2 3
3 2
1 2
2 3
1 3 3
2 1
4 4
1 2
1 3
2 4
1 1 1
2 2
1 3 1
2 4
2 1
1 2
1 2 2
2 3
1 2
2 1
2 1
2 1
5 3
1 2
2 3
2 4
2 5
2 5
1 3 0
2 2
4 3
1 2
1 3
2 4
1 4 1
1 1 0
1 2 1
2 5
1 2
1 2 0
1 2 3
1 2 1
1 2 3
2 1
3 4
1 2
1 3
2 1
2 1
1 1 2
2 1
3 2
1 2
2 3
1 2 0
2 3
2 5
1 2
1 2 3
1 2 1
1 2 2
2 2
2 1
4 3
1 2
2 3
1 4
2 2
1 4 1
2 1
3 4
1 2
2 3
2 1
1 1 1
1 3 2
1 2 3
3 5
1 2
1 3
1 3 2
2 3
2 1
1 1 0
2 3
4 5
1 2
1 3
2 4
1 1 2
2 1
1 4 3
1 3 3
1 4 3
2 1
1 2
2 2
5 4
1 2
1 3
3 4
3 5
2 4
1 5 1
2 2
1 2 2
5 3
1 2
1 3
3 4
2 5
1 3 0
1 4 3
1 5 1
5 3
1 2
1 3
1 4
1 5
1 5 2
2 5
2 5
5 5
1 2
2 3
1 4
3 5
2 3
1 4 0
1 5 1
2 3
1 5 2
2 1
1 2
1 2 2
4 2
1 2
1 3
3 4
1 1 2
1 2 1
5 5
1 2
1 3
1 4
4 5
1 5 2
2 5
1 4 1
2 3
1 3 3
5 5
1 2
2 3
2 4
4 5
2 5
1 3 1
2 5
2 4
2 5
2 2
1 2
2 2
2 1
4 4
1 2
2 3
3 4
1 2 1
1 3 1
1 2 3
2 1
3 2
1 2
1 3
1 2 0
1 2 0
4 5
1 2
1 3
3 4
2 2
1 3 0
1 4 1
1 3 3
1 4 1
3 3
1 2
2 3
2 2
1 2 0
2 1
4 1
1 2
1 3
1 4
1 1 0
3 2
1 2
2 3
1 2 3
2 1
3 1
1 2
2 3
2 3
5 5
1 2
1 3
2 4
3 5
2 3
1 2 2
1 1 2
2 3
1 2 0
5 5
1 2
1 3
1 4
1 5
2 5
1 3 0
1 1 1
2 5
1 1 0
5 2
1 2
2 3
1 4
3 5
2 4
2 1
3 4
1 2
1 3
2 3
2 1
2 2
1 1 0
5 4
1 2
2 3
2 4
4 5
2 5
1 3 2
1 1 3
2 5
4 2
1 2
1 3
1 4
1 2 2
1 4 0
4 1
1 2
1 3
1 4
1 2 0
2 3
1 2
1 1 1
1 1 1
2 1
4 3
1 2
1 3
1 4
1 1 0
1 1 2
2 1
3 2
1 2
2 3
2 3
2 2
5 4
1 2
1 3
1 4
2 5
2 2
2 2
2 5
1 5 2
5 3
1 2
1 3
1 4
3 5
1 5 1
1 1 0
1 3 2
5 1
1 2
2 3
1 4
2 5
2 5
5 3
1 2
1 3
2 4
2 5
2 1
2 2
1 3 2
2 1
1 2
1 1 1
3 5
1 2
1 3
1 2 3
1 3 1
1 1 1
2 1
2 3
4 5
1 2
1 3
3 4
2 3
2 3
1 2 2
1 4 3
1 1 1
3 2
1 2
1 3
2 1
1 3 1
3 3
1 2
2 3
2 2
1 3 3
2 3
2 2
1 2
2 1
1 2 2
2 2
1 2
1 2 0
2 1
3 2
1 2
1 3
1 1 0
2 3
5 1
1 2
2 3
3 4
2 5
1 5 0
2 1
1 2
2 1
2 5
1 2
2 2
1 1 0
1 1 0
2 2
2 2
3 5
1 2
2 3
1 3 3
1 2 0
1 2 3
2 1
2 1
4 3
1 2
1 3
3 4
2 3
2 3
1 3 0
3 3
1 2
1 3
2 1
2 1
1 1 0
3 2
1 2
2 3
2 2
1 3 3
5 3
1 2
2 3
1 4
1 5
1 1 3
1 4 1
2 1
3 3
1 2
2 3
1 1 2
2 1
2 3
3 4
1 2
2 3
1 2 1
1 2 3
1 1 0
2 1
3 2
1 2
2 3
2 2
2 2
4 4
1 2
2 3
2 4
1 4 3
2 2
1 4 2
1 3 0
3 5
1 2
2 3
1 3 2
2 2
2 3
1 2 0
2 2
5 3
1 2
1 3
2 4
3 5
1 3 3
1 4 3
2 3
2 4
1 2
1 1 0
2 2
1 1 0
1 2 1
4 3
1 2
1 3
2 4
1 1 0
2 1
1 2 3
2 5
1 2
1 1 0
1 1 1
1 2 2
2 2
2 1
4 1
1 2
2 3
1 4
1 3 0
5 1
1 2
1 3
1 4
4 5
2 3
2 1
1 2
2 2
5 1
1 2
1 3
1 4
2 5
2 5
5 4
1 2
2 3
1 4
3 5
2 5
1 3 3
2 3
1 2 1
3 4
1 2
1 3
1 3 0
1 2 0
1 3 3
1 3 0
4 3
1 2
1 3
1 4
2 4
2 3
1 1 3
4 5
1 2
2 3
3 4
2 2
2 2
1 2 2
1 2 1
2 3
3 5
1 2
2 3
2 3
2 2
2 3
1 2 3
2 2
3 1
1 2
2 3
2 3
5 1
1 2
2 3
3 4
1 5
2 5
2 2
1 2
1 2 3
1 1 2
3 2
1 2
1 3
1 2 3
2 2
4 3
1 2
1 3
3 4
1 1 1
1 4 2
2 1
3 3
1 2
2 3
2 2
1 2 3
2 2
2 4
1 2
1 2 3
2 1
2 2
2 1
3 2
1 2
1 3
1 2 0
2 3
3 1
1 2
2 3
2 2
3 1
1 2
1 3
2 2
5 1
1 2
2 3
3 4
4 5
1 1 2
4 3
1 2
2 3
2 4
1 2 3
2 1
1 1 0
5 5
1 2
2 3
3 4
1 5
2 4
2 5
2 5
1 4 1
2 3
4 1
1 2
1 3
2 4
2 3
5 5
1 2
1 3
2 4
4 5
2 5
2 2
1 2 2
1 3 3
2 3
2 1`

type testCase struct {
	n     int
	q     int
	edges [][2]int
	ops   [][]int
}

// solve mirrors the stub 1254D.go: outputs 0 for each query of type 2.
func solve(tc testCase) string {
	var sb strings.Builder
	for _, op := range tc.ops {
		if op[0] == 2 {
			sb.WriteString("0\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	var tests []testCase
	for {
		if !scan.Scan() {
			break
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, err
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("invalid test file")
		}
		q, _ := strconv.Atoi(scan.Text())
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("invalid test file")
			}
			u, _ := strconv.Atoi(scan.Text())
			if !scan.Scan() {
				return nil, fmt.Errorf("invalid test file")
			}
			v, _ := strconv.Atoi(scan.Text())
			edges[i] = [2]int{u, v}
		}
		ops := make([][]int, q)
		for i := 0; i < q; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("invalid test file")
			}
			t, _ := strconv.Atoi(scan.Text())
			if t == 1 {
				scan.Scan()
				v, _ := strconv.Atoi(scan.Text())
				scan.Scan()
				d, _ := strconv.Atoi(scan.Text())
				ops[i] = []int{1, v, d}
			} else {
				scan.Scan()
				v, _ := strconv.Atoi(scan.Text())
				ops[i] = []int{2, v}
			}
		}
		tests = append(tests, testCase{n: n, q: q, edges: edges, ops: ops})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for _, op := range tc.ops {
		if op[0] == 1 {
			fmt.Fprintf(&sb, "1 %d %d\n", op[1], op[2])
		} else {
			fmt.Fprintf(&sb, "2 %d\n", op[1])
		}
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	expected := solve(tc)
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

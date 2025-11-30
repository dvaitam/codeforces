package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesH.txt. Each line: q followed by q operations (push x or pop).
const testcasesRaw = `6 push b push a pop push c push c push a
7 push a push b push c push a push c pop push b
1 push a
3 push a push c push c
3 push c push b push c
2 push c push c
7 push b push b push a pop push a push b push b
8 push a push b push b push b push c push c push b push b
8 push b push a pop push c push c pop push c pop
5 push c push b push c push a push b
3 push b push b push b
7 push a pop pop push b pop push b pop
3 push a push a push c
3 push a push c push c
6 push c pop pop push c push c pop
7 push b push c push a push a push a push a push c
3 push b push a pop
7 push b pop push b push b pop push b push a
2 push c push c
3 push a pop push a
1 push b
2 push b push a
6 push c push c push a push a pop pop
4 push b push a pop pop
8 push c push b push c pop pop push a push b push a
8 push c pop push c push c push c push b pop pop
8 push b pop push c push b pop push a push a push c
8 push a push c pop push b push c push a pop pop
2 push a push a
4 push a push a pop push c
7 push a pop push b push c push c push a push a
8 push c push a push b push a push b push c push a pop
1 push b
1 push c
8 push a push c pop push a push b push b push a push b
4 push a push b push a push b
2 push b pop
3 push c push a push a
8 push a push b push c pop push b pop push b push b
2 push a push b
8 push b push c push a push a pop push b push a push a
7 push c pop push a push c push c push b push c
3 push a push a push b
2 push a push a
8 push c push b push a push a pop pop push a push b
5 push a push c push b push a pop
1 push c
4 push a pop push c push b
7 push b push a push c pop push b pop pop
1 push c
5 push a push b pop push b push c
5 push a push b push b pop pop
6 push b push a push b push b push a push a
2 push a push c
1 push b
5 push a push c pop push c pop
7 push c pop push c push a push c pop push c
3 push c push a pop
4 push a push b pop push b
1 push c
4 push a push a pop push c
2 push a pop
8 push a pop push a push c pop push b push c push a
4 push a push c pop push c
3 push a push c push a
8 push c push b push a pop push b pop pop push b
2 push b pop
7 push a push a push c push b push c push a pop
4 push b push a push b push b
5 push b push b push b push b pop
2 push a pop
5 push a push c push b pop pop
5 push c push a push b push b push b
7 push c push c push c push c push c pop push b
8 push c push b push b push c pop push b push b push c
7 push c push a push b push b pop pop push a
4 push a push a push b push c
4 push c pop push c push b
5 push a push c pop push c push a
2 push a push c
8 push b pop pop push a push b pop push a push c
8 push a pop push a push c push c pop pop pop
2 push a push a
4 push a push b pop push c
7 push a push b push c push b push c push c push b
5 push c pop push c push a push b
5 push a push c push a push c push b
4 push c push b push b push b
2 push a push a
4 push b pop push c push a
4 push c push c push b push a
3 push c push b push c
6 push a pop push a push c pop push c
7 push b push a push b push b pop push a push b
7 push c push b push b push b pop push c push c
7 push a push a push c push c push b push c push c
6 push b push c pop push a push a push c
8 push b pop pop pop push b push c push c pop
4 push b push c push a push a
4 push c push c push a pop`

type testCase struct {
	q   int
	ops []string
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		q, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse q: %v", idx+1, err)
		}
		ops := make([]string, 0, q)
		i := 1
		for len(ops) < q && i < len(fields) {
			if fields[i] == "push" {
				if i+1 >= len(fields) {
					return nil, fmt.Errorf("line %d incomplete push", idx+1)
				}
				ops = append(ops, fields[i]+" "+fields[i+1])
				i += 2
			} else {
				ops = append(ops, fields[i])
				i++
			}
		}
		if len(ops) != q {
			return nil, fmt.Errorf("line %d expected %d ops got %d", idx+1, q, len(ops))
		}
		cases = append(cases, testCase{q: q, ops: ops})
	}
	return cases, nil
}

// Embedded solver logic from 1738H.go.
func isPalindrome(b []byte) bool {
	i, j := 0, len(b)-1
	for i < j {
		if b[i] != b[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func solve(tc testCase) []int {
	queue := make([]byte, 0)
	res := make([]int, 0, tc.q)
	for _, op := range tc.ops {
		parts := strings.Fields(op)
		if parts[0] == "push" {
			queue = append(queue, parts[1][0])
		} else if parts[0] == "pop" {
			if len(queue) > 0 {
				queue = queue[1:]
			}
		}
		seen := make(map[string]struct{})
		for i := 0; i < len(queue); i++ {
			for j := i; j < len(queue); j++ {
				if isPalindrome(queue[i : j+1]) {
					seen[string(queue[i:j+1])] = struct{}{}
				}
			}
		}
		res = append(res, len(seen))
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
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
		want := solve(tc)
		var input strings.Builder
		input.WriteString(strconv.Itoa(tc.q))
		input.WriteByte('\n')
		for _, op := range tc.ops {
			input.WriteString(op)
			input.WriteByte('\n')
		}

		gotStr, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotParts := strings.Fields(strings.TrimSpace(gotStr))
		if len(gotParts) != tc.q {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", idx+1, tc.q, len(gotParts))
			os.Exit(1)
		}
		for i := 0; i < tc.q; i++ {
			val, err := strconv.Atoi(gotParts[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d output parse error at %d: %v\n", idx+1, i, err)
				os.Exit(1)
			}
			if val != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at step %d expected %d got %d\n", idx+1, i+1, want[i], val)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

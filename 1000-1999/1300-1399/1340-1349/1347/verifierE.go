package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt to remove external dependency.
const testcasesRaw = `100
19
18
15
3
5
14
10
5
20
13
9
19
11
9
9
4
18
11
12
9
13
17
11
20
7
6
2
19
18
12
13
20
2
6
14
6
7
18
4
6
8
17
20
8
9
6
9
14
13
20
6
17
5
2
18
13
17
16
11
2
9
19
7
17
17
19
12
4
10
6
14
8
12
11
14
3
8
3
12
9
12
16
9
10
13
7
11
2
13
20
19
3
6
13
2
17
3
2
9
3`

type testCase struct {
	n int
}

// referenceSolution embeds logic from 1347E.go to produce the permutation or -1.
func referenceSolution(n int) []int {
	if n == 2 || n == 3 {
		return nil
	}
	if n == 4 {
		return []int{2, 4, 1, 3}
	}
	res := make([]int, 0, n)
	for i := 1; i <= n; i += 2 {
		res = append(res, i)
	}
	if n%2 == 0 {
		evenSeq := []int{n - 4, n, n - 2}
		for x := n - 6; x >= 2; x -= 2 {
			evenSeq = append(evenSeq, x)
		}
		res = append(res, evenSeq...)
	} else {
		evenSeq := []int{n - 3, n - 1}
		for x := n - 5; x >= 2; x -= 2 {
			evenSeq = append(evenSeq, x)
		}
		res = append(res, evenSeq...)
	}
	return res
}

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil
	}
	count, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	tests := make([]testCase, 0, count)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			panic("invalid testcase")
		}
		tests = append(tests, testCase{n: n})
	}
	return tests
}

func runExe(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := parseTestcases()

	for i, tc := range tests {
		p := referenceSolution(tc.n)
		expected := ""
		if p == nil {
			expected = "-1"
		} else {
			for idx, v := range p {
				if idx > 0 {
					expected += " "
				}
				expected += strconv.Itoa(v)
			}
		}

		input := fmt.Sprintf("1\n%d\n", tc.n)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("Test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

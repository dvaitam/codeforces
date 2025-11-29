package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n0 int
	n1 int
	n2 int
}

// Embedded testcases from testcasesF.txt.
const testcasesRaw = `100
1 2 0
5 3 3
1 0 0
0 3 4
2 0 1
4 4 2
2 1 0
2 1 0
5 2 2
1 1 2
2 5 5
2 0 4
2 5 3
4 1 1
1 3 2
0 4 2
0 2 4
5 2 4
1 3 3
4 2 3
3 1 1
2 2 0
0 0 3
5 2 4
4 5 3
5 2 1
5 1 0
3 1 5
5 3 2
1 2 3
5 4 2
5 4 1
2 0 0
5 1 2
4 4 1
0 2 1
2 3 0
0 2 5
0 2 5
5 2 0
2 2 2
1 5 3
4 5 0
2 4 1
3 2 1
2 3 4
1 2 4
0 2 0
3 1 2
2 2 4
0 3 1
3 1 0
0 0 0
5 1 4
5 1 4
0 4 3
4 1 2
0 0 4
2 3 5
1 3 1
1 3 3
3 0 1
3 3 1
5 3 1
3 1 0
0 2 2
1 4 1
1 3 2
1 2 0
2 4 0
4 3 5
5 5 5
0 3 3
0 3 1
4 1 2
2 5 3
5 2 3
4 1 5
5 2 2
3 3 0
2 5 5
1 0 3
4 1 2
5 0 1
5 5 3
4 3 5
3 3 1
0 1 1
0 4 2
0 3 3
1 4 0
1 1 5
4 2 4
3 4 3
0 0 0
5 4 0
3 4 2
4 1 0
2 0 4
0 2 2`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %v", err)
	}
	if len(fields) != 1+t*3 {
		return nil, fmt.Errorf("malformed testcases: expected %d numbers, got %d", 1+t*3, len(fields))
	}
	cases := make([]testCase, 0, t)
	idx := 1
	for i := 0; i < t; i++ {
		n0, _ := strconv.Atoi(fields[idx])
		n1, _ := strconv.Atoi(fields[idx+1])
		n2, _ := strconv.Atoi(fields[idx+2])
		idx += 3
		cases = append(cases, testCase{n0: n0, n1: n1, n2: n2})
	}
	return cases, nil
}

// referenceSolution mirrors the construction from 1352F.go in this repo.
func referenceSolution(n0, n1p, n2p int) string {
	// remap variables to match original code naming
	n1 := n2p
	n2 := n1p
	n3 := n0

	var sBuilder []byte
	var tBuilder []byte
	tBuilder = append(tBuilder, '1')
	cFlag := false

	if n2 == 0 {
		if n1 > 0 {
			sBuilder = append(sBuilder, '1')
			for i := 0; i < n1; i++ {
				sBuilder = append(sBuilder, '1')
			}
		} else {
			sBuilder = append(sBuilder, '0')
			for i := 0; i < n3; i++ {
				sBuilder = append(sBuilder, '0')
			}
		}
	} else if n2%2 == 1 {
		cnt := n2
		for cnt > 0 {
			if cFlag {
				tBuilder = append(tBuilder, '1')
				cFlag = false
			} else {
				tBuilder = append(tBuilder, '0')
				cFlag = true
			}
			cnt--
		}
		for i := 0; i < n1; i++ {
			sBuilder = append(sBuilder, '1')
		}
		sBuilder = append(sBuilder, tBuilder...)
		for i := 0; i < n3; i++ {
			sBuilder = append(sBuilder, '0')
		}
	} else {
		cnt := n2 - 1
		for cnt > 0 {
			if cFlag {
				tBuilder = append(tBuilder, '1')
				cFlag = false
			} else {
				tBuilder = append(tBuilder, '0')
				cFlag = true
			}
			cnt--
		}
		for i := 0; i < n3; i++ {
			tBuilder = append(tBuilder, '0')
		}
		tBuilder = append(tBuilder, '1')
		for i := 0; i < n1; i++ {
			sBuilder = append(sBuilder, '1')
		}
		sBuilder = append(sBuilder, tBuilder...)
	}
	return string(sBuilder)
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("1\n%d %d %d\n", tc.n0, tc.n1, tc.n2)
	out, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	// Use first token as candidate string.
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output produced")
	}
	s := fields[0]
	exp := referenceSolution(tc.n0, tc.n1, tc.n2)
	if s != exp {
		return fmt.Errorf("output mismatch: expected %q got %q", exp, s)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

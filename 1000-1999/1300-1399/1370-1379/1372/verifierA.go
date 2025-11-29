package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
100
7
7
1
5
9
8
7
5
8
6
10
4
9
3
5
3
2
10
5
9
10
3
5
2
2
6
8
9
2
6
7
6
10
4
9
8
8
9
5
1
9
1
2
7
1
10
8
6
4
6
2
4
10
4
4
3
9
8
2
2
6
9
8
2
5
9
5
2
9
6
9
4
10
9
10
5
8
2
10
7
6
10
4
5
3
4
3
1
10
5
8
2
2
3
3
1
2
9
7
9
`

type testCase struct {
	n int
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

// solve replicates 1372A: output n ones.
func solve(tc testCase) []int {
	arr := make([]int, tc.n)
	for i := range arr {
		arr[i] = 1
	}
	return arr
}

func parseTestcases() ([]testCase, string, error) {
	fields := strings.Fields(testcasesData)
	if len(fields) == 0 {
		return nil, "", fmt.Errorf("no testcases")
	}
	pos := 0
	t := 0
	if v, err := strconv.Atoi(fields[pos]); err == nil {
		t = v
		pos++
	}
	if t == 0 {
		return nil, "", fmt.Errorf("missing testcase count")
	}
	cases := make([]testCase, 0, t)
	var inputBuilder strings.Builder
	inputBuilder.WriteString(strconv.Itoa(t))
	inputBuilder.WriteByte('\n')
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, "", fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, "", err
		}
		pos++
		cases = append(cases, testCase{n: n})
		inputBuilder.WriteString(strconv.Itoa(n))
		inputBuilder.WriteByte('\n')
	}
	return cases, inputBuilder.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, input, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	output, err := run(bin, input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outFields := strings.Fields(output)
	pos := 0
	for idx, tc := range testcases {
		want := solve(tc)
		if pos+tc.n > len(outFields) {
			fmt.Printf("case %d missing output values\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(outFields[pos+i])
			if err != nil {
				fmt.Printf("case %d invalid number %q\n", idx+1, outFields[pos+i])
				os.Exit(1)
			}
			if val < 1 || val > 1000 {
				fmt.Printf("case %d value out of range: %d\n", idx+1, val)
				os.Exit(1)
			}
			if val != want[i] {
				fmt.Printf("case %d failed: expected %d got %d\n", idx+1, want[i], val)
				os.Exit(1)
			}
		}
		pos += tc.n
	}
	if pos != len(outFields) {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

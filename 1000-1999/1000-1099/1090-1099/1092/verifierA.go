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

const testcasesRaw = `50 25
54 2
34 17
63 13
39 16
46 19
28 17
18 10
18 4
80 26
33 18
91 26
78 5
40 4
94 3
88 11
61 18
13 6
56 11
79 21
27 18
62 15
67 9
8 1
12 12
52 23
86 21
1 1
43 8
94 11
91 3
25 19
29 8
19 18
58 3
11 6
66 16
14 5
71 10
91 4
71 11
70 7
78 18
76 10
57 3
77 26
50 11
74 8
38 6
25 6
5 5
85 9
61 3
12 11
97 5
20 2
11 9
88 13
91 17
36 17
31 7
87 19
54 19
36 15
64 22
83 23
46 3
42 20
15 8
76 21
43 7
32 1
94 9
15 12
29 12
22 11
55 2
13 13
19 8
6 5
82 18
78 22
10 1
16 7
78 19
16 13
12 6
15 1
78 1
25 6
92 4
62 7
94 26
8 1
70 14
80 4
34 3
29 3
83 10
45 14`

// solve mirrors 1092A.go logic.
func solve(n, k int) string {
	var sb strings.Builder
	p := n / k
	for i := 0; i < k; i++ {
		ch := byte('a' + i)
		for j := 0; j < p; j++ {
			sb.WriteByte(ch)
		}
	}
	rem := n - k*p
	if rem > 0 {
		last := byte('a' + k - 1)
		for i := 0; i < rem; i++ {
			sb.WriteByte(last)
		}
	}
	return sb.String()
}

type testCase struct {
	n int
	k int
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
		k, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{n: n, k: k})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
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
	expected := solve(tc.n, tc.k)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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

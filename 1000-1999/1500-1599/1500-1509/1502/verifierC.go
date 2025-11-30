package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesC.txt.
const testcaseData = `
2
33
2
35
40
24
11
22
37
5
5
35
11
16
12
16
20
16
16
33
29
9
28
35
9
2
40
37
11
32
2
20
4
12
29
39
15
29
33
10
21
8
30
35
3
34
5
33
21
0
5
6
27
38
22
36
28
21
24
32
23
40
7
8
20
1
11
8
1
21
38
12
2
26
3
19
24
3
38
10
22
4
26
3
28
22
38
39
16
19
36
29
26
11
1
29
16
12
24
4
`

func parseTestcases() ([]int, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]int, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("case %d bad integer: %v", i+1, err)
		}
		res = append(res, v)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

// fib mirrors 1502C.go.
func fib(n int) uint64 {
	if n <= 1 {
		return uint64(n)
	}
	a, b := uint64(0), uint64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func runCandidate(bin string, n int) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%d\n", n))
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, n := range tests {
		expect := fib(n)
		gotStr, err := runCandidate(bin, n)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseUint(gotStr, 10, 64)
		if err != nil || got != expect {
			fmt.Printf("test %d failed: expected %d got %s\n", idx+1, expect, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

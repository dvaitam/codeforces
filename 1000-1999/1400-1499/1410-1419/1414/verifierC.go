package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `100
5
0
38
12
21
10
15
14
40
28
24
36
26
2
25
36
26
2
10
28
4
16
10
28
33
31
35
38
0
2
31
20
19
29
3
26
12
35
40
5
8
0
25
26
20
0
13
0
0
33
39
6
12
7
38
12
19
17
11
6
30
25
40
5
1
17
28
7
16
8
33
22
7
9
17
1
2
2
13
16
35
20
23
36
2
38
31
29
40
27
23
34
11
13
24
37
18
0
8
9`

func fib(n int) int64 {
	if n <= 1 {
		return int64(n)
	}
	a, b := int64(0), int64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func parseTestcases() ([]int, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %v", err)
	}
	if len(fields) != 1+t {
		return nil, fmt.Errorf("malformed testcases: expected %d numbers, got %d", 1+t, len(fields))
	}
	vals := make([]int, t)
	for i := 0; i < t; i++ {
		v, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return nil, fmt.Errorf("parse value %d: %v", i+1, err)
		}
		vals[i] = v
	}
	return vals, nil
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	vals, err := parseTestcases()
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

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(vals))
	for _, v := range vals {
		fmt.Fprintf(&input, "%d\n", v)
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var expected strings.Builder
	for i, v := range vals {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(strconv.FormatInt(fib(v), 10))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(vals))
}

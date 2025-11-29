package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `100
12
13
1
8
16
15
12
9
15
11
18
6
16
4
9
4
3
19
8
17
19
4
9
3
2
10
15
17
3
11
13
10
19
20
6
17
15
14
16
8
1
17
0
2
12
20
0
19
15
10
7
10
2
6
18
7
7
4
17
14
2
2
10
16
15
3
9
17
9
3
17
10
17
6
19
17
18
9
14
2
19
12
10
18
7
9
5
6
5
1
19
8
15
2
2
4
4
1
2
17`

func factorial(n int64) int64 {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func parseTestcases() ([]int64, error) {
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
	vals := make([]int64, t)
	for i := 0; i < t; i++ {
		v, err := strconv.ParseInt(fields[1+i], 10, 64)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		expected.WriteString(strconv.FormatInt(factorial(v), 10))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(vals))
}

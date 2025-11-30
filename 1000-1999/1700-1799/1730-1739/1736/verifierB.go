package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
3 10 2 5
2 8 8
8 7 4 2 8 1 7 7 10
1 8
5 4 10 2 6 1
1 1
1 7
4 7 1 9 4
8 8 9 4 6 4 4 8 5
1 7
2 3 5
2 6 9
7 9 4 5 5 10 8 9
7 10 1 8 4 7 7 3
6 9 6 2 8 9 2
3 9 7 6
8 1 8 1 5 10 10 10 7
3 3 9 4
1 4
4 7 9 6 10
6 8 5 9 10 1 7
3 9 9 4
7 1 8 6 10 9 4 9
7 8 6 7 6 1 9 9
6 8 10 1 4 3 9
3 2 9 5
1 2
2 1 8
1 5
4 5 2 10 3
6 5 2 3 3 5 9
3 5 5 8
6 8 8 2 1 5 7
6 7 4 5 2 5 9
4 10 7 1 4
1 7
3 1 3 8
7 9 4 9 8 4 9 1
7 10 6 7 1 5 3 4
1 5
2 2 5
5 3 7 10 5 3
1 9
1 10
4 10 8 3 10
1 7
4 6 2 4 10
7 10 4 8 2 7 5 9
8 1 6 10 7 5 1 3 4
6 10 3 6 7 4 5
2 7 9
6 9 8 9 4 2 1
2 3 3
3 9 4 5
6 10 9 5 6 6 6
2 5 4
8 3 10 9 2 6 1 7 2
7 3 3 6 2 10 10 7
2 10 9
4 10 2 5 6
5 10 9 2 8 5
2 1 5
1 10
1 2
7 2 1 4 4 10 7 3
2 8 3
4 3 2 7 7
5 9 5 8 6 2
4 6 1 1 1
5 10 6 8 7 6
7 2 2 6 10 8 2 5
4 10 9 8 6
5 3 9 4 5 4
4 6 2 5 2
8 2 10 6 4 7 5 1 6
3 6 10 5
4 6 2 9 10
2 4 4
1 4
7 2 5 9 2 2 1 1
5 6 8 8 3 2
6 2 9 3 3 3 3
6 5 2 9 10 5 3
4 3 9 1 6
4 3 5 7 9
3 1 4 5
2 8 7
5 9 8 9 8 1
7 6 3 5 8 1 7 10
1 1
6 10 3 10 3 3 5
5 7 10 7 3 10
2 4 8
1 3
6 9 8 4 4 6 8
8 4 7 6 9 10 5 4 1
2 9 6
3 9 4 5
5 5 9 6 3 8
2 2 10
`

type testCase struct {
	n int
	a []int64
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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func solve(arr []int64) string {
	n := len(arr)
	ok := true
	for i := 1; i < n-1; i++ {
		l1 := lcm(arr[i-1], arr[i])
		l2 := lcm(arr[i], arr[i+1])
		if gcd(l1, l2) != arr[i] {
			ok = false
			break
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, a: arr})
	}
	return cases, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		want := solve(tc.a)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(v, 10))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

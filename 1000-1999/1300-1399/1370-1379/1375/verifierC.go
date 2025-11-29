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

const testcasesC = `2 -4 -4
7 -3 5 -1 -1 4 -2 4
2 4 5
4 1 5 1 3
7 3 2 3 -1 -5 -5 0
9 0 1 1 3 -3 3 -3 -2 -2
2 -3 0
4 -3 3 3 0
10 5 3 -3 2 1 3 0 4 0 0
9 -3 1 2 5 3 -2 2 -1 2
10 3 0 5 2 2 0 4 3 2 2
5 0 -3 4 -1 2
6 -1 3 3 3 3 5
8 -1 -2 2 3 0 5 4 -4
7 -5 -2 -4 -5 4 5 -5
6 4 -2 5 -4 3 -3
6 -2 -2 -5 1 -5 -5
7 0 -3 -2 5 -5 -4 -4
3 -5 -5 -5
7 -1 -3 -3 -3 3 -5 1
2 -2 -3
2 -5 0
3 -1 0 2
2 -1 2
10 4 -5 -1 1 4 -3 2 -2 -4 5
7 -4 -5 2 -3 3 4 1
9 3 0 -3 0 -1 -1 4 1 5
2 3 -3
2 -1 -5
4 -3 -3 -4 2
5 3 -5 -2 -2 2
3 -1 -4 4
5 4 4 0 -1 5
8 -1 3 -5 -3 -5 1 1 -3
3 3 -4 -2
3 -4 -5 -3
5 -4 -2 -5 3 5
9 2 -1 3 5 1 -2 5 -2 1
8 3 -5 4 4 -5 1 3 4
4 -4 5 2 0
2 3 -4
7 -1 0 -1 -5 5 1 -4
3 -1 -2 5
2 2 -5
8 5 2 2 -2 4 4 -4 -5
6 -5 0 -1 -4 -2 2
5 -4 4 0 1 2
4 0 1 -4 -1
3 -4 -4 4
7 5 1 -2 -4 -5 4 5
9 -5 2 -1 0 2 -3 0 -1 2
10 2 1 2 5 -1 1 -2 -3 2 4
6 3 1 5 -4 4 4
3 -4 0 -3
10 -3 1 -4 -4 5 5 -5 -3 -1 1
5 5 5 0 2 -3
10 -1 -4 -3 3 1 -4 0 3 -2 3
6 -3 -3 2 -2 1 0
4 2 2 -5 4
8 -3 1 3 -5 2 -1 1 -1
8 5 2 0 3 0 5 -4 -2
10 4 -2 1 5 1 5 -5 0 2 3
9 5 -3 -4 -5 1 -2 4 4 1
5 -4 1 3 -2 -1
5 2 4 -3 -5 4
8 2 -1 3 4 -3 2 -2 -4
7 -5 2 3 5 5 -4 4
9 5 0 2 -1 3 2 -5 -4 4
7 -3 1 -1 5 5 -3 -5
4 2 1 2 5
6 -3 -5 -1 3 2 -5
7 -5 3 1 4 2 -2 5
6 2 5 -3 2 3 -1
3 -1 0 -1
7 5 -1 5 5 1 3 -4
10 5 -2 1 4 3 -3 3 5 -4 -1
2 -2 2
10 -2 3 -1 -5 -4 -4 5 1 0 -2
7 0 -4 0 2 0 -3 2
9 -1 2 -3 2 5 -2 -1 0 -3
3 -2 2 -2
7 -3 0 -3 -3 -2 -1 3
8 1 0 -1 4 3 4 0 1
6 3 4 5 5 -4 0
6 1 2 -3 -1 0 2
9 -4 -3 0 1 -3 -5 -4 0 -3
7 -4 5 1 -5 3 0 -2
8 3 -1 2 5 -3 0 0 -2
9 -4 -3 -2 0 -1 -3 1 0 -1
3 0 -2 -2
5 4 -5 0 0 5
2 -3 -3
3 1 2 -1
4 0 3 4 -4
7 5 4 1 -2 -5 1 2
9 4 0 3 4 4 -4 4 3 3
9 1 2 -3 1 1 3 2 -5 -4
9 4 -3 -4 5 3 -3 -4 1 -1
9 -5 -1 -4 5 0 -2 -3 -5 -3
8 5 -4 0 5 2 -5 2 -2`

type testCase struct {
	n   int
	arr []int
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

func referenceSolve(arr []int) string {
	if arr[0] < arr[len(arr)-1] {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesC))
	scan.Split(bufio.ScanWords)
	cases := make([]testCase, 0)
	for {
		if !scan.Scan() {
			break
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n: %w", err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("missing value %d for case %d", i+1, len(cases)+1)
			}
			v, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("parse value %d case %d: %w", i+1, len(cases)+1, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		want := referenceSolve(tc.arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != want {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

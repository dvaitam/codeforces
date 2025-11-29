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

const testcases = `100
7 7 2 9 17 16 13 10 16
6 10 7 17 5 10 5 4
10 5 18 20 5 10 4 3 11 16 18 4
6 7 11 20 7 18 16 15
9 5 2 18 1 3 13 1 20 16 11
4 6 3 7 19 8
4 3 18 15 3 3
6 9 16 4 10 18 10 4
9 6 18 7 20 18 19 10 15 3 20
7 6 19 8 10 6 7 6 2
10 5 16 3 3 5 5 2 3 18 13 17
5 9 8 7 19 14 19
5 8 16 12 3 11 20
2 8 19 11
4 4 1 9 4 8
6 3 11 14 2 4 5 8
1 10 18
10 2 1 4 7 20 19 4 13 3 12 4
1 10 1
4 3 4 16 7 2
1 9 14
10 2 9 3 8 3 10 12 14 6 2 17
8 1 20 4 13 7 9 12 16 19
3 4 2 6 6
6 9 9 4 20 15 6 1
8 7 19 17 10 12 13 9 5 18
1 8 3
6 1 18 9 5 8 16 12
10 5 12 19 20 5 10 13 14 3 1 20
4 6 6 8 8 15
7 10 14 2 13 19 14 2 6
8 2 9 6 15 17 16 18 20 1
1 8 11
5 8 2 14 7 18 3
3 1 13 14 11
1 4 1
1 9 20
2 4 4 20
4 5 9 6 4 16
7 2 1 9 15 4 9 5 17
6 2 5 9 1 2 2 7
5 9 11 12 19 2 20
8 8 14 12 18 6 7 13 19 10
1 3 5
5 6 11 12 3 11 20
1 1 9
3 3 19 10 12
7 9 5 10 4 16 8 2 10
3 9 3 10 13
6 5 14 4 4 18 16 16
6 6 4 16 4 16 14 2
5 6 5 6 19 13 3
2 2 7 8
1 7 1
2 7 18 17
5 8 16 19 7 14 3
6 4 9 19 6 14 7 12
2 2 1 17
8 4 4 16 13 9 7 2 7 20
3 2 7 15 13
6 9 5 4 20 16 5 19
7 7 17 16 11 16 16 7 18
10 4 1 11 11 11 2 17 5 9 20 5
7 10 10 16 3 3 17 2 3
4 3 2 10 1 15
6 3 5 15 12 17 13 17
9 1 19 3 17 20 3 14 7 10 18
10 7 16 13 20 19 8 1 1 6 10 17
10 5 11 3 16 9 10 14 13 13 2 6
3 4 10 11 2
1 8 14
3 8 20 3 5
6 7 2 20 15 13 15 2
2 8 5 1
1 10 20
3 6 4 18 12
4 7 16 4 2 20
8 10 11 4 20 10 5 13 10 4
9 4 2 13 15 12 7 15 12 3 2
1 8 9
1 9 19
10 4 8 3 17 17 14 17 10 4 5 14
10 7 3 4 14 3 4 14 5 1 15 14
7 1 16 11 9 3 12 3 4
6 1 12 12 6 1 8 12
2 10 5 7
1 4 4
1 5 12
1 10 8
3 3 15 4 16
6 5 5 1 7 12 11 16
5 5 18 11 6 19 3
2 9 19 10
3 7 5 5 8
6 9 8 8 6 10 12 14
1 3 20
1 7 3
2 3 14 10
9 7 5 19 14 10 12 3 8 15 12
9 1 13 14 1 14 11 15 7 12 10
`

type testCase struct {
	n int
	x int
	a []int
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

func referenceSolve(n, x int, arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	if sum%x != 0 {
		return n
	}
	left := -1
	for i, v := range arr {
		if v%x != 0 {
			left = i
			break
		}
	}
	if left == -1 {
		return -1
	}
	right := -1
	for i := n - 1; i >= 0; i-- {
		if arr[i]%x != 0 {
			right = i
			break
		}
	}
	removeLen := left + 1
	if n-right < removeLen {
		removeLen = n - right
	}
	return n - removeLen
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcases))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("empty testcases")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing n for case %d", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n case %d: %w", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing x for case %d", i+1)
		}
		x, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse x case %d: %w", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("missing value %d case %d", j+1, i+1)
			}
			v, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("parse value %d case %d: %w", j+1, i+1, err)
			}
			arr[j] = v
		}
		cases = append(cases, testCase{n: n, x: x, a: arr})
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.x))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		want := referenceSolve(tc.n, tc.x, tc.a)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		g, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil {
			fmt.Printf("case %d failed: non-integer output %q\n", idx+1, got)
			os.Exit(1)
		}
		if g != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%d\ngot:\n%d\n", idx+1, input, want, g)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

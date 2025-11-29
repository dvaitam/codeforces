package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	maxN = 500000
	B    = 710
)

// Embedded source for the reference solution (was 1207F.go).
const solutionSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxN = 500000
	B    = 710 // ~sqrt(maxN)
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	arr := make([]int, maxN+1)
	small := make([][]int, B)
	for i := 0; i < B; i++ {
		small[i] = make([]int, B)
	}

	for ; q > 0; q-- {
		var t, x, y int
		fmt.Fscan(reader, &t, &x, &y)
		if t == 1 {
			arr[x] += y
			for m := 1; m < B; m++ {
				small[m][x%m] += y
			}
		} else {
			if x < B {
				fmt.Fprintln(writer, small[x][y])
			} else {
				sum := 0
				start := y
				if start == 0 {
					start = x
				}
				for j := start; j <= maxN; j += x {
					sum += arr[j]
				}
				fmt.Fprintln(writer, sum)
			}
		}
	}
}
`

const testcasesRaw = `5 2 6 5 1 8 2 1 3 1 2 8 2 2 9 1
5 1 1 2 2 5 2 2 3 1 1 10 5 2 3 2
1 1 4 2
2 1 5 3 1 9 2
2 1 7 3 1 6 4
2 1 5 1 2 5 5
5 1 10 3 1 5 3 2 8 3 1 8 4 1 1 3
1 2 7 1
5 2 6 4 1 8 1 1 10 2 1 4 4 2 9 3
5 2 8 1 2 5 1 2 2 2 2 9 5 2 3 3
3 1 5 3 2 3 1 1 5 4
2 1 2 5 2 1 2
5 2 5 4 2 3 1 1 8 3 1 3 5 1 7 1
2 2 6 2 1 7 3
2 2 10 2 2 8 3
4 2 5 4 2 3 1 2 9 2 2 6 2
1 2 5 5
5 2 2 3 1 5 3 2 8 3 2 6 2 1 8 5
3 2 5 4 2 9 3 2 5 3
4 2 3 4 2 6 5 1 9 2 1 6 4
3 1 7 2 2 5 5 2 1 2
2 2 10 2 1 3 1
4 1 3 1 1 2 3 1 8 2 1 7 4
3 2 10 1 1 4 3 1 6 4
3 2 2 5 2 1 5 2 2 3
5 2 10 3 2 3 4 2 10 5 2 8 2 1 10 4
5 2 4 2 1 6 1 2 2 3 1 6 5 2 7 4
2 2 5 4 2 7 2
5 2 5 4 2 7 1 2 5 4 2 3 4 1 2 5
4 1 5 1 1 7 1 2 9 5 2 4 4
1 1 8 3
2 2 5 4 2 9 5
1 1 3 3
3 2 10 3 1 8 3 2 10 2
1 1 5 5
4 1 9 2 1 7 4 2 7 4 2 4 3
4 1 5 1 2 10 3 2 5 2 1 8 5
4 2 9 2 2 10 5 1 2 2 2 2 1
1 2 7 1
4 2 1 1 1 4 5 1 3 3 2 3 2
5 1 7 2 1 10 2 1 9 1 2 2 3 1 10 5
1 2 10 2
1 2 2 3
5 2 8 3 2 1 2 2 3 5 2 9 5 1 5 1
5 1 5 3 1 4 3 2 7 2 1 1 3 1 10 1
1 1 9 5
4 1 7 4 2 6 1 1 9 4 1 3 4
2 1 2 4 1 3 5
4 2 4 3 2 10 3 2 7 2 2 5 4
3 1 10 5 1 3 4 1 2 1
1 1 4 2
1 2 8 3
1 2 8 3
5 2 6 4 2 2 2 1 3 1 2 7 4 1 9 3
1 2 3 3
2 1 8 1 2 6 3
1 1 8 4
2 1 6 3 2 4 2
1 1 1 5
2 1 10 4 2 5 2
4 2 5 1 2 3 2 1 4 1 2 5 1
4 1 8 4 1 1 1 2 2 1 1 8 1
4 1 4 1 2 7 1 1 5 4 2 5 1
1 1 3 5
4 1 9 2 1 6 1 1 9 2 1 4 5
1 2 9 3
4 2 10 4 1 1 3 2 4 5 1 3 3
4 1 5 4 2 3 2 1 8 2 1 7 5
2 2 1 4 1 1 2
3 2 2 1 2 6 4 1 5 2
5 1 8 2 1 1 4 2 2 4 1 6 1 2 3 2
3 2 6 3 1 2 5 1 1 1
1 1 6 4
3 1 4 5 2 3 3 1 6 3
1 1 10 5
2 1 5 5 1 3 4
4 2 2 3 2 8 2 2 10 4 1 8 1
2 2 8 5 1 2 2
4 1 4 1 1 10 2 2 6 4 2 4 2
3 1 5 5 2 4 3 2 4 1
5 2 7 3 1 3 2 2 7 1 1 3 4 1 5 2
3 2 8 5 1 4 4 1 2 5
5 1 7 1 2 1 5 2 9 4 2 9 4 2 10 5
1 1 7 3
5 1 9 2 2 5 1 1 2 3 2 5 1 2 4 1
1 1 5 3
5 1 8 5 1 2 2 2 9 4 2 1 4 2 3 5
2 2 7 3 1 9 3
2 1 4 3 2 5 5
1 1 2 4
2 1 2 5 2 7 2
4 2 9 5 1 5 4 2 2 5 2 3 4
1 1 2 3
5 1 7 2 2 10 1 2 2 1 2 3 5 1 6 2
4 2 4 5 2 2 1 1 3 1 1 6 2
3 1 10 4 2 3 5 1 6 1
1 1 4 3
3 1 8 1 1 7 5 2 7 1
1 2 3 1
4 2 6 1 2 5 1 1 8 3 1 6 4
`

var _ = solutionSource

type testCase struct {
	q   int
	ops [][3]int
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("invalid line %d", idx+1)
		}
		q, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+3*q {
			return nil, fmt.Errorf("mismatched query count on line %d", idx+1)
		}
		ops := make([][3]int, q)
		for i := 0; i < q; i++ {
			t, _ := strconv.Atoi(fields[1+3*i])
			x, _ := strconv.Atoi(fields[1+3*i+1])
			y, _ := strconv.Atoi(fields[1+3*i+2])
			ops[i] = [3]int{t, x, y}
		}
		cases = append(cases, testCase{q: q, ops: ops})
	}
	return cases, nil
}

func expected(tc testCase) string {
	arr := make([]int, maxN+1)
	small := make([][]int, B)
	for i := 0; i < B; i++ {
		small[i] = make([]int, B)
	}
	var out strings.Builder
	for _, op := range tc.ops {
		t, x, y := op[0], op[1], op[2]
		if t == 1 {
			arr[x] += y
			for m := 1; m < B; m++ {
				small[m][x%m] += y
			}
		} else {
			if x < B {
				fmt.Fprintln(&out, small[x][y])
			} else {
				sum := 0
				start := y
				if start == 0 {
					start = x
				}
				for j := start; j <= maxN; j += x {
					sum += arr[j]
				}
				fmt.Fprintln(&out, sum)
			}
		}
	}
	return strings.TrimSpace(out.String())
}

func run(bin string, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.q))
		sb.WriteByte('\n')
		for _, op := range tc.ops {
			fmt.Fprintf(&sb, "%d %d %d\n", op[0], op[1], op[2])
		}
		input := sb.String()
		want := expected(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed\nexpected:\n%s\nGot:\n%s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

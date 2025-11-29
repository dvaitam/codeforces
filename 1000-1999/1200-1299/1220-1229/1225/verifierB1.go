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
1 9 1 4
10 8 3 1 7 4 6 2 4 7 4 8 2
7 5 5 4 1 3 5 4 3 1
3 4 2 2 3 4
4 5 1 4 5 3 5
8 9 4 2 1 2 3 3 3 9 4
5 6 5 5 3 3 3 3
2 5 1 5 4
3 10 3 2 6 1
7 2 4 1 1 2 1 2 1 1
10 2 5 2 2 1 2 2 1 1 2 1 1
2 7 1 7 7
1 4 1 4
3 2 2 1 1 1
2 7 2 7 5
5 9 3 8 6 2 4 6
1 1 1 1
10 6 8 4 3 4 1 1 3 5 4 1 3
4 10 4 6 5 3 9
4 5 2 2 3 1 3
2 8 1 6 4
7 5 1 3 2 3 5 3 2 3
2 9 1 4 4
1 4 1 1
5 9 1 2 1 1 5 6
8 8 3 2 6 2 3 3 3 3 6
5 2 5 2 1 1 1 1
6 10 6 9 4 3 5 7 9
3 1 3 1 1 1
8 7 5 5 4 7 5 4 1 4 7
6 3 3 2 1 3 2 3 1
1 6 1 5
3 3 2 2 2 3
7 3 5 1 1 2 1 1 3 2
9 8 4 4 6 8 8 4 7 6 5 4
1 2 1 1
9 4 5 3 3 3 2 4 1 1 4 2
3 5 2 2 5 1
8 7 6 4 5 7 2 5 6 1 5
2 5 1 3 1
3 10 3 2 8 4
7 7 4 2 3 4 2 5 4 2
2 7 2 1 6
5 5 2 4 5 1 2 5
8 10 1 1 10 4 5 4 3 5 3
9 4 5 3 3 4 2 3 4 4 1 2
10 7 4 3 7 1 7 1 1 5 6 1 5
5 3 1 3 2 3 2 2
9 6 9 3 1 1 4 6 4 3 3 5
7 6 7 6 6 5 4 1 6 4
7 4 5 1 3 2 4 4 3 2
8 10 4 6 9 1 7 10 7 7 6
10 10 2 8 4 5 1 7 3 7 5 3 2
10 1 6 1 1 1 1 1 1 1 1 1 1
5 9 1 7 2 6 2 8
1 3 1 3
2 7 2 5 3
4 9 2 4 6 5 2
2 9 2 8 9
9 1 3 1 1 1 1 1 1 1 1 1
10 6 4 3 5 6 2 6 1 5 4 3 4
4 5 2 1 2 5 4
10 3 10 2 2 3 1 1 1 3 2 2 2
7 4 1 2 3 1 1 2 4 3
8 2 3 1 1 1 1 1 2 2 2
5 2 5 1 1 1 2 1
8 8 7 3 4 4 5 8 7 4 8
5 6 4 5 1 2 1 1
1 1 1 1
7 10 3 4 7 3 3 1 1 7
3 9 1 7 5 3
2 8 2 1 1
9 1 9 1 1 1 1 1 1 1 1 1
3 5 3 2 4 4
6 5 3 2 2 1 5 5 2
6 7 5 6 5 6 5 1 3
9 7 9 2 6 5 4 6 1 6 3 6
10 2 5 1 1 1 1 1 2 1 1 1 2
9 6 2 3 1 2 5 1 4 6 2 4
8 1 5 1 1 1 1 1 1 1 1
5 6 2 3 4 1 6 3
2 7 1 5 5
4 6 3 5 4 5 4
2 3 2 3 3
10 9 9 1 5 3 4 6 7 9 6 2 7
6 3 5 1 1 2 3 3 2
7 5 3 3 3 3 5 5 1 5
2 3 2 3 2
6 10 1 8 5 8 8 6 7
2 10 1 3 1
9 8 5 4 6 6 6 7 5 8 6 3
1 3 1 3
4 10 2 2 3 7 10
1 2 1 1
4 5 1 5 5 1 1
4 3 4 1 3 2 2
5 4 2 4 2 4 4 3
9 4 8 1 3 4 2 1 4 4 1 4
10 9 10 7 1 6 8 1 4 5 1 9 2
5 9 3 9 9 5 9 7
`

type testCase struct {
	n   int
	k   int
	d   int
	arr []int
}

func solve(tc testCase) int {
	n, d := tc.n, tc.d
	a := tc.arr
	cnt := make(map[int]int)
	unique := 0
	for i := 0; i < d; i++ {
		if cnt[a[i]] == 0 {
			unique++
		}
		cnt[a[i]]++
	}
	best := unique
	for i := d; i < n; i++ {
		if cnt[a[i]] == 0 {
			unique++
		}
		cnt[a[i]]++
		rem := a[i-d]
		cnt[rem]--
		if cnt[rem] == 0 {
			unique--
		}
		if unique < best {
			best = unique
		}
	}
	return best
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.d))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	want := solve(tc)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid testcase: %q", line)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		d, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, err
		}
		if len(parts) != 3+n {
			return nil, fmt.Errorf("expected %d numbers got %d", n+3, len(parts))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[3+i])
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		tests = append(tests, testCase{n: n, k: k, d: d, arr: arr})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
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

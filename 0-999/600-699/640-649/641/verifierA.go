package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	n     int
	dirs  string
	jumps []int
}

// Embedded testcases from testcasesA.txt.
const testcasesRaw = `7 ><>>>>> 3 5 2 5 2 3 2
2 >< 3 1
2 >> 5 1
6 >><>>> 1 5 1 1 4 1
8 ><><<<<< 5 4 1 1 3 5 4 1
5 ><><> 4 1 5 4 3
4 ><<< 1 5 3 4
2 << 2 1
2 >> 5 2
4 >>>> 3 1 3 5
2 >> 2 2
1 > 1
4 ><>> 1 1 2 2
1 < 1
2 << 4 1
6 <<<<<< 4 2 1 1 5 4
2 >< 2 1
5 >><<> 1 5 1 4 2
5 >><<< 2 2 3 5 3
2 >< 1 4
7 >>>><<> 1 3 1 5 3 2 2
8 >>><>>>< 1 5 2 3 2 2 2 4
7 ><>><<> 1 3 2 4 5 4 5
1 < 4
6 >><><< 2 1 4 4 3 1
4 <<<< 1 5 2 3
5 <<>>< 1 3 4 1 3
3 ><< 3 1 1
1 < 3
6 ><>>>> 5 2 2 4 5 3
1 < 2
5 >>><> 5 1 1 3 2
3 >>> 5 2 3
2 >< 1 3
3 <>> 3 3 4
2 <> 4 3
6 <><>>< 3 3 2 2 5 4
2 << 2 2
1 > 1
2 >> 4 4
4 ><>< 3 5 2 4
4 ><<< 5 4 2 1
8 >><<<<<< 4 4 3 5 2 1 5 4
3 >>> 3 4 4
4 <<>> 3 1 5 2
5 <>>>< 1 5 1 1 2
3 <>< 4 3 2
3 >>> 5 5 1
2 <> 2 3
7 >><<<<> 5 5 3 3 1 4 3
5 >>><< 2 2 3 3 1
1 > 4
3 ><< 3 4 1
8 >><<><<< 5 5 2 3 1 5 3 2
7 ><<>><> 2 4 3 1 5 2 1
7 >><>><< 1 4 3 1 5 5 5
4 <<>> 1 2 4 5
7 <<><<>< 1 4 4 4 1 4 3
5 <><<> 1 3 3 2 1
4 ><<< 1 2 1 1
5 ><<<< 4 1 4 3 3
3 <<> 3 4 3
5 ><<<> 2 4 2 2 2
6 <<<>>> 1 2 5 1 4 1
2 <> 3 5
7 <>>><<> 3 5 1 4 4 1 4
6 ><>>>< 2 1 3 1 5 5
3 >>< 4 4 2
4 >><> 4 1 4 2
5 <>><< 3 1 3 5 5
3 >>< 4 2 5
7 ><<><<> 4 5 5 4 4 1 3
6 ><<<<> 3 5 3 1 5 2
3 <>> 3 5 2
4 <>>> 3 4 3 5
3 <<< 4 4 5
8 <>>>><>> 5 3 4 1 4 3 2 4
1 < 1
6 >><<<> 1 3 1 1 5 1
5 ><<>< 1 4 4 3 4
3 >>> 2 4 2
6 <<<>>> 4 4 4 2 2 4
4 <><< 1 2 3 1
3 <>< 5 3 3
7 ><>>>>> 2 3 3 1 1 1 2
6 <><<<> 2 5 4 5 2 4
4 ><<> 3 5 4 3
5 <<<>< 2 5 3 2 2
5 >>>>> 3 2 1 3 5
2 <> 4 4
1 > 4
8 ><<<<<<> 2 4 5 1 4 5 4 1
3 <>< 2 3 3
6 ><><>> 5 5 4 5 3 3
4 <<<< 4 2 2 3
1 > 1
2 >> 1 5
6 <><<<> 3 2 3 4 3 5
3 <<> 3 5 1
3 ><< 5 1 2`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testcase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		dirs := fields[1]
		if len(dirs) != n {
			return nil, fmt.Errorf("line %d: dir length mismatch", idx+1)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d jumps got %d", idx+1, n, len(fields)-2)
		}
		jumps := make([]int, n)
		for i := 0; i < n; i++ {
			jumps[i], err = strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad jump", idx+1)
			}
		}
		res = append(res, testcase{n: n, dirs: dirs, jumps: jumps})
	}
	return res, nil
}

// Embedded solver logic from 641A.go.
func solve(n int, dirs string, jumps []int) string {
	visited := make([]bool, n)
	pos := 0
	for {
		if pos < 0 || pos >= n {
			return "FINITE"
		}
		if visited[pos] {
			return "INFINITE"
		}
		visited[pos] = true
		if dirs[pos] == '>' {
			pos += jumps[pos]
		} else {
			pos -= jumps[pos]
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.n, tc.dirs, tc.jumps)

		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		sb.WriteString(tc.dirs)
		sb.WriteByte('\n')
		for i, v := range tc.jumps {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

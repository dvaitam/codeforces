package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test data from testcasesD.txt.
const testData = `100
3
18 17 4
5
19 15 20 18 2
9
0 15 8 17 7 6 15 17 17
7
12 20 4 7 20 4 16
7
0 2 5 18 1 9 0
5
15 19 12 13 12
9
14 4 11 3 1 4 15 6 8
7
20 9 13 16 12 18 11
9
18 13 18 7 10 0 8 19 5
5
17 18 18 3 20
3
20 18 8
5
3 2 15 20 15
1
11
1
13
3
0 9 13
7
3 1 19 19 1 12 18
5
17 8 16 7 1
5
0 2 3 19 17
1
6
7
9 19 8 4 1 10 10
5
4 12 12 14 16
7
20 19 17 3 19 16 8
7
20 7 9 13 8 16 9
9
10 0 13 18 10 0 12 19 18
3
1 20 20
5
14 11 11 19 8
7
0 18 1 0 11 8 20
7
9 18 19 10 5 11 5
5
11 19 8 9 12
1
0
9
4 9 16 7 20 8 7 10 5
7
20 3 3 19 10 10 7
7
5 2 10 20 6 18 14
5
7 3 1 16 6
5
18 5 8 10 20
1
19
5
18 4 13 9 16
5
14 11 20 13 9
7
18 13 1 13 4 6 0
7
19 16 13 17 7 1 14
9
9 17 10 7 2 18 9 3 7
1
1
9
6 13 18 1 0 15 3 5 16
5
7 0 16 17 13
1
19
1
10
3
8 17 15
1
11
3
6 3 17
1
5
3
8 4 0
7
20 18 12 1 8 7 8
9
16 16 13 1 15 10 0 1 4
1
3
1
2
7
1 2 16 16 15 10 5
5
2 11 12 20 12
9
9 11 8 6 10 13 3 4 17
1
12
1
18
3
1 11 14
9
20 17 12 20 1 19 13 1 11
7
10 13 13 14 0 7 6
9
8 18 2 13 7 13 4 0 10
5
17 8 3 14 3
9
12 3 10 18 17 3 18 0 15
3
7 12 1
9
2 18 3 12 5 0 10 3 0
1
15
5
18 9 2 1 18
9
16 7 3 17 3 17 1 17 10
9
5 2 7 5 20 7 14 19 12
5
11 19 12 11 17
7
2 12 16 7 13 5 13
9
18 16 15 4 20 12 4 5 3
7
15 16 14 18 5 4 8
3
4 18 16
5
7 17 9 13 19
9
18 8 6 9 0 8 15 12 6
3
18 11 7
5
15 4 13 15 19
3
14 18 20
9
0 15 2 12 1 14 7 7 20
1
6
5
7 6 8 4 5
9
1 8 5 1 10 5 13 2 2
1
2
5
9 1 11 14 18
5
0 0 10 10 13
7
15 2 6 20 18 15 12
3
17 10 3
5
2 13 3 14 16
5
3 16 11 11 14
5
20 8 3 10 18
9
16 3 15 16 11 1 9 18 5
3
5 11 20
7
3 3 17 4 10 20 20
9
13 17 9 20 5 14 15 9 5
1
3
3
17 17 18`

func runCandidate(bin string, input string) (string, error) {
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

func solveCase(a []int64) int64 {
	n := len(a)
	if n == 1 {
		return a[0]
	}
	even := make([]int64, 2*n+1)
	odd := make([]int64, 2*n+1)
	for i := 0; i < 2*n; i++ {
		v := a[i%n]
		if i%2 == 0 {
			even[i+1] = even[i] + v
			odd[i+1] = odd[i]
		} else {
			odd[i+1] = odd[i] + v
			even[i+1] = even[i]
		}
	}
	k := n / 2
	var best int64
	for i := 0; i < n; i++ {
		var s int64
		if i%2 == 0 {
			s = even[i+2*k+1] - even[i]
		} else {
			s = odd[i+2*k+1] - odd[i]
		}
		if s > best {
			best = s
		}
	}
	return best
}

func parseTests() ([]([]int64), error) {
	tokens := strings.Fields(testData)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("invalid test count")
	}
	idx++
	tests := make([][]int64, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("unexpected end of data at case %d", caseIdx+1)
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid n at case %d", caseIdx+1)
		}
		idx++
		if idx+n > len(tokens) {
			return nil, fmt.Errorf("not enough numbers for case %d", caseIdx+1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(tokens[idx+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value in case %d", caseIdx+1)
			}
			arr[i] = v
		}
		idx += n
		tests = append(tests, arr)
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, arr := range tests {
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", len(arr))
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		expect := solveCase(arr)
		out, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Printf("case %d: invalid output\n", i+1)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

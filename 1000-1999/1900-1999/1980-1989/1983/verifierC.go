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
100
4 1 9 5 2 1 8 5 3 3 5 7 2
3 1 1 2 1 2 1 2 1 1
1 5 5 5
1 1 1 1
4 1 1 1 1 1 1 1 1 1 1 1 1
5 1 2 9 1 1 1 1 10 1 1 10 1 1 1 1
2 1 1 1 1 1 1
5 1 1 2 8 1 1 3 3 2 4 1 1 6 1 4
1 1 1 1
2 1 2 1 2 2 1
2 1 3 2 2 1 3
1 5 5 5
4 2 1 1 3 1 2 1 3 1 4 1 1
4 1 1 6 1 2 5 1 1 1 4 3 1
3 1 1 7 4 4 1 2 2 5
4 10 1 2 2 6 1 1 7 1 1 12 1
4 7 7 2 1 2 5 3 7 2 1 5 9
3 7 2 1 3 6 1 1 3 6
1 5 5 5
5 1 2 1 1 3 1 1 4 1 1 1 1 3 2 1
1 2 2 2
3 3 6 2 8 2 1 3 2 6
5 3 1 3 4 3 1 1 6 5 1 4 1 1 2 6
1 3 3 3
2 2 8 7 3 5 5
1 1 1 1
5 1 1 1 16 1 14 2 2 1 1 13 1 3 1 2
1 1 1 1
4 1 10 1 8 3 12 3 2 6 1 1 12
3 2 1 1 3 2 1 3 2 1
4 1 1 1 1 2 2 2 2 3 3 3 3 4 4 4 4
5 6 3 4 8 7 1 4 8 5 2 4 2 2 2 2
2 4 1 1 2 3 1
1 3 3 3
1 1 1 1
5 2 5 8 8 2 2 5 11 1 1 9 1 3 1 6
3 5 2 3 2 2 6 2 5 7
1 2 2 2
3 1 2 3 4 5 6 7 8 9
4 1 2 1 1 3 3 2 1 2 4 4 6 2
3 10 1 2 6 1 1 5 8 1 1
3 11 1 8 1 1 2 2 1 5 1
2 1 4 2 2 3 3
5 2 1 3 1 2 1 4 1 2 1 2 1 6 1 2
5 2 1 1 1 1 8 1 1 9 1 1 6 1 1 3
5 1 2 7 1 1 2 3 3 4 5 1 2 3 1 3 6
4 1 1 2 2 3 4 1 1 2 2 1 1 1
4 1 1 1 2 2 2 1 1 1 1 1 1
4 3 1 3 1 1 1 2 3 1 4 3 1 4
3 1 1 3 5 4 1 2 6
2 2 4 2 1 3 4
4 7 5 2 1 1 9 1 3 2 2 2 2 2 1 8 1
1 3 3 3
5 1 1 4 4 1 1 2 2 4 3 1 5 1 1 3
5 1 3 6 1 1 6 4 3 6 2 1 4 1 3 1
2 1 3 1 2 1 1
1 10 10 10
1 2 2 2
3 4 1 1 2 3 4 5 1 2 3
1 2 2 2
4 1 2 1 2 3 3 3 4 2 1 1 1
5 2 1 1 1 1 4 6 3 1 7 6 2 1 4 7 3
1 2 2 2
2 1 1 1 1 1 1
5 1 1 1 6 2 4 4 1 1 6 4 1 1 6 5
1 1 1 1
5 1 1 2 3 4 1 1 3 2 5 1 1 7 1
5 8 2 1 1 4 2 3 6 5 4 3 6 1 1 4 4
5 5 1 5 2 5 2 2 7 2 1 8 1 1 3 9
3 8 1 2 2 2 1 2 2 8 8
5 3 3 2 3 1 4 2 4 1 1 2 5 1 2 1
5 3 2 5 1 1 2 5 2 2 2 2 4 1 3 4
4 1 1 1 1 1 1 4 1 3 2 1 2
3 10 7 6 4 6 1 1 6 4 1
4 9 3 1 2 2 2 2 1 1 1 2 2 1
2 3 8 3 1 8 1
4 2 1 1 1 1 2 1 2 1 2 3 2 3
5 7 4 1 2 2 1 3 1 2 2 1 3 1 3 3
1 2 2 2
1 4 4 4
2 1 1 1 3 4 3
1 7 7 7
5 1 1 1 1 2 2 2 2 1 3 3 3 3 4 4 4
5 1 1 4 1 3 8 7 1 7 6 1 1 5 6 1
4 3 1 5 6 2 8 3 8 2 1 1 3
5 9 1 2 9 1 1 7 2 1 2 1 2 2 1
5 3 1 1 1 4 3 1 1 2 7 1 7 2
4 3 2 2 2 3 3 3 1 1 1 3 4
3 1 1 1 1 1 2 2 2 3
2 1 1 2 5 3 3
4 3 1 4 2 1 5 2 6 2 8 1 1 6
3 2 1 1 2 6 1 1 4 4
2 5 1 1 2 4 3
5 2 3 4 1 1 6 6 6 1 1 8 5 2 1 4
1 2 2 2
5 5 1 3 1 4 2 2 4 4 2 1 5 3 1
5 4 1 1 2 5 2 2 6 3 3 2 2 2 2
3 6 8 3 1 4 8 7 1 1
5 3 1 1 1 2 6 1 1 2 2 6 1 3 1
2 3 1 1 2 3 1
5 4 2 1 1 3 3 3 2 1 1 4 3 2 1 4
3 1 1 7 2 2 1 1 8 4
2 1 1 1 2 3 4
2 1 1 4 6 3 1
5 1 1 1 6 4 4 3 3 1 3 4 6 2 1
1 1 1 1
3 1 1 1 2 2 2 3 3 3
4 1 1 2 2 5 5 1 4 5 5 1 3
1 2 2 2
1 3 3 3
3 1 1 2 1 1 4 2 1 1
1 1 1 1
5 1 1 5 3 5 1 3 1 3 1 3 1 3 1 3
3 1 1 7 7 1 1 3 4
3 4 1 2 2 4 2 1 3 2 1
1 2 2 2
4 4 2 4 1 1 2 4 2 4 3 4 1
4 4 2 1 2 3 3 3 4 4 1 1 1
1 5 5 5
5 1 1 5 3 1 1 4 2 4 6 5 2 6 3 4
2 1 1 1 1 1 1
3 2 3 5 6 3 1 1 3 6
2 2 4 6 2 2 1
1 1 1 1
1 1 1 1
3 2 1 2 2 3 2 2 2
1 1 1 1
`

type testCase struct {
	n   int
	arr [3][]int64
}

func parseTests(raw string) ([]testCase, error) {
	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty testcases")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for c := 0; c < t; c++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("truncated before case %d", c+1)
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("bad n in case %d", c+1)
		}
		idx++
		var arr [3][]int64
		for k := 0; k < 3; k++ {
			arr[k] = make([]int64, n)
			for i := 0; i < n; i++ {
				if idx >= len(tokens) {
					return nil, fmt.Errorf("missing values in case %d", c+1)
				}
				v, err := strconv.ParseInt(tokens[idx], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("bad value in case %d", c+1)
				}
				arr[k][i] = v
				idx++
			}
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	return cases, nil
}

func solve(tc testCase) string {
	n := tc.n
	arr := tc.arr
	need := int64(0)
	for _, v := range arr[0] {
		need += v
	}
	need = (need + 2) / 3
	perms := [6][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	var seg [3][2]int
	for _, p := range perms {
		ptr := 0
		ok := true
		for i := range seg {
			seg[i][0], seg[i][1] = 0, -1
		}
		for _, x := range p {
			start := ptr
			cur := int64(0)
			for ptr < n && cur < need {
				cur += arr[x][ptr]
				ptr++
			}
			if cur < need {
				ok = false
				break
			}
			seg[x][0] = start
			seg[x][1] = ptr - 1
		}
		if ok {
			// The original solution prints a rune newline with fmt.Fprint,
			// which renders as numeric 10 after the segment pairs.
			return fmt.Sprintf("%d %d %d %d %d %d 10",
				seg[0][0]+1, seg[0][1]+1,
				seg[1][0]+1, seg[1][1]+1,
				seg[2][0]+1, seg[2][1]+1)
		}
	}
	return "-1"
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for k := 0; k < 3; k++ {
		for i, v := range tc.arr[k] {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		want := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

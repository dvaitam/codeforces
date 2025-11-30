package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	input    string
	expected string
}

type fig struct{ cnt, cost int64 }

const testcaseData = `100
1
2 2
3
3 8 13
5
4 19
1 18
3 13
7 16
6 17
4
9 14 15 16
3
8 10
7 13
9 5
5
3 7 11 12 15
3
3 4
9 16
6 16
5
3 11 18 27 33
5
6 11
8 5
7 14
9 7
8 8
4
9 18 24 32
4
6 18
9 14
8 7
6 5
5
5 13 18 23 32
5
9 16
10 18
7 9
4 15
9 11
5
2 8 9 13 15
1
10 20
1
5
5
4 3
9 4
5 7
4 1
7 1
1
6
3
3 7
1 2
2 2
1
1
1
6 8
2
3 6
5
1 12
10 1
4 4
1 0
6 19
1
5
3
8 0
5 14
9 19
1
5
4
10 4
8 7
2 10
2 0
4
3 12 22 29
4
9 10
3 10
5 8
10 13
1
9
2
1 8
1 4
2
3 5
4
4 16
1 7
4 14
2 8
1
10
2
10 19
6 8
4
5 14 15 18
1
7 13
2
2 11
1
4 3
1
1
2
4 3
4 0
5
8 16 21 30 37
2
4 13
7 16
1
10
5
1 13
9 18
3 3
8 11
1 16
1
10
3
5 11
5 0
7 3
1
5
2
1 14
1 13
4
8 12 22 32
1
1 9
1
6
3
2 7
8 6
2 18
3
7 15 18
3
7 3
5 3
2 2
5
6 13 17 19 20
5
8 1
8 9
6 14
3 11
5 15
5
8 15 23 28 35
2
3 15
10 8
5
7 9 19 29 31
1
6 5
5
3 10 12 14 15
2
5 12
4 10
4
3 12 17 19
2
9 13
2 10
5
4 13 18 21 24
4
4 12
6 18
3 14
8 0
5
7 10 17 26 27
4
5 12
5 13
8 11
9 10
1
4
5
10 6
7 12
1 10
8 16
8 20
2
2 3
4
4 18
10 12
4 3
7 17
2
5 15
5
4 15
10 4
1 19
7 15
5 16
5
3 11 15 17 23
1
8 17
1
10
4
6 14
5 16
8 0
2 19
3
3 10 15
2
1 5
8 12
4
5 8 9 14
5
8 0
6 1
9 12
10 14
4 9
4
3 11 20 25
1
5 10
3
6 11 18
5
2 16
4 12
10 16
3 16
2 9
1
4
4
9 7
9 8
1 3
2 12
3
4 10 16
1
6 14
3
3 11 19
3
8 4
8 20
4 8
3
3 5 9
4
4 11
3 11
3 4
4 8
5
7 14 20 25 35
5
10 10
7 20
5 17
10 20
2 11
3
7 15 18
3
6 14
8 2
3 10
4
3 4 6 12
2
6 2
7 0
5
6 10 20 27 36
3
8 20
3 11
6 6
4
2 5 9 15
3
3 13
6 8
2 10
2
4 8
5
1 10
6 20
10 1
3 5
2 13
4
5 8 14 23
5
2 10
10 12
4 1
7 15
8 19
3
9 19 29
1
10 16
5
8 15 23 26 33
4
9 14
1 3
8 18
3 3
5
3 5 12 17 25
1
5 3
3
4 7 8
2
7 2
6 20
4
1 9 13 15
4
3 17
1 4
9 17
1 1
2
9 10
5
6 16
4 4
6 15
1 4
9 3
2
2 10
2
1 19
4 20
4
6 16 23 32
5
3 16
2 4
4 5
7 6
5 10
4
3 10 13 20
3
5 3
9 3
8 8
3
9 17 22
2
7 4
9 3
1
10
5
4 6
4 12
10 1
3 20
1 8
4
9 10 14 17
5
6 1
4 3
3 20
9 5
2 14
3
4 7 13
3
9 18
2 13
7 1
4
5 7 12 13
2
7 10
5 17
4
10 19 23 30
2
3 14
8 11
4
8 18 23 33
2
10 15
8 6
4
10 16 21 23
2
6 19
8 7
5
10 13 18 22 31
3
2 0
1 6
6 1
3
9 14 20
4
2 13
8 0
5 18
10 4
2
3 6
5
7 2
10 14
5 20
2 15
8 7
2
10 15
2
4 19
6 18
5
7 16 23 27 31
5
1 8
4 4
10 12
7 3
8 12
4
8 15 20 24
2
4 1
9 16
1
10
5
1 1
7 13
7 7
9 8
2 11
5
6 15 23 33 35
4
4 8
1 0
8 1
3 20
2
4 10
2
9 1
10 4
3
2 11 20
1
3 13
2
1 6
5
5 15
1 17
6 10
2 19
6 3
5
6 12 17 25 30
5
10 4
1 1
6 13
1 11
9 1
1
9
5
10 13
7 13
4 5
3 19
1 0
5
6 9 14 15 16
2
10 7
7 2
3
2 12 14
2
4 17
4 3
1
7
1
9 8
5
4 5 14 23 32
4
7 4
3 13
3 14
6 1
5
3 12 20 27 37
4
3 15
10 4
6 4
1 8
2
3 10
5
5 14
8 14
4 13
7 8
4 11
1
7
5
1 13
5 0
9 15
10 8
5 7
4
8 14 23 33
4
4 17
9 5
8 9
6 13
1
9
2
7 3
7 19
4
10 19 27 29
4
8 19
6 17
6 5
3 7
2
7 15
4
3 13
5 10
10 12
5 20
3
6 7 14
5
1 6
8 3
2 0
6 10
10 10
4
3 9 11 20
5
8 12
10 7
8 3
10 0
6 0
3
8 11 12
1
6 12
4
10 11 19 29
2
4 19
6 5
3
5 12 22
`

func expectedCase(figs []fig, p []int64) int64 {
	sort.Slice(figs, func(i, j int) bool { return figs[i].cost < figs[j].cost })
	var total int64
	for _, f := range figs {
		total += f.cnt
	}
	seg := make([]int64, 0, len(p)+1)
	var prev int64
	for _, v := range p {
		if prev >= total {
			break
		}
		up := v
		if up > total {
			up = total
		}
		cur := up - prev
		if cur > 0 {
			seg = append(seg, cur)
			prev += cur
		}
	}
	if prev < total {
		seg = append(seg, total-prev)
	}
	var ans int64
	si := 0
	var rem int64
	if len(seg) > 0 {
		rem = seg[0]
	}
	for _, f := range figs {
		cnt := f.cnt
		for cnt > 0 && si < len(seg) {
			take := cnt
			if take > rem {
				take = rem
			}
			ans += take * f.cost * int64(si+1)
			cnt -= take
			rem -= take
			if rem == 0 {
				si++
				if si < len(seg) {
					rem = seg[si]
				}
			}
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: unexpected end of data", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		figs := make([]fig, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d: incomplete figure data", caseNum+1)
			}
			k, err1 := strconv.Atoi(fields[pos])
			c, err2 := strconv.Atoi(fields[pos+1])
			pos += 2
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("case %d: bad figure values", caseNum+1)
			}
			figs[i] = fig{int64(k), int64(c)}
			sb.WriteString(fmt.Sprintf("%d %d\n", k, c))
		}
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing p length", caseNum+1)
		}
		t2, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad p length: %w", caseNum+1, err)
		}
		pos++
		p := make([]int64, t2)
		sb.WriteString(fmt.Sprintf("%d\n", t2))
		for i := 0; i < t2; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: missing p value", caseNum+1)
			}
			v, err := strconv.Atoi(fields[pos])
			pos++
			if err != nil {
				return nil, fmt.Errorf("case %d: bad p value: %w", caseNum+1, err)
			}
			p[i] = int64(v)
			sb.WriteString(fmt.Sprintf("%d", v))
			if i+1 < t2 {
				sb.WriteByte(' ')
			} else {
				sb.WriteByte('\n')
			}
		}
		expected := expectedCase(figs, p)
		cases = append(cases, testCase{input: sb.String(), expected: fmt.Sprintf("%d", expected)})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
100
3 2
cb
cc
cc
1 2
ac
1 1
a
2 2
ab
ca
3 1
a
c
a
2 2
ab
aa
1 3
cba
1 1
a
1 1
a
1 2
ba
3 3
caa
cab
bab
2 1
a
b
1 2
bc
3 1
c
c
c
2 1
b
b
2 2
cb
ab
2 3
aab
acb
2 1
c
b
2 2
ca
ba
3 1
c
a
a
1 2
bc
2 3
bba
ccb
2 1
b
a
1 2
cc
2 1
b
b
3 3
abc
bba
aca
3 3
bba
caa
ccb
1 1
c
3 2
bb
cb
aa
3 1
b
b
a
1 3
cac
2 1
a
b
2 1
a
b
2 1
b
c
1 3
bbc
3 2
bb
bb
ba
1 2
ca
3 2
ba
ab
bc
3 3
bab
ccc
abb
3 3
cbb
bcc
bbc
1 3
abc
2 2
cb
bb
3 3
cbb
bcb
cbb
1 3
bbb
3 1
c
a
a
2 2
bc
ac
3 3
bac
ccc
bbc
3 3
bca
aac
bcc
1 1
c
1 3
cab
1 1
a
1 1
b
1 2
ac
1 2
bb
2 3
cac
aac
2 1
b
b
2 2
ac
cb
1 3
cba
2 3
cbc
bba
2 2
cc
cb
2 1
a
c
2 3
baa
cab
3 1
b
a
b
3 3
cab
ccb
bba
2 2
bc
bb
1 3
cbc
2 2
bb
ab
2 2
ba
ba
1 3
cba
2 1
a
b
1 2
cc
2 1
b
a
1 2
ba
3 2
bb
cb
cc
3 3
aaa
bbc
cbc
2 3
cab
bbc
3 3
aaa
bcb
cac
3 1
a
b
b
3 3
ccc
cbb
bbc
3 1
b
b
a
2 1
c
b
3 2
cb
ba
ab
3 3
bbc
abc
caa
3 1
b
a
a
3 1
b
c
b
1 2
cb
1 1
a
1 3
cac
1 2
cb
1 1
c
1 2
aa
3 1
a
c
a
3 2
ca
ba
cc
1 3
bca
1 2
ab
3 3
bbb
cba
acb
1 3
ccc
2 3
cca
bac
1 2
bc
1 1
b
`

type testCase struct {
	n    int
	m    int
	grid []string
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	idx := 0
	t, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
	if err != nil {
		return nil, err
	}
	idx++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		parts := strings.Fields(strings.TrimSpace(lines[idx]))
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header on case %d", i+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		idx++
		if idx+n > len(lines) {
			return nil, fmt.Errorf("not enough rows for case %d", i+1)
		}
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			grid[r] = strings.TrimSpace(lines[idx+r])
		}
		idx += n
		cases = append(cases, testCase{n: n, m: m, grid: grid})
	}
	return cases, nil
}

// solve mirrors 1185E.go.
func solve(tc testCase) string {
	n, m := tc.n, tc.m
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = int(tc.grid[i][j] - 'a')
		}
	}
	const INF = 1000000000
	mx := -1
	lastR, lastC := 0, 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] > mx {
				mx = grid[i][j]
				lastR, lastC = i, j
			}
		}
	}
	row := make([]int, mx+1)
	col := make([]int, mx+1)
	mnr := make([]int, mx+1)
	mxr := make([]int, mx+1)
	mnc := make([]int, mx+1)
	mxc := make([]int, mx+1)
	for c := 0; c <= mx; c++ {
		row[c] = -1
		col[c] = -1
		mnr[c] = n
		mnc[c] = m
		mxr[c] = -1
		mxc[c] = -1
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c := grid[i][j]
			if row[c] == -1 {
				row[c] = i
			} else if row[c] != i {
				row[c] = INF
			}
			if col[c] == -1 {
				col[c] = j
			} else if col[c] != j {
				col[c] = INF
			}
			if i < mnr[c] {
				mnr[c] = i
			}
			if i > mxr[c] {
				mxr[c] = i
			}
			if j < mnc[c] {
				mnc[c] = j
			}
			if j > mxc[c] {
				mxc[c] = j
			}
		}
	}
	b := make([][]int, n)
	for i := range b {
		b[i] = make([]int, m)
		for j := range b[i] {
			b[i][j] = -1
		}
	}
	ops := make([][4]int, mx+1)
	ok := true
	for c := 0; c <= mx; c++ {
		if row[c] == -1 {
			ops[c] = [4]int{lastR, lastC, lastR, lastC}
		} else if row[c] == INF && col[c] == INF {
			ok = false
			break
		} else if row[c] != INF {
			r := row[c]
			ops[c] = [4]int{r, mnc[c], r, mxc[c]}
			for x := mnc[c]; x <= mxc[c]; x++ {
				b[r][x] = c
			}
		} else {
			cc := col[c]
			ops[c] = [4]int{mnr[c], cc, mxr[c], cc}
			for i := mnr[c]; i <= mxr[c]; i++ {
				b[i][cc] = c
			}
		}
	}
	if ok {
		for i := 0; i < n && ok; i++ {
			for j := 0; j < m; j++ {
				if b[i][j] != grid[i][j] {
					ok = false
					break
				}
			}
		}
	}
	if !ok {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	sb.WriteString(strconv.Itoa(mx + 1))
	sb.WriteByte('\n')
	for c := 0; c <= mx; c++ {
		r1, c1, r2, c2 := ops[c][0]+1, ops[c][1]+1, ops[c][2]+1, ops[c][3]+1
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r1, c1, r2, c2))
	}
	return strings.TrimSpace(sb.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := buildInput(tc)
		expect := solve(tc)
		got, err := run(bin, fmt.Sprintf("1\n%s", input))
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

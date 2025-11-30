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

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `3
cab
1
2 b
8
cbaababb
5
1 c
8 b
4 c
2 b
1 a
1
c
5
1 b
1 b
1 c
1 b
1 c
4
baca
4
3 a
4 c
1 a
3 a
6
cccbcc
2
3 b
5 b
9
bcabacbbc
2
6 c
6 a
8
ccaacbbb
1
8 a
5
ccccb
2
2 c
2 a
4
ccab
5
3 c
3 b
3 c
1 b
2 c
9
ababbccac
4
8 b
7 b
1 c
9 c
10
bbcaacacca
1
9 b
1
c
1
1 a
8
ababacab
3
2 a
3 b
3 c
5
ccbbc
3
4 b
1 a
3 b
6
bababc
5
2 c
4 a
2 a
4 a
1 c
3
bcc
4
3 a
3 c
3 b
1 c
1
b
5
1 c
1 a
1 a
1 a
1 a
2
bb
2
2 c
2 a
1
c
1
1 c
8
acccabab
1
4 c
7
cabacbb
5
4 a
3 c
7 b
3 a
2 a
6
cabbab
1
4 c
6
ccbcaa
1
1 a
3
aca
3
2 c
3 b
2 b
6
abaccb
2
5 c
1 b
1
b
1
1 a
3
bac
5
2 a
3 c
1 c
1 b
2 b
10
cabbaabacc
1
2 b
2
aa
2
2 a
1 b
3
caa
1
2 b
9
bcbcbbaac
3
1 a
1 b
6 b
7
bbaabcb
1
3 a
10
ccbcbbacab
2
4 b
2 b
2
ba
5
2 a
2 b
1 b
1 b
2 a
6
acccca
2
2 a
2 b
2
bc
1
1 a
1
b
3
1 b
1 a
1 a
9
caaaabbac
5
5 a
4 a
9 c
1 b
9 c
4
abbc
2
1 c
2 b
2
cb
4
2 c
2 c
2 a
2 b
3
bba
4
3 a
1 c
2 c
1 c
3
abb
4
3 b
1 c
1 a
2 a
3
cbc
4
3 c
3 a
1 b
2 c
8
acbbcccc
3
4 a
2 c
6 a
9
abbcbcbac
4
2 a
9 c
7 a
3 b
7
accabcb
3
4 c
7 a
5 c
1
c
1
1 c
2
bc
1
1 c
2
ba
4
2 b
1 b
2 a
2 a
2
bc
5
2 a
2 b
1 b
1 a
2 c
1
a
5
1 b
1 a
1 a
1 b
1 c
5
cbacb
4
4 a
2 c
4 a
3 a
1
a
5
1 c
1 c
1 a
1 c
1 b
9
cbcbaabcb
3
5 c
7 b
8 a
7
bacabcc
5
2 b
5 c
4 c
6 b
6 a
8
cccabcac
4
7 b
6 c
2 b
4 c
5
cabcc
2
4 b
2 a
10
abbcbccbab
3
8 a
8 c
1 b
9
accbabacb
1
3 c
3
cab
3
3 b
1 c
1 a
6
baaccc
3
4 c
5 c
1 a
5
ccccb
3
5 c
2 b
5 b
3
bbc
3
3 a
2 c
3 a
1
c
4
1 b
1 b
1 a
1 c
8
ccacbbca
2
3 c
8 b
5
baaca
3
1 a
2 b
3 b
2
aa
1
1 a
1
b
5
1 b
1 a
1 a
1 b
1 b
8
baaabbcc
4
4 b
5 b
8 c
2 a
2
aa
1
2 b
7
cbabaca
1
1 b
3
cca
5
2 b
1 a
2 c
2 a
1 c
1
c
2
1 b
1 b
2
aa
4
1 c
2 c
1 c
2 b
6
cbbcca
2
1 c
5 a
6
bccccc
1
3 c
7
caccbca
3
6 c
6 a
3 a
2
aa
2
2 a
1 c
2
cb
5
2 a
2 a
1 c
1 b
1 b
8
accbabba
3
1 b
1 c
5 b
3
bba
3
1 b
1 c
3 a
6
bcbcba
2
6 b
5 c
10
cccabcaabb
5
6 a
7 b
3 c
2 a
5 c
9
bbbbbbbcc
5
9 a
9 a
3 b
6 b
2 b
5
bbbcb
1
5 a
3
acb
5
2 a
3 c
3 b
2 c
2 b
5
bcbcc
2
1 a
3 c
4
caaa
4
1 a
3 c
1 a
3 a
10
ccaaacacba
5
6 b
5 a
4 c
8 a
7 b
6
cabcab
4
2 a
6 c
4 c
4 a
7
ccccbab
4
1 a
3 c
6 c
1 c
2
bc
3
2 c
2 c
2 c
10
bbbacbcaca
3
1 b
10 a
6 b
7
bccaaaa
4
3 b
3 b
6 c
7 b
6
bbabba
4
2 a
2 b
3 a
5 b
7
bcbcbcb
4
3 b
2 c
7 b
4 c
7
aaaaaac
1
1 b
3
bab
2
1 a
2 c
1
c
2
1 b
1 c
2
cc
4
1 b
2 a
2 c
1 a`

type query struct {
	m int
	c byte
}

type testCase struct {
	n       int
	s       string
	queries []query
}

func solveCase(tc testCase) ([]int, error) {
	n := tc.n
	s := tc.s
	if len(s) != n {
		return nil, fmt.Errorf("string length mismatch: got %d want %d", len(s), n)
	}
	best := make([][]int, 26)
	for i := range best {
		best[i] = make([]int, n+1)
	}
	bytesStr := []byte(s)
	for ch := 0; ch < 26; ch++ {
		prefix := make([]int, n+1)
		c := byte('a' + ch)
		for i := 0; i < n; i++ {
			if bytesStr[i] != c {
				prefix[i+1] = prefix[i] + 1
			} else {
				prefix[i+1] = prefix[i]
			}
		}
		minChanges := make([]int, n+1)
		for i := 1; i <= n; i++ {
			minChanges[i] = n + 1
		}
		for l := 0; l < n; l++ {
			for r := l; r < n; r++ {
				length := r - l + 1
				mism := prefix[r+1] - prefix[l]
				if mism < minChanges[length] {
					minChanges[length] = mism
				}
			}
		}
		for length := 1; length <= n; length++ {
			m := minChanges[length]
			if m <= n {
				for k := m; k <= n; k++ {
					if best[ch][k] < length {
						best[ch][k] = length
					}
				}
			}
		}
		for k := 1; k <= n; k++ {
			if best[ch][k] < best[ch][k-1] {
				best[ch][k] = best[ch][k-1]
			}
		}
	}
	ans := make([]int, len(tc.queries))
	for i, q := range tc.queries {
		m := q.m
		if m > n {
			m = n
		}
		ch := q.c - 'a'
		ans[i] = best[ch][m]
	}
	return ans, nil
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
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
		if !scan.Scan() {
			return nil, fmt.Errorf("missing string for n=%d", n)
		}
		s := scan.Text()
		if !scan.Scan() {
			return nil, fmt.Errorf("missing q")
		}
		qCnt, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse q: %w", err)
		}
		qs := make([]query, qCnt)
		for i := 0; i < qCnt; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing m", len(cases)+1)
			}
			m, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d: parse m: %w", len(cases)+1, err)
			}
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing c", len(cases)+1)
			}
			c := scan.Text()
			if len(c) != 1 {
				return nil, fmt.Errorf("case %d: invalid char %q", len(cases)+1, c)
			}
			qs[i] = query{m: m, c: c[0]}
		}
		cases = append(cases, testCase{n: n, s: s, queries: qs})
	}
	return cases, nil
}

func runCandidate(bin string, input string) ([]string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	fields := strings.Fields(out.String())
	return fields, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expectInts, err := solveCase(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: solve error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect := make([]string, len(expectInts))
		for i, v := range expectInts {
			expect[i] = strconv.Itoa(v)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
		for _, q := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %c\n", q.m, q.c))
		}
		input := sb.String()
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if len(got) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d outputs got %d\nInput:\n%s\n", idx+1, len(expect), len(got), input)
			os.Exit(1)
		}
		for i := range expect {
			if got[i] != expect[i] {
				fmt.Fprintf(os.Stderr, "case %d query %d failed: expected %s got %s\nInput:\n%s\n", idx+1, i+1, expect[i], got[i], input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

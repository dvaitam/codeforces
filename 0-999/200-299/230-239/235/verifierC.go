package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt.
const testcasesCData = `
100
a 1 acc
bcaca 3 bc ccbc cbaa
bbbbca 3 aa a aac
cccabb 3 bcbbb bc ccab
bccbc 2 bccc bcab
cbb 2 ccc ccccb
cabcb 3 abcaa a cabca
ca 2 aa caab
aacaaa 1 a
b 2 ac cc
b 3 a aa ccc
bb 2 b ccca
bccab 1 c
aabacc 2 cbab bcb
c 3 ca aaa ab
ccaa 1 abac
cccb 2 bcaa b
aaccaaa 1 a
aaac 3 bbcc acac bcac
b 3 aacbb c c
bcbbac 2 a aca
abcbbacc 1 b
b 2 a aacb
cbabbab 1 a
cbacaa 3 accb bba bbc
ccbbcbba 1 cbcb
cc 3 a aca aacc
a 2 accc bac
aacba 2 accba bc
bbcc 1 bcac
cabcabb 2 ccb bcbc
ca 3 abcbc b ccbc
aab 1 cbaab
bccc 1 caac
bbccabc 1 b
b 3 c cbbb baacb
bbc 3 aa bbcb ab
abacbcba 3 bca cccb b
bbcbcc 2 accab cacca
aabca 3 aaa babb b
babbbbac 2 bb aa
acbabaaa 2 cbbcb ccc
cbcccc 2 cccab bba
bbbaa 2 aaab ba
acbacbc 2 cabb ba
abb 1 bbab
acac 3 b cca aa
bbabcca 2 baabb cbcc
cc 3 bcba bcba b
acc 1 b
bcaba 3 aaa bc b
abaabaca 1 caaac
c 2 aabba ca
abaa 3 cb ccb ccaca
caa 2 bb abab
bacabb 2 bbabc cc
ac 3 aa caca c
ccbca 3 ac aca a
acb 3 aab cca bcca
bcacccba 1 bbcc
cccabac 1 bbbb
cacbb 1 cbba
bcc 2 cc acbac
aaaab 1 cbc
babbca 2 caaaa bcacc
bcabbaac 2 ac ccc
cbacaca 2 ac babb
bbbaaaa 3 accca b bacb
bc 2 bcacb ba
b 1 ca
baca 3 cb c cacca
caabccb 2 c bca
acbbcb 2 cca a
bcabcc 3 c c ccbbb
aaca 1 cbcab
a 1 ababa
aa 3 aa acbc ac
baababc 2 c cb
ccbabca 2 ab ac
ccbbbba 2 babc b
b 2 c cbcb
bcbb 3 bcacc bb bac
cccb 1 cbcc
acbbcccc 2 baaac bb
ccabbbcb 2 cbc b
a 2 a b
cbcbac 2 c bbcab
ca 2 b acca
b 2 ccab aaccb
bbb 3 cbbba acbab aac
acba 2 cac ab
cbcacacc 3 bc c ac
acab 3 ca bca b
ba 3 bba b ab
ba 2 accbb c
abaa 3 baa baa cb
ccaaa 1 b
ba 2 c ababc
bcaa 2 ca aacbb
bcbbac 1 a
`

type state struct {
	next map[byte]int
	link int
	len  int
	cnt  int
}

// buildSuffixAutomaton constructs the SAM for s and returns its states.
func buildSuffixAutomaton(s string) []state {
	st := make([]state, 2*len(s))
	st[0].next = make(map[byte]int)
	st[0].link = -1
	last, sz := 0, 1
	extend := func(c byte) {
		cur := sz
		sz++
		st[cur].len = st[last].len + 1
		st[cur].cnt = 1
		st[cur].next = make(map[byte]int)
		p := last
		for p != -1 && st[p].next[c] == 0 {
			st[p].next[c] = cur
			p = st[p].link
		}
		if p == -1 {
			st[cur].link = 0
		} else {
			q := st[p].next[c]
			if st[p].len+1 == st[q].len {
				st[cur].link = q
			} else {
				clone := sz
				sz++
				st[clone].len = st[p].len + 1
				st[clone].next = make(map[byte]int)
				for k, v := range st[q].next {
					st[clone].next[k] = v
				}
				st[clone].link = st[q].link
				st[clone].cnt = 0
				for p != -1 && st[p].next[c] == q {
					st[p].next[c] = clone
					p = st[p].link
				}
				st[q].link = clone
				st[cur].link = clone
			}
		}
		last = cur
	}
	for i := 0; i < len(s); i++ {
		extend(s[i])
	}
	st = st[:sz]
	return st
}

// solveCase mirrors the logic from 235C.go: for each query return the count of cyclic rotations.
func solveCase(s string, queries []string) []int64 {
	st := buildSuffixAutomaton(s)

	// accumulate end position counts
	maxLen := 0
	for i := 0; i < len(st); i++ {
		if st[i].len > maxLen {
			maxLen = st[i].len
		}
	}
	cntLen := make([]int, maxLen+1)
	for i := 0; i < len(st); i++ {
		cntLen[st[i].len]++
	}
	for i := 1; i <= maxLen; i++ {
		cntLen[i] += cntLen[i-1]
	}
	order := make([]int, len(st))
	for i := len(st) - 1; i >= 0; i-- {
		l := st[i].len
		cntLen[l]--
		order[cntLen[l]] = i
	}
	for i := len(st) - 1; i > 0; i-- {
		v := order[i]
		if st[v].link >= 0 {
			st[st[v].link].cnt += st[v].cnt
		}
	}

	stamp := make([]int, len(st))
	curStamp := 0

	results := make([]int64, len(queries))
	nS := len(s)
	for qi, x := range queries {
		m := len(x)
		if m > nS {
			results[qi] = 0
			continue
		}
		curStamp++
		P := x + x[:m-1]
		v, l := 0, 0
		var ans int64
		for i := 0; i < len(P); i++ {
			c := P[i]
			for v != -1 && st[v].next[c] == 0 {
				v = st[v].link
				if v >= 0 {
					l = st[v].len
				} else {
					l = 0
				}
			}
			if v == -1 {
				v = 0
				l = 0
			}
			if st[v].next[c] != 0 {
				v = st[v].next[c]
				l++
			}
			if l >= m {
				u := v
				for st[st[u].link].len >= m {
					u = st[u].link
				}
				if stamp[u] != curStamp {
					stamp[u] = curStamp
					ans += int64(st[u].cnt)
				}
			}
		}
		results[qi] = ans
	}
	return results
}

type testCase struct {
	s       string
	queries []string
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Fields(testcasesCData)
	if len(lines) == 0 {
		return nil, fmt.Errorf("no data")
	}
	idx := 0
	toInt := func() (int, error) {
		if idx >= len(lines) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(lines[idx])
		idx++
		return v, err
	}

	t, err := toInt()
	if err != nil {
		return nil, fmt.Errorf("parse T: %w", err)
	}
	cases := make([]testCase, 0, t)
	for ci := 0; ci < t; ci++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("missing case %d", ci+1)
		}
		s := lines[idx]
		idx++
		q, err := toInt()
		if err != nil {
			return nil, fmt.Errorf("case %d q: %w", ci+1, err)
		}
		if idx+q > len(lines) {
			return nil, fmt.Errorf("case %d missing queries", ci+1)
		}
		queries := make([]string, q)
		for i := 0; i < q; i++ {
			queries[i] = lines[idx+i]
		}
		idx += q
		cases = append(cases, testCase{s: s, queries: queries})
	}
	return cases, nil
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
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		var input strings.Builder
		fmt.Fprintf(&input, "%s\n%d\n", tc.s, len(tc.queries))
		for _, q := range tc.queries {
			input.WriteString(q)
			input.WriteByte('\n')
		}
		expect := solveCase(tc.s, tc.queries)
		gotStr, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(gotStr)
		if len(gotFields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d outputs, got %d\n", idx+1, len(expect), len(gotFields))
			os.Exit(1)
		}
		for i, g := range gotFields {
			var v int64
			if _, err := fmt.Sscan(g, &v); err != nil {
				fmt.Fprintf(os.Stderr, "case %d parse output %q: %v\n", idx+1, g, err)
				os.Exit(1)
			}
			if v != expect[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at query %d: expected %d got %d\n", idx+1, i+1, expect[i], v)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

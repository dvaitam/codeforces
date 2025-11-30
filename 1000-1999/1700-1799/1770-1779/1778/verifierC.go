package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input    string
	expected string
}

const testcaseData = `100
7 3
abcbbbb
bcacaba
2 2
bc
cc
3 2
aca
cbb
2 1
bb
cc
4 3
bcba
caac
7 0
cbbacbc
aacaaac
8 0
abcbabcb
cacbcacc
5 3
acbbc
abaaa
1 1
b
a
2 2
aa
aa
7 2
caaccbc
bbbcccb
2 1
ca
bc
6 1
aacbac
ababba
2 0
ca
ac
2 0
ac
ac
2 1
ab
aa
1 0
a
c
2 1
ac
ac
1 1
c
a
5 0
aacbb
baacb
1 0
c
b
4 2
bcbc
acca
1 0
a
b
5 0
cbcaa
bcbcc
5 2
bcbac
cabca
6 0
cbaabb
cbcbcc
3 2
bcb
caa
4 2
aaac
bbcc
7 0
bccbcca
ababcab
8 0
abbbbaba
ccacaabc
7 2
aaacacc
caaacca
5 2
caabb
caabb
2 1
ac
cc
6 0
abaaaa
cbcbbc
1 1
c
c
8 3
bcaabcba
aabbbbca
6 0
abaacb
bbcaba
8 1
abaccabb
bbbaacbb
6 2
abacbb
abbcca
3 3
caa
aac
4 0
baab
ccbb
8 1
bababcab
abaacacb
4 0
bbba
caac
3 0
abb
bca
2 2
ba
cb
7 3
cbbbcca
ccaabcc
6 2
acabca
bcbccb
2 0
ca
aa
3 0
bab
baa
8 2
cbccacac
ccacbabc
7 3
bccaaca
cabccbb
2 1
bb
bb
7 0
acaabcb
aabbabc
2 2
ca
bb
1 1
b
b
1 0
b
a
1 0
c
c
3 2
acc
cba
7 3
aaccbcc
bcacccb
3 3
bcc
aca
1 1
b
b
4 3
bcaa
abba
4 1
accc
cbcb
2 0
bc
ba
2 1
aa
ba
1 1
b
c
7 0
bbcbaba
abcabba
1 0
b
a
3 1
aac
cca
1 1
b
c
1 0
a
a
8 0
bbcbaaab
bbbbccba
2 0
cc
ba
7 1
aabcaaa
bbbcaac
1 1
a
c
2 0
bb
cb
3 3
bcb
aab
6 0
bbabcb
babbba
3 0
bac
cca
8 3
abbaabbc
abbccaba
5 0
cbcba
abacb
3 3
cba
cba
7 2
caabcaa
bbcccbb
2 2
bb
bc
4 0
aabb
cbac
4 1
abbb
cacc
4 0
bccb
cbbb
6 1
aacabb
cbaccb
6 3
cbabbc
cbbcab
5 1
cbaca
abbba
2 2
ac
cc
7 0
baacabc
cbcaacb
4 0
bbcb
babb
7 1
bcacbaa
abbbbbb
4 1
bacc
abaa
2 0
ba
cc
3 1
cbc
acc
5 2
bbacc
ccccc
7 3
bbcbabb
aaaabab
1 0
a
b`

func solve(a, b string, k int) int64 {
	n := len(a)
	letterIndex := make(map[byte]int)
	letters := make([]byte, 0, 10)
	for i := 0; i < n; i++ {
		c := a[i]
		if _, ok := letterIndex[c]; !ok {
			letterIndex[c] = len(letters)
			letters = append(letters, c)
		}
	}
	m := len(letters)
	if k >= m {
		return int64(n) * int64(n+1) / 2
	}
	ans := int64(0)
	maxMask := 1 << m
	for mask := 0; mask < maxMask; mask++ {
		if bits.OnesCount(uint(mask)) > k {
			continue
		}
		cur := int64(0)
		curLen := int64(0)
		for i := 0; i < n; i++ {
			if a[i] == b[i] {
				curLen++
			} else {
				idx := letterIndex[a[i]]
				if (mask>>idx)&1 == 1 {
					curLen++
				} else {
					cur += curLen * (curLen + 1) / 2
					curLen = 0
				}
			}
		}
		cur += curLen * (curLen + 1) / 2
		if cur > ans {
			ans = cur
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
		if pos+3 > len(fields) {
			return nil, fmt.Errorf("case %d: incomplete data", caseNum+1)
		}
		n, err1 := strconv.Atoi(fields[pos])
		k, err2 := strconv.Atoi(fields[pos+1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("case %d: bad n/k", caseNum+1)
		}
		a := fields[pos+2]
		b := fields[pos+3]
		pos += 4

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(k))
		sb.WriteByte('\n')
		sb.WriteString(a)
		sb.WriteByte('\n')
		sb.WriteString(b)
		sb.WriteByte('\n')

		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(a, b, k), 10),
		})
	}
	return cases, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

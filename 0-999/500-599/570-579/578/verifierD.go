package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesD = `2 3
ac
4 3
aaaa
4 3
aacc
3 3
aab
2 2
bb
2 2
bb
3 2
bba
2 2
bb
1 3
a
3 3
cab
4 3
bbaa
3 3
aaa
4 3
cccb
3 2
aab
2 3
ba
3 3
ccb
5 2
baaab
5 2
ababb
1 2
b
1 3
c
3 2
bbb
2 3
cc
1 3
c
2 3
ba
3 3
cab
5 2
babab
3 3
cab
2 3
aa
1 2
a
2 2
ab
5 2
baabb
2 3
aa
4 3
baab
4 2
baba
1 2
b
3 2
aab
3 2
bab
5 2
babba
4 2
abbb
3 3
cac
3 3
bba
3 2
aba
3 2
abb
4 3
aaaa
1 3
a
4 3
acaa
2 3
cb
5 3
aaacc
1 3
c
3 2
aba
5 2
bbaab
4 2
bbab
4 2
aabb
4 3
acac
3 3
bac
2 3
aa
2 3
bb
3 2
aab
1 2
b
5 3
cacab
3 3
aac
1 2
b
5 3
cbcbc
5 3
babcb
3 2
baa
5 2
bbaba
1 3
c
2 3
cb
3 2
baa
5 2
baaab
5 2
aabab
4 2
babb
2 3
cc
5 3
bcbbb
4 3
caba
5 3
cccbb
4 2
aabb
5 3
abbac
4 2
aaab
2 2
ab
5 3
acaaa
2 2
ba
1 3
a
3 2
aba
1 3
c
1 2
b
2 2
bb
5 2
baabb
5 3
bbaba
4 2
baba
1 2
a
5 2
aaaba
4 3
bacc
2 2
bb
3 3
cca
2 3
cc
5 3
bcbaa
3 2
abb
3 2
aab
5 2
aaaba`

type testCase struct {
	n int
	m int
	s string
}

// Embedded solver from 578D.go.
func solve(tc testCase) int {
	n, m := tc.n, tc.m
	s := tc.s

	alphabet := make([]byte, m)
	for i := 0; i < m; i++ {
		alphabet[i] = byte('a' + i)
	}

	unique := make(map[string]struct{})

	for i := 0; i < n; i++ {
		u := s[:i] + s[i+1:]
		for j := 0; j < n; j++ {
			for _, c := range alphabet {
				if j == i && byte(c) == s[i] {
					continue
				}
				t := u[:j] + string(c) + u[j:]
				if t == s {
					continue
				}
				unique[t] = struct{}{}
			}
		}
	}

	return len(unique)
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesD), "\n")
	trimmed := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			trimmed = append(trimmed, strings.TrimSpace(line))
		}
	}
	if len(trimmed)%2 != 0 {
		return nil, fmt.Errorf("testcase lines not paired")
	}
	cases := make([]testCase, 0, len(trimmed)/2)
	for i := 0; i < len(trimmed); i += 2 {
		header := strings.Fields(trimmed[i])
		if len(header) != 2 {
			return nil, fmt.Errorf("line %d: bad header", i+1)
		}
		n, err1 := strconv.Atoi(header[0])
		m, err2 := strconv.Atoi(header[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: parse error", i+1)
		}
		cases = append(cases, testCase{n: n, m: m, s: trimmed[i+1]})
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
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc)
		input := fmt.Sprintf("%d %d\n%s\n", tc.n, tc.m, tc.s)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output\n", idx+1)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

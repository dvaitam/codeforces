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
bacb
aaaabcba
ccbb
aba
c
bcaa
bbb
aabaa
ccaccbca
b
cabccabcb
aaaaa
baaababa
bbbabbabb
cbbc
baababb
cac
acabbbca
cbbbcac
abbccbb
abbccabba
aabbcabc
abacacbca
acbab
baccccbcab
abbacbbaa
babc
acaaab
ca
acaccabbc
accacbc
abcc
cbaa
bbacba
baaa
cbcacbcaa
cababab
ab
abcca
abbcabc
abac
abbab
ccacc
ccca
bcbccabcc
ccbbacc
cabcaa
bbabc
accccbacb
cabbcab
bcbaca
a
cbbcab
acbcbb
abca
cacb
accbcbaac
babb
bb
baabb
bab
cbcabb
acc
ccb
bcbbbbcb
aa
ccccbccc
abbbbbaca
bcc
accba
cbaaab
abbc
cbabb
bcbbac
cb
bbaba
aaaca
cabcbbbbb
babab
baca
bccbb
bc
bbcbababc
cbabab
aaaa
bcccc
cccacbc
bacaca
abccbbcc
bbabcbcbbb
aacc
aab
cabbbcab
abababab
aaac
bcab
bcabbabb
ccacbb
abcaacba
bcbccabc
`

func isPalindrome(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func allSame(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func solveCase(s string) (bool, []string) {
	n := len(s)
	if n == 1 {
		return false, nil
	}
	if !isPalindrome(s) {
		return true, []string{s}
	}
	if allSame(s) {
		return false, nil
	}
	if n%2 == 0 {
		for _, j := range []int{2, 3} {
			if j <= n-2 {
				left, right := s[:j], s[j:]
				if !isPalindrome(left) && !isPalindrome(right) {
					return true, []string{left, right}
				}
			}
		}
		for j := 2; j <= n-2 && j <= 6; j++ {
			left, right := s[:j], s[j:]
			if !isPalindrome(left) && !isPalindrome(right) {
				return true, []string{left, right}
			}
		}
		return false, nil
	}
	if n >= 2 {
		x := s[0]
		y := s[1]
		if x != y {
			alt := true
			for i := 0; i < n; i++ {
				if i%2 == 0 && s[i] != x {
					alt = false
					break
				}
				if i%2 == 1 && s[i] != y {
					alt = false
					break
				}
			}
			if alt {
				return false, nil
			}
		}
	}
	outerChar := s[0]
	outerSame := true
	mid := n / 2
	for i := 0; i < n; i++ {
		if i == mid {
			continue
		}
		if s[i] != outerChar {
			outerSame = false
			break
		}
	}
	if outerSame && s[mid] != outerChar {
		return false, nil
	}
	for _, j := range []int{2, 3} {
		if j <= n-2 {
			left, right := s[:j], s[j:]
			if !isPalindrome(left) && !isPalindrome(right) {
				return true, []string{left, right}
			}
		}
	}
	for j := 2; j <= n-2 && j <= 6; j++ {
		left, right := s[:j], s[j:]
		if !isPalindrome(left) && !isPalindrome(right) {
			return true, []string{left, right}
		}
	}
	return false, nil
}

type testCase struct {
	s string
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	tests := make([]testCase, 0)
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tests = append(tests, testCase{s: line})
		_ = idx
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("no tests")
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(tc.s)
	sb.WriteByte('\n')
	return sb.String()
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		ok, parts := solveCase(tc.s)
		var expected strings.Builder
		if !ok {
			expected.WriteString("NO")
		} else {
			expected.WriteString("YES\n")
			expected.WriteString(strconv.Itoa(len(parts)))
			expected.WriteByte('\n')
			for i, p := range parts {
				if i > 0 {
					expected.WriteByte(' ')
				}
				expected.WriteString(p)
			}
		}

		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected.String() {
			fmt.Printf("case %d failed: expected %q got %q\n", idx+1, expected.String(), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

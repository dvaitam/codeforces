package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// testcasesRaw embeds the original testcasesB.txt content.
const testcasesRaw = `acbccbcbcbbbccac
abccbbac
cabbccbaacaacab
bb
baacbbcbccacacc
cba
bbbccacabbbcc
ccbbabbbbb
ac
bcbaabba
caccabbacbbbaabaaacc
accbabbaaccbcacc
cacccbabbcc
caab
bbcaaab
cccaaaccccbc
abcb
baab
b
abbcbabaabcabb
baaaaacccacba
cacaabcac
cbacba
bbbbccbacaac
baaabc
bbcaacabbb
abcbcaca
aabcbccabbab
babbbaaabac
cbbab
bbcabbcc
bacbaabacba
abcbacbacbbcaabcaabb
cbbacbbbcacaaabbcba
aabbbcac
cabccbc
bbbccaaac
cbacbbcbcbabbccbac
cabbacbacaabbbac
ccc
aacabccbcaac
bcbbabbcb
abcacbaaabbbbb
abaccaba
aacccaa
cacbbbbaccccac
cbcbbbcbabbaaaabab
abcac
babcacab
bbacaabaacbaabc
bcbbbbcaab
abb
bb
aaaa
babca
baaabaabbbbacbccac
ccbcbbaaaccaacccbc
bccaccc
ab
ccba
cbacaacbcbabcbbcaa
cbbccaabaacaaabcca
abccacbab
cacc
baaccacbcbcbcab
cabbaaabcccabb
bccccabbcab
caccaabcaacacca
abaabcbaabba
acbbabbccacb
acabbaacaccccca
aaaa
cccaca
accbcaccabbccbacabab
babacbbbbbba
bcbaccbcabbcabbcbac
ab
bc
caccbbbaac
cccccbaa
c
abbbcbbbbcbaaacabab
ccbccaacbbacaacc
cccbbbbbabcc
cccbacbaccc
acccbbbacbccaabca
ccbaccacabab
accbcaaababcaccb
bbbcbabab
abbcccbcaabb
caccbbabbcaa
baacacbaa
ba
cbbca
bcaaaabccccac
cc
abcaacababbcc
cbacaaccaccba
aacbcabcacb
bcaca`

type testCase struct {
	input string
}

// referenceSolution embeds the algorithm from 126B.go (KMP border search).
func referenceSolution(s string) string {
	n := len(s)
	if n == 0 {
		return "Just a legend"
	}
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	best := 0
	for i := 0; i < n-1; i++ {
		if pi[i] > best {
			best = pi[i]
		}
	}
	k := pi[n-1]
	for k > 0 {
		if best >= k {
			return s[:k]
		}
		k = pi[k-1]
	}
	return "Just a legend"
}

func parseTestcases() []testCase {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	cases := make([]testCase, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		cases = append(cases, testCase{input: line})
	}
	return cases
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := parseTestcases()
	for idx, tc := range tests {
		input := tc.input + "\n"
		exp := referenceSolution(tc.input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

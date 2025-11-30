package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	a string
	b string
	c string
}

// Embedded testcases (from testcasesB.txt) so the verifier needs no external files.
const rawTestcases = `bcbb bb b
abaac b ccc
ca cb b
bbcca c bb
caacb cc c
b acb ca
aa cba a
abc bca cb
ccbba c bbc
aa ac c
aca aa a
bccbca acc bcb
ccba bc abc
aacbac ab a
aac aa c
caaaca cca bab
a a aca
acac b cab
c b b
cbac a c
bcbc a cc
aa b cba
aabcb cc bcb
acca bca bc
a abb cb
ccacbb cb caa
baaac b bcc
bccbc ca a
cabc b cc
b b ba
caca a bcb
aca c c
accab b c
bc a ab
accc c ba
aa ac b
caccc cb cc
bcaa bcb aa
bb ca bc
a a cb
aba bc aab
ab bbb baa
bbaba cb ba
caa cc bca
c a a
bccb b b
ababa bca bab
a c bca
b ac aa
abbbc a a
cbccb cb c
cca cc aa
bba cab cab
cbaac aa aaa
baa c bb
cacac cc acb
cb bb cca
c abc c
bbb bb b
a a bcb
b a bc
cabbac b bba
a ac c
ac ccb ab
ccbc c b
ccbabb c cac
bb a bb
abbacc c c
cc c c
aabc baa ba
c ab b
bbcbab aa b
baaaba c aa
c c aca
aca aa bab
aaa bbb bb
acaac cba ba
ca a ab
aac ab aca
cb ca cb
aab cbc ca
bcbb ab b
abac c c
ba bba ab
bbc cab a
bcb a aba
cabcba cb acb
abc aab b
bbacb bbc aca
b c ba
abbbc a c
cabcc bcb b
caaa ca bbc
cbbc b cba
cbbc ab bac
aabb b aca
ccbaba a cab
caacba abb cb
bccb a bc
aa abb bb
aba c abc
aba bc cc`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(rawTestcases, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 strings, got %d", idx+1, len(parts))
		}
		cases = append(cases, testCase{a: parts[0], b: parts[1], c: parts[2]})
	}
	return cases, nil
}

func solve(tc testCase) string {
	var freqA, freqB, freqC [26]int
	for i := 0; i < len(tc.a); i++ {
		freqA[tc.a[i]-'a']++
	}
	for i := 0; i < len(tc.b); i++ {
		freqB[tc.b[i]-'a']++
	}
	for i := 0; i < len(tc.c); i++ {
		freqC[tc.c[i]-'a']++
	}
	maxB := len(tc.a) / len(tc.b)
	for i := 0; i < 26; i++ {
		if freqB[i] > 0 {
			if v := freqA[i] / freqB[i]; v < maxB {
				maxB = v
			}
		}
	}
	bestX, bestY, bestSum := 0, 0, 0
	for x := 0; x <= maxB; x++ {
		var rem [26]int
		ok := true
		for i := 0; i < 26; i++ {
			rem[i] = freqA[i] - x*freqB[i]
			if rem[i] < 0 {
				ok = false
				break
			}
		}
		if !ok {
			break
		}
		y := len(tc.a) / len(tc.c)
		for i := 0; i < 26; i++ {
			if freqC[i] > 0 {
				if v := rem[i] / freqC[i]; v < y {
					y = v
				}
			}
		}
		if x+y > bestSum {
			bestSum = x + y
			bestX = x
			bestY = y
		}
	}
	var res []byte
	bBytes := []byte(tc.b)
	cBytes := []byte(tc.c)
	for i := 0; i < bestX; i++ {
		res = append(res, bBytes...)
	}
	for i := 0; i < bestY; i++ {
		res = append(res, cBytes...)
	}
	for i := 0; i < 26; i++ {
		used := bestX*freqB[i] + bestY*freqC[i]
		left := freqA[i] - used
		for j := 0; j < left; j++ {
			res = append(res, byte('a'+i))
		}
	}
	return string(res)
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%s\n%s\n%s\n", tc.a, tc.b, tc.c)
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := solve(tc)
		input := buildInput(tc)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}

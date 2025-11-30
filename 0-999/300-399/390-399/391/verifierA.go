package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt.
const testcasesAData = `
C 0
C 0
A 0
G 0
C 0
C 0
G 0
C 0
G 0
T 0
T 0
G 0
T 0
A 0
G 0
TG 0
AA 1
GC 0
AG 0
CG 0
TC 0
CG 0
AA 1
AC 0
AC 0
GT 0
GA 0
TT 1
TT 1
CA 0
AGC 0
AGG 1
AGT 0
GCA 0
CGT 0
GTT 1
TAG 0
CAA 1
TTA 1
ACG 0
TTC 1
GCC 1
GAG 0
ACG 0
TTA 1
GATG 0
TGCA 0
ATTA 1
AAAT 0
ACAG 0
AAAT 0
TACT 0
AACA 1
GATA 0
GGCT 1
ACAA 1
CTGG 1
CTTA 1
TTGG 2
ACTA 0
CCGGC 2
GTACA 0
GAGTT 1
CGGGT 0
GCCAA 2
TGTTT 0
CCCAC 0
CATCA 0
GTCCA 1
ACGGC 1
ACTAT 0
ACCGA 1
TAAAT 0
ATGGT 1
ACCAA 2
GCAGTG 0
ATGAAA 0
TGGGAC 0
CCGTTC 2
GATTGG 2
GGAGAA 2
GTTGGC 2
TGACTA 0
GTAGCG 0
GCAACC 2
GGACAC 1
CAGGTT 2
CAAATT 1
ACAACG 1
CCTCAG 1
TGTCTGA 0
AACTACC 2
GTATTAT 1
CCGTACT 1
CCCGCCT 1
TAGGGAT 0
GTCGCAA 1
AATTAGA 2
CGTTCGC 1
AAACTGC 0
`

type testCase struct {
	s string
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesAData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields, got %d", idx+1, len(parts))
		}
		cases = append(cases, testCase{s: parts[0]})
	}
	return cases, nil
}

// solve mirrors 391A.go: count even-length runs.
func solve(s string) int {
	n := len(s)
	res := 0
	for i := 0; i < n; {
		j := i + 1
		for j < n && s[j] == s[i] {
			j++
		}
		if (j-i)%2 == 0 {
			res++
		}
		i = j
	}
	return res
}

func runCandidate(bin string, s string) (int, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(s + "\n")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	valStr := strings.TrimSpace(out.String())
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("parse output %q: %v", valStr, err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc.s)
		got, err := runCandidate(bin, tc.s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

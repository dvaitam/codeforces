package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n  int
	s1 string
	s2 string
}

// solve embeds the logic from 1730D.go.
func solve(tc testCase) string {
	counts := make(map[[2]byte]int)
	b1 := []byte(tc.s1)
	b2 := []byte(tc.s2)
	for i := 0; i < tc.n; i++ {
		a := b1[i]
		b := b2[tc.n-1-i]
		if a > b {
			a, b = b, a
		}
		counts[[2]byte{a, b}]++
	}

	oddSame := 0
	ok := true
	for pair, c := range counts {
		if pair[0] != pair[1] {
			if c%2 == 1 {
				ok = false
				break
			}
		} else {
			if c%2 == 1 {
				oddSame++
			}
		}
	}
	if ok {
		if tc.n%2 == 0 {
			if oddSame > 0 {
				ok = false
			}
		} else {
			if oddSame != 1 {
				ok = false
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

// Embedded copy of testcasesD.txt.
const testcaseData = `
6 caccec abebce
1 d d
3 ebc ced
4 adcb eabd
4 abac abbd
1 c e
4 aebd ccdd
4 abec bbbb
3 bbe cea
2 bd bb
5 bebbe cadba
4 bded bdee
4 abae aaab
6 ddcdcd cabebb
1 a e
1 d b
4 beec aebc
6 cccecb adbedc
4 cdbd ccce
4 dede cebb
4 bceb ebee
5 edccd cdeab
3 daa bea
3 abb ebc
5 caaae bcbdd
6 ccaeda abaeec
3 bdb cbb
4 beee deba
1 d e
1 b e
1 d c
3 eba cbc
6 cddecc adbaec
5 deddb edbae
1 b d
3 adc cac
5 eaaae abadd
4 acba deba
2 da ae
4 becd dead
5 bebbe baedc
3 cbd aae
1 a c
2 be ba
3 acb cea
1 c b
3 dcd dee
1 a c
4 dded edea
2 bd ec
2 ab ba
4 eaae cecc
6 ddeced ecbead
3 cea bea
4 acde beae
5 aaeca cbccd
4 dabd cbad
5 bdced badee
6 bedcec debedd
2 cd cb
2 cb ee
3 dbe ceb
3 bbc ebe
6 eaddbb ecebae
4 bcdc ebcd
3 aab beb
5 dacec ddcda
1 a c
2 ce ed
3 edd cba
6 aaeaee ceaeba
1 a c
1 e d
6 beeebc cbbebb
2 bc dc
2 ec ba
2 db be
1 c e
6 dbcaec ebabbe
6 adbbde aecbbd
3 dbc bac
4 ccee cbce
4 abda eaeb
6 dbdceb bacdcd
5 acdde beeba
3 ede edb
4 bdaa edae
6 dbabbb aeddae
2 eb bd
1 b a
5 bebee eeded
6 cdeacc babdbe
6 babbdd aeccdc
6 caeadc eeabbd
2 ad cb
5 cbaec cacbb
3 edd dab
4 acbe bedc
2 eb ed
3 bab bde
`

var expectedOutputs = []string{
	"NO",
	"YES",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"YES",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
	"NO",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields, got %d", i+1, len(parts))
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		tests = append(tests, testCase{n: n, s1: parts[1], s2: parts[2]})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	input := fmt.Sprintf("1\n%d\n%s\n%s\n", tc.n, tc.s1, tc.s2)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := runCase(bin, tc, expectedOutputs[i]); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

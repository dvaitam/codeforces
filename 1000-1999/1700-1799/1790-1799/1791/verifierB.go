package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
3 LUL
8 DDDRLDLD
7 LDURLUL
1 L
9 LDRDLRDDR
6 RRDULD
9 LRULUDRUU
10 DDLDRDDRUU
2 DL
3 DUD
1 D
1 U
10 DRRRLRRDUU
8 ULDRRDLD
6 RDDUDU
1 U
8 LRRRLULL
2 LD
1 U
3 RDD
9 RDRDRLLDR
1 L
4 DLDL
6 LDUDLU
6 LDLLLU
9 RDDDRLURL
5 URUDD
1 R
10 ULULDLULDD
10 URRLUDUDDL
4 LRUR
8 DUUURLLU
5 RLRRR
5 LLLRL
3 ULL
8 DDDDLLRD
2 LL
7 RLRRRDD
9 UDDRLUULD
6 UURLUD
2 RR
5 RDLLL
8 LULRRRUD
8 DLURDUDU
4 RDLR
3 LLD
7 RLUDRDU
6 LRRRUD
7 UUDDRUD
3 UUU
10 UULLLDDDUU
8 DDDLRUUU
7 UURRRLR
2 UL
2 LD
5 LDDDL
8 ULUDRRUU
1 L
3 DDU
7 RDUDLRR
9 ULDLULLUR
10 DUUDDLDDUD
9 LDDLRLLRU
10 RDDLRRRRLL
3 DDR
9 LDDULULRD
4 UDUU
9 UURDRLDDR
6 UURDDL
3 UUD
5 LURLD
6 LLLLRR
8 URLURDDD
8 RDLDRLRL
4 RLRL
1 R
2 DL
5 RLLRD
5 UUDLL
6 DLLRRU
10 LDDDRUDLLL
5 DRLDL
1 L
3 RDD
2 UL
7 ULULRLL
6 RRUUUD
5 DLDLL
9 RUDLDLUDR
7 DLDDRRR
10 UDLDRURDDR
7 LDDRUDR
9 RLDUDDLUU
10 URRDDLDLDL
1 L
8 RRDDULUR
5 LUDLD
3 RDL
1 D
3 RUD`

type testCase struct {
	input    string
	expected string
}

func solve(n int, s string) string {
	x, y := 0, 0
	for _, ch := range s {
		switch ch {
		case 'L':
			x--
		case 'R':
			x++
		case 'U':
			y++
		case 'D':
			y--
		}
		if x == 1 && y == 1 {
			return "YES"
		}
	}
	return "NO"
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
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: incomplete data", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		s := fields[pos+1]
		pos += 2

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		sb.WriteString(s)
		sb.WriteByte('\n')

		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solve(n, s),
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
		fmt.Println("usage: verifierB /path/to/binary")
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

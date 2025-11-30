package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `TA
5
3
2 2
1 1
1 4

C
2
2
1 1
1 2

T
3
2
1 1
2 2

TATGCC
2
1
4 5

TATTC
5
3
4 5
4 7
3 6

AAA
1
3
3 3
4 4
2 3

CGTC
3
2
4 4
5 5

A
1
1
1 1

CTG
1
0

ATT
2
0

TCG
2
0

CGTCT
8
3
5 5
3 3
3 3

TT
2
2
1 2
2 3

A
1
0

CCAGCT
1
0

TGCGC
4
1
4 5

C
1
3
1 1
2 2
2 2

T
3
2
1 1
1 2

TA
7
2
1 2
1 4

G
2
1
1 1

GTGCTT
7
3
4 4
7 7
4 4

CTACG
1
3
4 5
2 2
4 5

TAGGC
2
0

TG
4
3
1 1
1 3
5 5

TCATA
3
0

AAG
2
3
1 2
4 4
4 5

A
3
2
1 1
2 2

TG
2
3
1 1
2 2
2 4

GTT
3
0

AT
1
0

GACG
6
2
3 3
2 2

T
5
3
1 1
2 2
1 2

TCG
3
0

CT
4
2
1 2
1 1

CTTTCT
5
0

A
1
2
1 1
1 2

C
3
3
1 1
2 2
2 2

CTG
1
3
3 3
3 4
4 5

CGT
11
3
2 2
2 4
1 4

TAT
4
3
1 3
4 6
2 4

T
2
1
1 1

TACCA
2
1
1 3

T
1
0

CCGCA
5
2
1 5
3 9

TG
2
0

CT
4
3
2 2
2 2
4 4

T
1
0

AC
1
0

CA
1
3
1 2
4 4
2 3

TTGA
4
1
3 4

TACTAA
3
2
5 6
6 6

AGCACA
2
2
5 5
6 7

AACTA
1
0

CCTC
1
2
3 4
6 6

GG
2
0

TCGTC
7
1
1 4

TGG
5
2
3 3
3 4

A
5
3
1 1
1 2
1 2

C
1
2
1 1
1 2

CCAC
1
0

TCA
2
1
3 3

CTATTT
1
1
2 3

C
1
0

GACCGG
1
3
5 6
6 7
7 7

CT
2
1
2 2

CTGTA
3
1
2 4

CAAAC
5
2
3 5
8 8

GCCGGT
1
2
1 4
3 10

GCATTA
4
0

AC
3
1
2 2

AGGCCA
5
1
2 6

GT
5
2
1 2
3 4

ACT
1
1
2 3

C
1
1
1 1

TAACC
5
2
5 5
3 6

GT
1
0

GAA
4
2
3 3
1 2

AGATCC
1
1
2 2

A
2
2
1 1
1 2

CAAACG
6
1
1 5

TTCAA
2
0

T
5
3
1 1
1 2
4 4

CGAC
3
2
4 4
2 4

TTC
3
3
3 3
3 4
6 6

C
2
2
1 1
2 2

A
1
0

CCCTGG
6
1
6 6

TATA
2
0

G
1
0

GA
3
1
2 2

CAT
4
2
1 2
5 5

GAG
4
1
1 2

CAATA
8
1
2 5

TGTGTG
3
0

AGGT
3
1
2 3

T
1
0

CGCTG
6
3
4 4
4 4
7 7

TCCT
3
2
2 3
5 6

GGTA
10
3
3 3
5 5
1 4

CAA
3
1
2 2`

type operation struct {
	l int
	r int
}

type testCase struct {
	s   string
	k   int
	ops []operation
}

// solveLogic mirrors 217E.go: maps each query position back through the operations.
func solveLogic(tc testCase) string {
	s := tc.s
	k := tc.k
	n := len(tc.ops)
	L := make([]int, n)
	R := make([]int, n)
	segLen := make([]int, n)
	half := make([]int, n)
	for i, op := range tc.ops {
		L[i] = op.l
		R[i] = op.r
		ln := op.r - op.l + 1
		segLen[i] = ln
		half[i] = ln / 2
	}
	out := make([]byte, k)
	for pos := 1; pos <= k; pos++ {
		p := pos
		for i := n - 1; i >= 0; i-- {
			li := L[i]
			ri := R[i]
			ln := segLen[i]
			if p > ri && p <= ri+ln {
				idx := p - ri
				h := half[i]
				if idx <= h {
					p = li + idx*2 - 1
				} else {
					p = li + (idx-h-1)*2
				}
			} else if p > ri+ln {
				p -= ln
			}
		}
		out[pos-1] = s[p-1]
	}
	return string(out)
}

func parseCase(lines []string, start int) (testCase, int, error) {
	if start >= len(lines) {
		return testCase{}, start, fmt.Errorf("no data")
	}
	s := strings.TrimSpace(lines[start])
	start++
	if start >= len(lines) {
		return testCase{}, start, fmt.Errorf("missing k")
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[start]))
	if err != nil {
		return testCase{}, start, fmt.Errorf("parse k: %w", err)
	}
	start++
	if start >= len(lines) {
		return testCase{}, start, fmt.Errorf("missing n")
	}
	n, err := strconv.Atoi(strings.TrimSpace(lines[start]))
	if err != nil {
		return testCase{}, start, fmt.Errorf("parse n: %w", err)
	}
	start++
	ops := make([]operation, n)
	for i := 0; i < n; i++ {
		if start >= len(lines) {
			return testCase{}, start, fmt.Errorf("missing op %d", i+1)
		}
		fields := strings.Fields(lines[start])
		if len(fields) != 2 {
			return testCase{}, start, fmt.Errorf("op %d: wrong field count", i+1)
		}
		l, err := strconv.Atoi(fields[0])
		if err != nil {
			return testCase{}, start, fmt.Errorf("op %d: parse l: %w", i+1, err)
		}
		r, err := strconv.Atoi(fields[1])
		if err != nil {
			return testCase{}, start, fmt.Errorf("op %d: parse r: %w", i+1, err)
		}
		ops[i] = operation{l: l, r: r}
		start++
	}
	return testCase{s: s, k: k, ops: ops}, start, nil
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for i := 0; i < len(lines); {
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			break
		}
		tc, next, err := parseCase(lines, i)
		if err != nil {
			return nil, fmt.Errorf("case starting at line %d: %w", i+1, err)
		}
		cases = append(cases, tc)
		i = next
	}
	return cases, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(tc.s)
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", tc.k))
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.ops)))
	for _, op := range tc.ops {
		fmt.Fprintf(&sb, "%d %d\n", op.l, op.r)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expect := solveLogic(tc)
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

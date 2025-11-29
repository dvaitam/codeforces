package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	a int64
	b int64
}

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `100
700 401
858 723
538 283
535 832
242 870
221 917
696 604
846 973
430 594
282 462
505 677
657 718
939 813
366 85
333 628
119 499
602 646
344 866
195 249
17 750
278 120
723 226
381 814
175 341
437 836
64 104
802 150
876 715
225 47
837 588
650 932
959 548
617 697
76 28
128 651
194 621
851 590
123 401
94 380
854 119
38 621
23 200
985 994
190 736
127 491
216 745
820 63
960 696
24 558
436 636
104 856
267 72
227 74
663 309
359 447
185 63
516 479
41 611
104 717
401 205
267 368
927 750
482 859
924 941
584 174
715 689
209 990
786 60
808 693
163 866
166 351
543 257
121 612
944 453
682 180
14 483
698 420
922 583
896 521
940 319
665 366
398 858
674 257
158 575
708 13
469 760
81 344
757 47
558 288
139 246
781 977
494 361
625 295
690 368
605 970
914 649
875 636
136 733
318 398
767 425`

func parseTestcases() ([]testcase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %v", err)
	}
	if len(fields) != 1+t*2 {
		return nil, fmt.Errorf("malformed testcases: expected %d numbers, got %d", 1+t*2, len(fields))
	}
	cases := make([]testcase, 0, t)
	idx := 1
	for i := 0; i < t; i++ {
		aVal, err := strconv.ParseInt(fields[idx], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse a at case %d: %v", i+1, err)
		}
		bVal, err := strconv.ParseInt(fields[idx+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse b at case %d: %v", i+1, err)
		}
		idx += 2
		cases = append(cases, testcase{a: aVal, b: bVal})
	}
	return cases, nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// solve implements the logic from 1414B.go for a single testcase.
func solve(tc testcase) string {
	return strconv.FormatInt(gcd(tc.a, tc.b), 10)
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&input, "%d %d\n", tc.a, tc.b)
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var expected strings.Builder
	for i, tc := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(solve(tc))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

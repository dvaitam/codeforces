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
100
813382118
-172692001
627694678
911784257
-96829398
-913060454
-443980515
97954097
43521778
-130410564
971893209
683194652
782095536
-348640892
23484162
-231094825
252803391
914826684
950157577
-530897812
83806779
-700911987
-394757831
-699898217
623077180
-796352494
327937310
716703952
-462041733
953665205
143671689
514345872
739928271
292574851
937386629
-684402796
-333963156
-787907338
567301757
-841639335
930240526
826378634
468844309
-290906865
13918763
202190735
-783745798
-240238909
-67623087
-320972758
311869788
375298783
960676318
-560887388
186535554
24370691
-49323255
858238927
119598413
-440597023
-866255615
728784093
972388339
178322771
967083172
-969845676
-799700194
545554038
804082152
-143532968
525257611
771341094
685877224
434848064
342748146
-997545822
314038990
59950399
778252349
863162773
-284597743
-476205387
568261309
-301628955
511060853
869322740
-864742298
-589686554
969283239
218720039
-523894505
-487576197
725170358
-693995623
724814782
166062048
-37992670
-804115231
-827243928
-312687983
`

type testCase struct {
	n int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesAData), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("no data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("parse T: %w", err)
	}
	if len(lines)-1 < t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(lines)-1)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		val, err := strconv.Atoi(strings.TrimSpace(lines[i+1]))
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", i+2, err)
		}
		cases = append(cases, testCase{n: val})
	}
	return cases, nil
}

// solve mirrors 313A.go logic.
func solve(n int) int {
	if n >= 0 {
		return n
	}
	a := n / 10
	b := (n/100)*10 + n%10
	if a > b {
		return a
	}
	return b
}

func runCandidate(bin string, n int) (int, error) {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc.n)
		got, err := runCandidate(bin, tc.n)
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

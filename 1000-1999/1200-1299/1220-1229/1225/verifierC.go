package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `53125 -965
80189 -208
822297 -260
944000 -105
196801 875
879648 -670
5155 674
446578 -936
257839 -873
948959 -196
118043 -244
392477 131
884576 678
630049 -727
896137 -221
548781 308
98553 525
333821 3
662799 -993
490404 -519
423300 150
83869 -114
889474 231
618480 -230
19984 -296
539653 -400
243036 -185
520561 966
423845 -244
52251 724
64709 -166
635500 -123
615861 889
164833 414
152120 -323
199910 -712
515703 871
433403 -358
790627 580
6247 812
383267 -237
380156 494
83989 604
788396 -215
138176 464
676526 943
38003 789
148376 -899
768486 -368
277356 -319
42531 -838
133847 -719
251586 -350
221043 544
165098 -91
208780 -164
390559 -731
326282 606
622429 -625
482490 -950
602002 -635
769259 938
657407 -751
33172 529
147728 -152
691794 -403
537882 545
222875 -884
941965 -855
632032 973
549890 588
538565 106
959153 910
744183 -132
767989 -451
726048 -38
586122 -97
488184 324
607750 268
948456 -834
592442 -298
111489 785
314615 -531
817760 -909
249107 -763
215571 -634
688026 540
412002 542
859961 822
424866 296
135052 693
592146 -527
589446 315
216697 -1
917407 413
616193 855
673320 -686
719154 978
777144 286
794152 -762`

// Embedded reference logic from 1225C.go.
func solve(n, p int64) int64 {
	for k := int64(1); k <= 60; k++ {
		s := n - p*k
		if s <= 0 {
			continue
		}
		if s < k {
			continue
		}
		if int64(bits.OnesCount64(uint64(s))) <= k {
			return k
		}
	}
	return -1
}

type testCase struct {
	n int64
	p int64
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields", idx+1)
		}
		n, err1 := strconv.ParseInt(fields[0], 10, 64)
		p, err2 := strconv.ParseInt(fields[1], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: invalid numbers", idx+1)
		}
		res = append(res, testCase{n: n, p: p})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.p)
		expected := fmt.Sprintf("%d", solve(tc.n, tc.p))

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", i+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod int64 = 1000000007

type testCase struct {
	n int64
	k int64
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
138 582
868 821
783 64
262 120
508 779
461 483
668 388
808 214
97 499
30 914
856 399
444 622
781 785
3 712
457 272
739 821
235 605
968 104
924 325
32 22
27 665
555 9
962 902
391 702
222 992
433 743
30 540
228 782
449 961
508 566
239 353
237 693
225 779
471 975
297 948
23 426
858 938
570 944
658 102
191 644
742 880
304 123
761 340
918 738
997 728
513 958
991 432
520 849
933 686
195 310
291 601
997 903
512 866
964 517
403 603
874 35
492 248
762 816
414 424
681 177
376 561
904 719
795 690
756 383
89 449
680 520
111 797
168 533
861 402
380 501
751 30
481 44
316 720
869 629
608 592
404 662
175 172
515 232
13 789
205 552
943 880
562 237
415 526
353 975
868 591
362 470
932 275
676 561
624 980
747 5
393 802
878 840
978 907
961 758
525 828
133 531
797 574
211 436
973 57
493 890`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("case %d invalid format", i+1)
		}
		n, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		k, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", i+1, err)
		}
		res = append(res, testCase{n: n, k: k})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

func powMod(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

// solve mirrors 1514B.go.
func solve(tc testCase) string {
	return fmt.Sprintf("%d", powMod(tc.n, tc.k))
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
100
980 884
972 870
59 94
88 370
857 174
755 829
687 875
317 258
622 218
623 37
597 698
164 442
655 403
824 741
882 522
974 381
559 959
457 515
276 923
38 892
30 373
478 955
328 930
391 434
915 906
540 169
575 182
243 237
26 181
334 178
141 523
524 369
528 691
575 187
917 457
817 425
754 538
930 931
783 373
810 608
364 371
881 985
458 166
979 773
411 733
758 473
672 544
257 502
287 948
512 513
529 852
817 363
679 905
467 922
926 473
361 582
745 943
572 742
469 499
676 228
965 333
836 717
857 171
899 930
633 275
793 934
493 317
312 981
820 724
853 517
577 531
521 668
632 603
418 320
750 213
502 525
377 957
702 639
905 78
805 841
351 744
10 930
836 196
764 109
62 589
670 51
281 606
234 699
898 938
110 773
536 140
876 273
252 845
217 967
903 62
435 920
736 778
34 59
373 369
178 256
`

type testCase struct {
	n int64
	k int64
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
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

func expectedValue(n, k int64) int64 {
	div := n - 1
	q := k / div
	r := k % div
	if r == 0 {
		return q*n - 1
	}
	return q*n + r
}

func parseTestcases() ([]testCase, error) {
	rawLines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, ln := range rawLines {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			lines = append(lines, ln)
		}
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	start := 0
	count := 0
	if fields := strings.Fields(lines[0]); len(fields) == 1 {
		if v, err := strconv.Atoi(fields[0]); err == nil {
			count = v
			start = 1
		}
	}
	if count == 0 {
		count = len(lines)
		start = 0
	}
	if start+count != len(lines) {
		return nil, fmt.Errorf("testcase count mismatch: declared %d actual %d", count, len(lines)-start)
	}
	cases := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		parts := strings.Fields(lines[start+i])
		if len(parts) != 2 {
			return nil, fmt.Errorf("case %d expected 2 numbers got %d", i+1, len(parts))
		}
		nVal, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		kVal, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{n: nVal, k: kVal})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
		want := expectedValue(tc.n, tc.k)
		gotStr, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Printf("case %d failed: invalid output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

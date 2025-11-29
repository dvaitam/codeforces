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
	n int64
	k int64
}

const testcasesRaw = `100
866 395
778 912
432 42
267 989
525 498
416 941
804 850
312 992
490 367
599 914
931 224
518 143
290 144
775 98
635 819
258 932
547 723
831 617
925 151
319 102
749 76
922 871
702 339
485 574
105 363
446 324
627 656
936 210
991 566
490 454
888 534
268 64
826 941
563 938
16 96
738 861
410 728
846 804
686 641
3 627
507 848
890 342
251 748
335 721
893 65
197 940
583 228
246 823
992 146
824 557
460 94
84 328
898 521
957 502
113 309
566 299
725 128
562 341
836 945
555 209
988 819
619 561
603 295
457 94
612 818
396 325
591 248
299 189
195 842
193 34
629 673
268 488
72 92
697 776
135 898
155 946
41 863
84 920
718 946
851 554
701 401
859 723
539 283
536 832
243 870
222 917
697 604
847 973
431 594
283 462
506 677
658 718
940 813
367 85
334 628
120 499
603 646
345 866
196 249
18 750`

func parseTestcases(raw string) []testCase {
	fields := strings.Fields(raw)
	if len(fields) < 1 {
		panic("no testcase data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(fmt.Sprintf("bad testcase count: %v", err))
	}
	fields = fields[1:]
	if len(fields)%2 != 0 || len(fields)/2 != t {
		panic("testcase pair count mismatch")
	}
	res := make([]testCase, 0, t)
	for i := 0; i < len(fields); i += 2 {
		n, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("bad n at pair %d: %v", i/2+1, err))
		}
		k, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("bad k at pair %d: %v", i/2+1, err))
		}
		res = append(res, testCase{n: n, k: k})
	}
	return res
}

func smallestDivisor(n int64) int64 {
	if n%2 == 0 {
		return 2
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return i
		}
	}
	return n
}

// expected replicates the logic from 1350A.go.
func expected(tc testCase) int64 {
	n := tc.n
	k := tc.k
	if n%2 == 0 {
		n += 2 * k
	} else {
		d := smallestDivisor(n)
		n += d
		k--
		n += 2 * k
	}
	return n
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases(testcasesRaw)

	var input strings.Builder
	input.WriteString(strconv.Itoa(len(tests)))
	input.WriteByte('\n')
	for _, tc := range tests {
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	}
	gotStr, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	gotLines := strings.Fields(gotStr)
	if len(gotLines) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d lines got %d\n", len(tests), len(gotLines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := strconv.FormatInt(expected(tc), 10)
		if strings.TrimSpace(gotLines[i]) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

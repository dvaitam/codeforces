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
	a  int
	b  int
	m  int
	r0 int
}

// solve embeds the logic from 172B.go.
func solve(tc testCase) string {
	visited := make([]int, tc.m)
	for i := range visited {
		visited[i] = -1
	}
	cur := tc.r0
	for i := 1; ; i++ {
		cur = (tc.a*cur + tc.b) % tc.m
		if visited[cur] != -1 {
			return strconv.Itoa(i - visited[cur])
		}
		visited[cur] = i
	}
}

// Embedded copy of testcasesB.txt.
const testcaseData = `
958 91 271 22
323 641 159 149
496 64 578 539
445 327 491 197
893 24 747 119
375 325 498 125
534 312 832 352
414 680 170 618
92 863 932 507
548 811 613 577
118 305 414 111
130 941 278 814
714 309 783 374
598 667 1000 707
453 851 310 57
520 443 83 40
18 660 190 110
120 436 672 354
33 157 336 133
629 654 963 55
431 210 822 148
473 125 232 225
38 851 144 74
159 758 928 397
381 19 837 112
156 640 564 434
666 240 840 192
560 860 868 284
82 914 812 44
630 745 105 486
26 861 226 694
7 684 442 295
856 677 308 96
361 189 329 279
892 330 473 418
688 381 976 292
925 207 467 110
332 547 414 91
84 327 422 219
945 438 626 304
343 881 804 284
918 514 427 780
440 12 565 404
966 995 142 257
655 612 867 311
443 460 408 33
731 923 505 731
275 756 854 453
558 315 856 495
120 712 742 214
5 699 877 269
309 564 323 885
777 986 460 732
138 423 184 18
160 812 388 150
786 428 682 494
102 457 946 538
853 526 501 258
305 963 314 225
532 518 590 374
652 966 440 860
268 846 852 488
278 146 334 150
46 201 771 397
361 311 988 67
639 780 267 269
656 264 748 116
182 197 811 59
351 902 659 231
332 986 542 699
767 442 853 597
194 835 503 506
532 674 922 567
754 787 31 133
820 175 38 400
75 132 778 495
638 868 953 705
207 947 477 451
626 205 829 388
168 95 735 581
658 79 262 42
66 179 774 150
976 419 332 755
367 487 784 663
773 174 430 144
280 986 1000 19
21 637 659 353
582 977 371 42
219 824 120 686
153 68 274 864
880 786 756 399
280 934 330 82
870 524 382 767
158 349 484 390
245 929 33 161
679 174 975 371
489 96 593 8
972 244 957 895
485 3 986 910
555 838 189 978
439 404 38 928
548 822 529 236
277 506 832 219
918 892 978 619
`

var expectedOutputs = []string{
	"54",
	"52",
	"272",
	"98",
	"246",
	"22",
	"12",
	"60",
	"10",
	"196",
	"382",
	"4",
	"180",
	"282",
	"84",
	"4",
	"36",
	"21",
	"42",
	"9",
	"44",
	"12",
	"102",
	"95",
	"11",
	"498",
	"36",
	"30",
	"21",
	"12",
	"6",
	"20",
	"486",
	"2",
	"22",
	"210",
	"2",
	"2",
	"66",
	"1",
	"36",
	"53",
	"4",
	"28",
	"144",
	"6",
	"99",
	"310",
	"182",
	"30",
	"13",
	"162",
	"2",
	"462",
	"1",
	"30",
	"1",
	"208",
	"156",
	"20",
	"226",
	"148",
	"10",
	"58",
	"3",
	"125",
	"6",
	"172",
	"732",
	"360",
	"202",
	"64",
	"6",
	"14",
	"20",
	"102",
	"5",
	"380",
	"282",
	"112",
	"3",
	"82",
	"23",
	"10",
	"40",
	"153",
	"40",
	"10",
	"12",
	"25",
	"6",
	"3",
	"478",
	"3",
	"51",
	"419",
	"179",
	"66",
	"12",
	"390",
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
		if len(parts) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 fields, got %d", i+1, len(parts))
		}
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad a: %v", i+1, err)
		}
		b, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad b: %v", i+1, err)
		}
		m, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %v", i+1, err)
		}
		r0, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad r0: %v", i+1, err)
		}
		tests = append(tests, testCase{a: a, b: b, m: m, r0: r0})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	input := fmt.Sprintf("%d %d %d %d\n", tc.a, tc.b, tc.m, tc.r0)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

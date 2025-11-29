package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases previously stored in testcasesD.txt.
const testcasesDData = `100
979 884 891 94
87 47 57 754
829 686 764 258
621 218 528 37
596 163 383 654
403 371 403 973
381 279 335 515
275 19 33 373
477 164 358 434
914 906 914 169
574 182 302 237
25 6 16 178
140 131 139 369
527 187 415 816
425 377 410 929
931 782 875 809
608 363 455 880
985 457 622 978
773 410 646 671
544 256 506 286
948 511 767 528
852 816 838 678
905 466 702 360
582 571 582 468
499 338 394 964
333 86 310 930
632 275 520 317
311 259 294 531
520 417 456 749
213 126 191 376
957 701 739 804
841 350 721 9
930 835 859 763
109 8 81 669
51 18 32 699
897 109 881 535
140 69 100 845
216 16 124 920
735 33 91 372
369 89 216 689
25 3 6 978
70 4 9 747
941 22 404 262
131 41 64 536
709 2 396 604
45 16 20 995
38 1 23 961
631 116 408 346
501 16 173 460
565 47 317 774
412 319 409 158
485 116 163 677
704 324 376 25
459 404 459 971
131 101 116 528
336 74 248 266
269 215 256 19
717 572 607 687
59 17 19 135
166 44 56 465
651 238 498 939
726 33 285 239
732 456 493 257
83 76 79 640
811 639 731 263
701 434 576 539
769 5 159 37
394 210 251 114
525 90 213 105
103 3 26 769
238 27 82 26
534 476 505 318
549 390 444 702
929 779 832 747
826 445 662 524
22 19 19 903
429 269 417 186
943 97 776 822
492 188 197 532
984 944 951 626
376 149 325 955
382 158 162 896
702 423 474 108
314 102 300 689
846 17 479 62
421 327 389 475
214 151 160 6
292 13 203 314
958 741 760 225
774 503 601 119
586 383 483 734
475 72 457 354
405 63 193 125
126 11 89 871
343 329 335 983
218 178 184 26
634 482 493 741
723 510 584 367
983 469 613 818
384 138 261 539`

type testCase struct {
	n int64
	l int64
	r int64
	s int64
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// solve mirrors the logic from 1063D.go for a single test case.
func solve(n, l, r, s int64) int64 {
	var n1, n2 int64
	n1 = r - l + 1
	if n1 <= 0 {
		n1 += n
	}
	n2 = n - n1

	check1_t := func(t, sVal, flag int64) bool {
		denom := t + n
		temp := sVal / denom
		temp2 := sVal % denom
		if temp2 >= n1+flag && temp2 <= n1+min(t, n1) {
			rem := t - (temp2 - n1)
			if rem >= 0 && rem <= n2 {
				return true
			}
		}
		if temp != 0 {
			temp2 += denom
			if temp2 >= n1+flag && temp2 <= n1+min(t, n1) {
				rem := t - (temp2 - n1)
				if rem >= 0 && rem <= n2 {
					return true
				}
			}
		}
		return false
	}

	check1 := func(t int64) bool {
		return check1_t(t, s, 0) || check1_t(t, s+1, 1)
	}

	solve1 := func() int64 {
		for x := n; x >= 0; x-- {
			if check1(x) {
				return x
			}
		}
		return -1
	}

	div1 := func(x, y int64) int64 {
		if x == 0 {
			return 0
		}
		if x > 0 {
			return x / y
		}
		x = -x
		return -((x-1)/y + 1)
	}

	div2 := func(x, y int64) int64 {
		if x == 0 {
			return 0
		}
		if x > 0 {
			return (x-1)/y + 1
		}
		x = -x
		return -(x / y)
	}

	calc2 := func(k, sVal, flag int64) int64 {
		sTemp := sVal - k*n - n1
		l1 := div2(sTemp, k+1)
		l2 := div2(sTemp-n1, k)
		l3 := div1(sTemp-flag, k)
		l4 := div1(sTemp+n2, k+1)
		l1 = max(l1, l2)
		l3 = min(l3, l4)
		if l1 <= l3 {
			return l3
		}
		return -1
	}

	solve2 := func() int64 {
		var ans int64 = -1
		if s >= n1 && s <= n1*2 {
			ans = max(ans, s-n1+n2)
		}
		if s+1 > n1 && s+1 <= n1*2 {
			ans = max(ans, s+1-n1+n2)
		}
		for k := int64(1); k*n <= s+1; k++ {
			ans = max(ans, calc2(k, s, 0))
			ans = max(ans, calc2(k, s+1, 1))
		}
		return ans
	}

	if n < s/n {
		return solve1()
	}
	return solve2()
}

func parseTestCases(data string) ([]testCase, error) {
	tokens := strings.Fields(data)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no embedded testcases found")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if idx+3 >= len(tokens) {
			return nil, fmt.Errorf("test %d missing numbers", i+1)
		}
		var nums [4]int64
		for j := 0; j < 4; j++ {
			val, err := strconv.ParseInt(tokens[idx+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d has invalid integer: %w", i+1, err)
			}
			nums[j] = val
		}
		idx += 4
		cases = append(cases, testCase{n: nums[0], l: nums[1], r: nums[2], s: nums[3]})
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("embedded data has %d extra tokens", len(tokens)-idx)
	}
	return cases, nil
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d %d %d\n", tc.n, tc.l, tc.r, tc.s)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := fmt.Sprintf("%d", solve(tc.n, tc.l, tc.r, tc.s))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	cases, err := parseTestCases(testcasesDData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(candidate, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

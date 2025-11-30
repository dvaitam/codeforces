package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
865 395 4 42 266 989 524
498 208 3 992 489 367
598 224 5 143 289 144 774 98
634 257 5 723 830 617 924 151
318 51 1 921
871 701 3 484 574 104
363 223 3 626 656 935
210 142 4 454 887 534 267
64 2 1 737
861 409 1 627
506 424 3 250 748 334
721 65 2 940 582
228 62 2 823 557
459 47 1 328
897 521 4 112 309 565 299
724 128 5 341 835 945 554 209
987 819 5 561 602 295 456 94
611 395 3 590 248 298
189 49 2 34 628
673 267 4 71 92 696 776
134 39 1 863
83 70 4 858 723 538 283
535 242 2 917 696
604 430 5 282 462 505 677 657
718 366 1 333
628 119 4 602 646 344 866
195 63 1 750
278 60 2 381 814
175 86 4 836 64 104 802
150 57 1 837
588 548 5 697 76 28 128 651
194 156 5 123 401 94 380 854
119 5 5 23 200 985 994 190
736 127 4 216 745 820 63
960 696 1 558
436 318 1 856
267 36 2 74 663
309 180 4 185 63 516 479
41 39 1 717
401 103 3 368 927 750
873 830 2 948 739
679 193 2 724 512
950 257 2 772 280
827 755 4 124 762 831 400
458 86 2 538 69
483 270 5 730 209 987 614 92
273 183 1 804
479 289 4 329 933 238 543
808 401 2 123 972
740 319 1 640
144 96 5 729 359 739 270 48
626 186 2 537 612
823 676 3 741 822 422
293 128 5 400 867 802 110 672
493 50 1 698
670 216 1 325
168 11 4 344 419 818 760
530 96 5 670 997 338 367 287
507 438 4 220 990 322 839
485 274 2 947 683
196 165 3 168 835 493
638 545 2 852 760
789 60 1 173
711 115 5 472 372 865 36 106
768 362 1 251
120 77 3 557 871 483
539 118 4 966 63 721 97
996 621 5 903 505 914 447 458
451 139 5 37 881 34 679 657
336 204 2 193 640
965 156 2 33 586
792 125 2 718 996
595 238 2 79 194
527 87 3 602 213 797
134 23 5 923 593 277 764 30
242 77 5 292 638 779 81 717
467 284 1 456
746 329 2 293 258
539 507 2 80 746
222 124 4 640 317 976 800
496 459 3 74 595 411
928 452 5 949 791 8 726 601
140 127 2 696 711
383 304 1 10
146 43 1 791
636 66 1 820
91 21 1 150
530 176 2 81 916
753 342 1 664
695 389 2 701 346
602 54 4 140 992 870 523
68 62 2 241 893
510 259 1 948
399 125 1 692
262 44 5 751 473 193 638 83
37 6 2 785 475
804 497 2 55 506
966 694 1 949
254 115 3 591 802 670
495 205 1 396
801 787 4 935 796 152 470
761 456 1 445
500 291 5 872 459 803 273 564
502 322 1 689
560 244 3 845 871 506
`

func solve(a, b int64, xs []int64) int64 {
	ans := b
	limit := a - 1
	for _, x := range xs {
		if x > limit {
			ans += limit
		} else {
			ans += x
		}
	}
	return ans
}

type testCase struct {
	input    string
	expected string
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := []testCase{}
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		a, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad a: %w", idx+1, err)
		}
		b, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad b: %w", idx+1, err)
		}
		n, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %w", idx+1, err)
		}
		if len(parts) != 3+n {
			return nil, fmt.Errorf("line %d: expected %d values, got %d", idx+1, 3+n, len(parts))
		}
		xs := make([]int64, n)
		xStrs := make([]string, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(parts[3+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad x[%d]: %w", idx+1, i, err)
			}
			xs[i] = val
			xStrs[i] = parts[3+i]
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d %d %d\n", a, b, n)
		sb.WriteString(strings.Join(xStrs, " "))
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(a, b, xs), 10),
		})
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
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

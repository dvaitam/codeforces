package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `865 395 777 912 431 42
266 989 524 498 415 941
803 850 311 992 489 367
598 914 930 224 517 143
289 144 774 98 634 819
257 932 546 723 830 617
924 151 318 102 748 76
921 871 701 339 484 574
104 363 445 324 626 656
935 210 990 566 489 454
887 534 267 64 825 941
562 938 15 96 737 861
409 728 845 804 685 641
2 627 506 848 889 342
250 748 334 721 892 65
196 940 582 228 245 823
991 146 823 557 459 94
83 328 897 521 956 502
112 309 565 299 724 128
561 341 835 945 554 209
987 819 618 561 602 295
456 94 611 818 395 325
590 248 298 189 194 842
192 34 628 673 267 488
71 92 696 776 134 898
154 946 40 863 83 920
717 946 850 554 700 401
858 723 538 283 535 832
242 870 221 917 696 604
846 973 430 594 282 462
810 454 339 752 493 535
82 899 453 366 263 599
453 518 369 974 69 476
656 320 588 454 221 153
798 54 986 517 306 515
784 525 17 569 638 245
333 513 122 908 222 749
773 622 409 291 28 582
171 314 790 596 776 837
433 394 570 24 159 993
760 820 131 82 128 265
865 765 766 651 478 222
16 127 153 887 530 977
736 43 943 140 176 454
926 310 652 611 613 841
699 255 618 626 381 508
623 78 430 579 991 189
744 521 149 231 192 989
168 308 188 21 100 4
791 522 642 552 757 45
382 802 99 975 415 511
761 702 156 866 935 55
302 256 111 599 10 599
57 393 219 276 197 590
161 512 740 440 467 291
75 261 580 913 44 443
993 639 354 29 311 999
552 568 221 516 370 862
768 410 781 853 406 180
243 513 766 42 823 33
637 619 350 834 32 479
173 228 269 572 107 609
995 155 213 41 141 515
955 219 537 566 706 632
603 181 431 103 34 446
379 913 897 835 989 459
150 986 329 990 486 885
93 977 810 360 33 315
10 344 791 669 296 59
895 370 147 970 374 482
726 742 600 813 887 725
514 33 189 579 147 828
652 378 787 958 998 390
772 366 631 578 70 58
256 201 740 754 478 642
710 139 955 22 590 306
81 951 314 149 694 672
611 744 530 343 605 173
927 895 207 886 702 271
234 393 816 194 966 566
406 65 910 413 109 469
240 909 728 705 473 128
730 464 155 760 343 992
31 553 92 786 433 531
22 569 796 645 522 288
577 938 907 665 135 164
288 313 443 754 17 53
520 249 290 194 92 818
312 848 946 327 197 255
770 348 405 380 230 861
204 472 602 723 37 780
543 550 200 456 354 228
381 776 254 303 30 61
790 372 766 99 454 798
444 90 632 639 875 417
977 779 248 178 989 42
102 416 106 75 329 344
71 723 342 107 215 823
717 898 411 963 763 315
616 757 486 972 523 319
990 134 870 702 27 804
303 165 923 949 550 868
557 806 154 451 80 197
386 23 66 189 81 154
`

type testCase struct {
	x  int64
	y  int64
	z  int64
	t1 int64
	t2 int64
	t3 int64
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

// solve is the embedded reference logic from 1054A.go.
func solve(tc testCase) string {
	stair := abs(tc.x-tc.y) * tc.t1
	elev := (abs(tc.x-tc.z) + abs(tc.y-tc.x)) * tc.t2
	elev += 3 * tc.t3
	if elev <= stair {
		return "YES"
	}
	return "NO"
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d %d %d %d %d\n", tc.x, tc.y, tc.z, tc.t1, tc.t2, tc.t3)
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
	want := solve(tc)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(raw))
	scan.Split(bufio.ScanWords)
	vals := make([]int64, 0, 6)
	for scan.Scan() {
		v, err := strconv.ParseInt(scan.Text(), 10, 64)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	if len(vals)%6 != 0 {
		return nil, fmt.Errorf("testcase count not divisible by 6")
	}
	tests := make([]testCase, 0, len(vals)/6)
	for i := 0; i+5 < len(vals); i += 6 {
		tests = append(tests, testCase{
			x:  vals[i],
			y:  vals[i+1],
			z:  vals[i+2],
			t1: vals[i+3],
			t2: vals[i+4],
			t3: vals[i+5],
		})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

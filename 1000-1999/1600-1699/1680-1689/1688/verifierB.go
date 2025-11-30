package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve is the embedded logic from 1688B.go.
func solve(arr []int) int {
	evenCnt := 0
	hasOdd := false
	minTrailing := 31
	for _, v := range arr {
		if v%2 == 1 {
			hasOdd = true
		} else {
			evenCnt++
			tz := 0
			for v%2 == 0 {
				v /= 2
				tz++
			}
			if tz < minTrailing {
				minTrailing = tz
			}
		}
	}
	if hasOdd {
		return evenCnt
	}
	return minTrailing + len(arr) - 1
}

func runCase(exe string, arr []int, exp int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	if gotStr != fmt.Sprintf("%d", exp) {
		return fmt.Errorf("expected %d got %s", exp, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Println("failed to load embedded testcases:", err)
		os.Exit(1)
	}
	for i, arr := range tests {
		exp := solve(arr)
		if err := runCase(exe, arr, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

// Embedded copy of testcasesB.txt (t followed by cases).
const testcaseData = `
100
3 583 868 822
2 262 121
8 780 461 484 668 389 808 215 97
8 30 915 856 400 444 623 781 786
1 713
8 273 739 822 235 606 968 105 924
6 32 23 27 666 555 10
7 703 222 993 433 744 30 541
4 783 449 962 508
9 239 354 237 694 225 780 471 976 297
1 427
9 945 658 103 191 645 742 881 304 124
6 918 739 997 729 513 959
7 520 850 933 687 195 311 291
10 997 904 512 867 964 518 403 604 874 36
8 249 762 817 414 425 681 178 376
9 904 720 795 691 756 384 89 450 680
9 111 798 168 534 861 403 380 502 751
1 481
9 450 349 962 881 280 292 717 845 104
8 124 704 824 935 24 859 214 955
6 577 18 571 993 793 285
5 599 173 840 656 319
3 535 22 956
8 56 714 152 869 819 187 126 721
9 97 739 26 925 906 446 397 288 255
4 130 161 100 983
4 638 450 268 507
9 76 14 725 939 52 711 40 168 822
8 533 985 756 983 626 995 901 388
5 108 840 338 836 375
7 275 689 594 858 886 967 426
8 27 542 446 784 844 91 369 824
2 30 674
10 208 19 318 588 902 361 537 372 603 927
10 840 608 190 612 542 779 75 767 625 468
1 510
3 733 727 322
3 574 708 203
6 142 976 947 536 484 143
8 643 991 727 691 424 537 840 468
1 271
2 101 431
1 814
6 153 20 669 654 440 970
6 346 666 928 931 418 936
6 559 749 702 529 52 111
10 716 942 500 548 260 621 137 4 907 203
7 679 971 625 291 925 224 630
9 807 994 735 516 510 829 124 943 278
3 231 607 355
6 517 438 898 652 961 276
4 519 214 156 836
2 62 603
7 179 300 403 715 454 388 72
8 95 204 91 423 128 192 831 420
2 208 997
5 391 4 231 736 292
6 333 194 535 570 141 498
1 472
5 218 30 984 293 539
7 512 411 280 18 18 546 296
4 496 181 220 849
4 954 231 702 125
2 722 630
8 185 138 811 48 330 688 250 109
7 348 721 150 431 958 929 132
5 980 286 87 160 745
6 984 766 822 934 309 221
10 138 57 828 337 580 99 707 933 878 100
8 556 970 696 406 763 657 434 952
3 147 932 912
10 764 568 813 43 798 892 927 199 312 583
3 146 757 100
2 69 58
8 521 147 270 859 519 896 808 480
1 709
6 940 687 349 383 48 890
7 361 712 772 150 712 201 587
2 467 291
10 826 115 759 100 622 986 396 759 657 10
9 501 465 395 772 912 888 102 446 391
7 266 681 403 330 474 72 533
8 326 429 956 173 313 271 826 328
5 722 408 54 592 863
3 36 484 628
7 272 578 999 247 834 967 624
3 59 955 563
10 912 332 101 764 535 83 10 247 36 901
7 916 190 151 660 209 408 31
7 745 136 870 760 714 394 873
8 982 480 135 918 933 124 984 966
3 689 479 715
5 5 403 902 802 468
3 365 870 435
9 274 209 320 462 870 150 340 411 553
10 496 488 434 818 243 180 359 5 57 123
5 704 205 565 155 594
1 442
5 415 871 105 646 184
10 228 560 356 988 834 549 375 348 56 671
9 203 575 685 771 327 420 341 141 771
6 349 580 891 325 382 970
9 575 407 920 204 640 684 443 917 10
`

func loadTestcases() ([][]int, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	getInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := getInt()
	if err != nil {
		return nil, fmt.Errorf("bad count: %v", err)
	}
	tests := make([][]int, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n, err := getInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %v", caseIdx+1, err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := getInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: missing value %d: %v", caseIdx+1, i+1, err)
			}
			arr[i] = val
		}
		tests = append(tests, arr)
	}
	return tests, nil
}

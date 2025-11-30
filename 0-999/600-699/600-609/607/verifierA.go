package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `7 776 63 911 52 430 39 41 62 265 46 988 75 523 28
9 142 69 288 91 143 78 773 19 97 40 633 13 818 94 256 10 931 88
6 483 79 573 82 103 27 362 71 444 62 323 57
9 266 52 63 91 824 86 940 81 561 1 937 79 14 64 95 43 736 32
6 720 29 891 31 64 19 195 70 939 58 581 12
2 415 92 84 79
10 808 13 270 87 490 83 156 9 505 31 367 46 69 55 570 67 689 8 869 8
3 458 75 625 92 956 100
2 992 23 3 39
7 161 85 630 53 42 16 887 2 585 46 938 41 653 70
2 613 71 642 14
8 390 35 209 36 927 48 421 34 967 43 47 80 774 63 439 67
3 454 7 937 39 34 43
4 712 93 479 3 589 34 420 75
3 812 84 557 22 666 79
6 224 3 821 60 89 8 46 71 406 4 6 22
4 146 59 397 62 932 50 925 4
5 488 24 248 98 387 22 480 32 781 55
9 90 50 237 62 809 96 338 27 396 32 446 45 181 48 695 32 455 78
10 113 49 423 14 833 50 621 5 9 42 536 60 916 30 136 95 510 37 367 10
9 560 18 991 50 942 1 853 32 871 79 433 62 83 96 17 63 452 42
8 731 71 818 56 230 30 226 71 287 65 102 72 341 13 159 93
4 142 37 805 22 343 53 556 59
5 535 4 394 90 23 63 480 84 601 53
7 156 44 992 14 38 13 338 97 30 89 806 32 366 82
7 612 35 183 12 596 8 396 88 863 26 419 58 789 100
2 410 37 83 87
5 70 28 720 82 377 63 39 54 996 39
10 529 71 752 70 529 78 253 100 291 4 892 14 704 35 911 1 722 35 176 16
1 354 4
10 935 31 561 90 174 99 907 67 288 58 158 35 205 57 322 64 152 59 518 100
3 873 23 825 53 934 8
7 73 29 937 49 909 66 834 67 313 28 336 99 80 56
5 173 64 689 98 515 14 235 30 843 98
7 455 46 170 14 490 56 963 22 313 15 285 40 780 98
7 679 16 964 32 536 64 318 42 733 14 613 56 231 43
4 245 96 431 68 465 66 494 58
2 233 35 581 33
6 66 68 734 32 410 99 396 20 926 60 386 77
7 60 68 598 15 516 80 107 63 249 35 907 21 30 73
1 8 48
1 557 95
1 289 19
3 586 73 12 49 505 15
3 646 69 86 8 90 25
3 755 15 17 58 910 14
8 780 80 612 35 84 60 124 62 116 63 763 17 336 40 879 98
5 495 19 981 61 897 29 5 13 194 66
9 862 2 47 81 163 9 936 61 612 11 334 94 697 82 233 70 601 71
6 270 3 356 84 245 9 222 40 549 15 59 44
3 226 16 258 39 67 15
7 311 81 114 56 705 98 870 36 292 6 837 9 455 2
8 985 7 957 76 512 91 304 58 208 40 698 4 930 69 174 5
8 585 82 26 63 567 89 513 3 912 13 164 12 488 83 458 12
9 37 61 907 43 282 88 135 34 702 95 246 72 192 50 321 33 481 30
8 734 42 903 40 700 76 132 98 970 32 65 3 617 68 543 83
3 130 88 971 72 183 56
7 37 93 59 83 58 12 274 100 489 83 742 43 151 57
7 577 28 111 57 823 18 888 49 777 17 999 3 892 91
7 744 80 517 60 455 34 819 7 756 11 508 49 546 83
5 937 83 472 4 46 81 692 39 426 90
10 1 77 964 59 384 39 633 58 559 22 616 92 205 6 776 29 224 66 319 28
10 800 1 511 58 160 59 587 62 54 14 292 91 997 7 60 68 924 34 14 42
2 396 70 9 17
8 942 25 329 18 875 20 995 40 480 49 467 96 481 77 734 91
3 42 25 131 19 9 36
4 39 62 381 86 399 95 316 62
9 713 14 491 98 496 43 294 38 634 90 87 43 940 36 287 31 955 13
3 823 80 42 43 44 91
5 200 61 591 94 925 29 851 79 395 52
8 504 30 564 72 70 70 929 46 479 61 310 53 392 21 840 18
8 945 48 307 46 38 66 36 95 465 80 168 31 369 69 934 10
7 964 1 233 29 953 95 679 23 903 43 604 1 640 72
1 54 99
3 459 60 179 7 35 19
3 972 33 903 23 365 90
10 418 29 6 43 35 95 156 5 852 63 874 2 659 40 513 67 474 95 560 50
10 1 79 593 89 356 87 189 45 356 14 322 27 388 41 413 81 879 27 47 5
1 11 40
5 417 38 395 41 219 54 844 54 731 76
5 273 81 38 57 722 41 504 70 350 50
10 571 54 250 29 926 4 918 64 638 99 498 13 613 53 75 62 872 10 960 63
1 133 9
10 131 20 729 84 648 78 268 21 879 1 656 20 493 97 501 5 193 34 249 61
10 125 25 621 76 637 56 603 67 429 18 100 59 548 40 143 60 127 20 771 67
5 526 58 651 22 297 72 978 72 144 32
1 860 37
10 460 26 981 75 300 21 329 52 756 56 526 100 245 41 14 92 446 88 508 15
3 694 89 777 67 711 60
10 981 73 241 92 654 87 826 76 15 2 877 43 331 62 585 67 886 91 226 12
6 211 64 51 1 573 80 273 89 721 19 977 72
7 905 21 739 60 408 80 851 100 954 41 500 22 935 91
3 495 8 4 64 385 5
9 969 79 21 64 661 97 441 8 330 53 716 57 784 21 456 68 182 10
4 304 13 977 5 646 63 42 21
2 453 85 807 27
8 439 71 356 27 909 77 90 62 109 20 527 60 755 82 166 48
9 342 77 451 40 612 31 40 48 87 11 161 7 408 10 19 26 155 21
1 508 61
10 124 18 462 12 120 8 943 99 69 93 413 37 924 5 570 62 226 13 382 25
1 800 42
9 334 38 727 98 443 40 749 6 398 85 833 39 345 36 130 9 94 86
7 604 44 436 48 614 33 156 46 972 25 142 71 609 35
2 83 33 34 8
4 359 87 471 51 316 73 259 78
2 318 70 602 49
9 768 50 642 67 227 10 474 12 663 48 36 80 706 93 863 33 385 66
5 265 95 188 81 247 51 897 65 325 76
7 666 6 956 82 900 96 8 3 886 12 26 3 671 72
6 450 78 592 55 560 42 2 75 82 66 938 69
8 872 8 998 65 532 23 62 25 90 34 247 25 995 24 144 8
7 879 11 46 61 784 86 935 43 985 7 383 12 99 32
4 910 82 266 91 938 88 124 85
2 152 37 828 95
9 972 56 603 25 776 6 39 80 929 60 439 35 42 18 967 48 681 59
6 493 91 739 65 184 85 651 67 686 7 559 92
3 936 76 507 40 743 43`

type testCase struct {
	n     int
	pairs [][2]int
}

// solveCase mirrors 607A.go.
func solveCase(tc testCase) int {
	n := tc.n
	bs := make([]struct{ x, p int }, n)
	for i, pr := range tc.pairs {
		bs[i] = struct{ x, p int }{x: pr[0], p: pr[1]}
	}
	sort.Slice(bs, func(i, j int) bool { return bs[i].x < bs[j].x })
	pos := make([]int, n)
	pow := make([]int, n)
	for i := 0; i < n; i++ {
		pos[i] = bs[i].x
		pow[i] = bs[i].p
	}
	dp := make([]int, n)
	for i := 0; i < n; i++ {
		left := pos[i] - pow[i]
		j := sort.Search(i, func(k int) bool { return pos[k] >= left }) - 1
		if j >= 0 {
			dp[i] = dp[j] + (i - j - 1)
		} else {
			dp[i] = i
		}
	}
	ans := n
	for i := 0; i < n; i++ {
		destroyed := dp[i] + (n - i - 1)
		if destroyed < ans {
			ans = destroyed
		}
	}
	if n < ans {
		ans = n
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for lineIdx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", lineIdx, err)
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("line %d: expected %d values got %d", lineIdx, 1+2*n, len(fields))
		}
		pairs := make([][2]int, n)
		idx := 1
		for i := 0; i < n; i++ {
			x, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse x%d: %w", lineIdx, i+1, err)
			}
			y, err := strconv.Atoi(fields[idx+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse p%d: %w", lineIdx, i+1, err)
			}
			pairs[i] = [2]int{x, y}
			idx += 2
		}
		cases = append(cases, testCase{n: n, pairs: pairs})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, p := range tc.pairs {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expect := strconv.Itoa(solveCase(tc))
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\nInput:\n%s\n", idx+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

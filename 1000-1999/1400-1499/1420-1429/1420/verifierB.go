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

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `6 837 329 36 537 151 895
5 617 802 159 862 388
10 301 735 723 826 481 67 819 86 528 889
1 67
4 133 41 307 15
8 338 882 164 819 152 889 671 471
6 517 391 922 542 514 34
10 92 694 813 824 530 776 614 78 764 436
4 296 548 922 612
7 845 995 493 865 810 995 397
10 600 239 871 885 817 20 672 906 0 758
3 309 519 583
5 340 67 505 880 268
5 791 417 393 829 392
1 167
3 244 293 746
6 56 964 36 492 427 144
8 911 884 616 734 83 689 715 155
6 421 36 626 477 395 469
1 103
8 796 155 20 33 612 632 135 645
6 107 716 562 664 354 199
7 802 795 796 502 113 902 61
10 717 478 629 647 956 345 666 127 992 698
10 303 807 869 130 981 933 396 818 300 938
2 531 881
4 39 800 401 455
6 774 195 466 365 808 647
2 979 45
1 497
5 922 27 967 532 682
10 585 896 221 235 95 794 839 905 910 642
9 715 536 430 519 312 968 116 149 436
10 432 945 86 958 107 425 64 101 425 792
3 751 977 31
8 441 702 427 30 508 941 884 985
6 739 258 80 360 72 124
6 708 30 353 356 182 10
4 838 374 72 610
3 212 3 209
2 765 7
5 377 706 25 955 619
4 879 145 191 464
2 488 352
5 133 28 989 213 370
6 484 983 299 303 959 899
9 651 334 188 607 82 105 546 594 315
3 385 919 150
3 823 228 323
9 248 242 772 188 298 381 429 679 47
3 615 21 403
2 719 74
3 430 306 563
7 758 948 145 605 432 305 652
6 86 254 455 647 378 652
9 59 385 418 8 427 985 745 922 328
8 208 380 300 975 482 93 973 189
2 283 114
9 620 704 157 814 719 456 950 408 189
7 442 178 253 982 464 348 959
9 145 363 473 646 652 88 494 773 208
5 1 850 715 459 633
8 7 223 305 117 787 644 308 558
10 159 434 723 769 482 94 694 509 778 983
4 556 780 415 286
1 123
5 904 684 41 0 262
7 538 911 595 727 405 455 104
5 362 290 892 773 688
4 609 87 36 72
5 312 546 348 121 542
4 912 942 780 167
2 424 881
5 289 532 137 587 535
4 544 107 420 978
9 413 759 797 925 807 285 299 452 380
10 643 141 160 126 713 123 390 410 605 479
3 573 684 306
6 647 484 760 425 223 488
8 711 513 325 504 667 981 61 454
5 146 763 507 53 908
10 220 26 363 482 400 909 10 866 539 68
2 702 972
7 6 369 42 118 635 3 276
5 744 922 232 144 769
10 294 195 107 444 471 733 338 393 172 338
7 663 918 703 445 151 458 955
3 536 323 132
4 932 191 454 357
7 437 826 503 398 747 225 814
4 449 962 209 600
1 927
7 34 239 648 86 891 191 372
1 760
3 238 625 304
10 88 721 889 524 769 291 789 898 904 361
7 469 55 647 713 528 681 979
9 952 752 956 440 594 465 501 260 721
8 220 345 272 43 44 53 166 358
1 296`

type testCase struct {
	n   int
	arr []uint64
}

func solveCase(tc testCase) string {
	var cnt [32]int64
	for _, v := range tc.arr {
		if v == 0 {
			continue
		}
		k := bits.Len64(v) - 1
		cnt[k]++
	}
	var ans int64
	for _, c := range cnt {
		if c > 1 {
			ans += c * (c - 1) / 2
		}
	}
	return fmt.Sprint(ans)
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n: %v", idx+1, err)
		}
		if len(fields) != 1+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 1+n, len(fields))
		}
		arr := make([]uint64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseUint(fields[1+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = val
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&input, "%d", tc.n)
		for _, v := range tc.arr {
			fmt.Fprintf(&input, " %d", v)
		}
		input.WriteByte('\n')
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var expected strings.Builder
	for i, tc := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(solveCase(tc))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

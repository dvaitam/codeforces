package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `7 776 911 430 41 265 988 523
8 414 940 802 849 310 991 488 366
10 913 929 223 516 142 288 143 773 97 633
5 931 545 722 829 616
3 317 101 747
2 920 870
6 483 573 103 362 444 323
10 655 934 209 989 565 488 453 886 533 266
1 824
9 937 14 95 736 860 408 727 844 803
1 626
8 847 888 341 249 747 333 720 891
2 195 939
10 227 244 822 990 145 822 556 458 93 82
6 896 520 955 501 111 308
9 298 723 127 560 340 834 944 553 208
10 560 601 294 455 93 610 817 394 324 589
4 297 188 193 841
3 33 627 672
5 487 70 91 695 775
9 288 799 51 92 38 983 169 421 553
7 525 577 147 639 116 431 207
8 388 857 681 823 100 34 400 537
9 333 435 935 890 172 744 896 202 564
4 51 535 688 263
10 799 517 60 24 419 859 656 43 362 929
6 933 778 800 374 698 89
10 221 989 525 1 82 388 478 345 356 177
3 112 618 834
9 394 419 725 255 260 26 761 666 586
9 373 828 15 584 668 236 98 412 960
10 750 351 989 317 131 163 413 42 620 55
4 398 859 730 840
4 866 216 316 930
2 553 913
5 161 392 632 37 850
1 423
6 545 461 874 15 24 989
1 385
8 630 7 62 51 643 264 363 551
2 941 418
7 639 421 108 222 69 872 326
8 238 337 777 960 766 211 305 820
8 854 972 282 132 944 279 165 118
8 385 266 548 301 866 10 262 12
3 717 11 799
6 803 870 56 301 474 897
10 153 311 266 286 516 921 407 766 519 268
1 116
7 525 944 843 24 952 144 921
4 304 998 194 107
9 453 678 420 95 477 808 642 405 452
3 735 104 397
4 38 290 630 923
7 980 419 210 617 401 100 375
9 230 804 212 632 564 452 230 524 5
8 606 552 146 381 861 89 825 458
7 476 399 363 9 905 815 753
2 400 18
1 817
4 392 856 372 945
4 160 600 632 759
3 92 461 322
6 665 360 60 152 988 963
1 910
10 595 47 345 361 459 153 109 71 699 538
5 224 38 906 775 480
9 144 560 955 502 903 891 41 919 934
9 440 70 802 723 791 574 634 333 487
10 709 643 96 987 75 119 624 894 699 933
6 517 871 881 875 998 479
9 749 16 708 173 874 263 523 872 26
4 388 39 285 772
7 35 486 27 105 946 912 531
5 177 59 174 754 709
3 923 85 395
6 99 753 891 175 345 925
4 266 824 236 390
10 49 523 178 901 286 755 752 726 785 16
8 514 465 467 935 504 198 127 648
10 12 682 273 530 637 4 315 374 861 543
3 76 117 325
2 366 39
7 554 197 255 893 2 53 1000
10 991 354 561 281 280 700 278 980 735 94
5 100 883 861 335 804
8 728 277 188 123 424 136 390 57
9 1000 521 8 780 890 196 23 38 890
1 220
7 81 691 710 966 541 87 927
9 475 212 511 184 998 470 57 261 929
4 78 788 30 6
6 381 434 321 181 887 957
9 260 575 8 9 313 331 879 3 522
8 536 642 828 723 696 534 115 678
8 984 765 996 34 979 36 946 316
4 856 279 818 161
5 383 980 104 943 293
4 66 951 1000 590
8 788 577 624 125 717 897 474 893
8 563 899 812 325 866 351 619 49
5 773 970 31 216 299
7 205 794 521 677 176 368 124`

type testCase struct {
	n   int
	arr []int
}

func solveCase(tc testCase) int {
	sum := 0
	for _, v := range tc.arr {
		sum += v
	}
	return sum
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		tc := testCase{n: n, arr: make([]int, n)}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value: %v", idx+1, err)
			}
			tc.arr[i] = v
		}
		cases = append(cases, tc)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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

	for i, tc := range cases {
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for idx, v := range tc.arr {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(expected) {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

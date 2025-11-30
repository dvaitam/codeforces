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
	n   int
	arr []int
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
7 552 823 -139 -918 -470 977 47
8 -171 880 605 699 -379 982 -24 -267
10 826 859 -553 33 -715 -423 -714 547 -806 266
5 863 90 444 659 232
3 -365 -798 494
2 840 741
6 -34 146 -794 -276 -111 -353
10 311 869 -582 979 131 -24 -94 772 67 -467
1 648
9 875 -972 -809 473 720 -184 454 689 607
1 253
8 695 776 -318 -501 495 -334 441 782
2 -609 878
10 -546 -512 645 981 -709 644 112 -83 -814 -836
6 792 40 910 2 -777 -383
9 -404 447 -745 121 -319 668 888 106 -584
10 120 203 -411 -89 -813 221 634 -212 -351 178
4 -406 -624 -613 682
3 -933 254 344
5 -25 -859 -817 390 551
3 795 -694 891
1 725
2 839 432
9 399 -199 715 444 74 -436 68 662 -518
4 833 391 207 690
7 187 -437 -78 8 352 313 434
6 -832 -336 254 -764 -4 202
6 730 -611 -503 -967 498 -445
2 444 -549
6 627 -651 -319 -128 670 -873
2 603 -701
4 -908 673 175 298
9 233 393 -849 -946 -746 300 -614 241 700
10 -755 -199 -813 -242 707 -763 -926 240 -956 -602
3 470 -747 -19
4 489 639 -875 918
1 114
7 270 -793 711 -468 -857 -548 -853
5 -283 -107 -631 -875 31
8 -920 221 -794 432 -199 -592 -468 -266
8 716 847 881 166 -654 428 377 -584
1 615
3 731 -669 -299
9 -487 -760 222 887 -95 363 -642 -973 -35
7 843 165 791 41 879 -363 329
6 -205 715 346 -487 -686 148
1 -63
2 -313 513
1 114
5 -724 -509 560 952 -14
6 249 -411 379 -265 208 939
10 -729 465 -365 -206 532 -152 697 332 -835 -997
10 -607 430 -316 -673 -510 -544 305 -83 -225 454
10 791 -152 -936 -177 785 437 162 -144 581 356
1 -661
8 -870 -470 436 -678 -86 80 812 -3
9 236 547 -1000 810 -921 12 -333 -361 715
8 -898 656 684 792 995 663 -150 -615
9 972 296 -830 715 485 -733 -970 -178 944
7 -353 -994 -563 -971 469 545 -996
9 253 -800 -610 -757 245 329 -594 789 -381
5 410 -627 -795 -26 748
7 285 -834 -956 -438 871 -73 638
2 764 -475
3 338 66 673
6 -765 785 -684 -430 743 -962
1 -917
4 394 -469 143 -356
6 921 162 863 739 -914 733
10 342 12 458 319 848 -61 311 -109 -238 784
9 -635 -575 -231 202 -404 -982 -717 -691 -445
6 -309 617 -248 471 -809 -308
10 -927 -916 -448 -665 -694 194 -408 -261 -192 123
3 -400 -765 -21
4 912 -902 -370 -633
9 492 -855 -381 -175 711 -328 -388 -151 -778
2 148 860
8 -30 -310 722 634 999 663 -297 -746
8 -763 432 19 -127 -923 -382 -314 504
3 883 -659 283
10 -231 650 995 308 -822 -866 654 -827 -595 535
4 -875 -212 -984 -800
7 139 62 -407 -82 884 0 614
10 463 391 -556 -134 -829 -246 -550 -466 198 591
3 -117 -607 -266
2 -870 683
1 849
9 -76 540 386 -587 -757 18 -185 -475 -576
1 941
4 276 -701 -786 -595
8 -226 -260 119 693 -691 -786 221 -1
3 154 -169 307
7 797 67 14 391 878 818 -340
8 21 300 372 791 -587 111 248 907
4 -981 -304 445 971
6 675 -341 -928 75 -697 790
5 234 605 -681 725 -224
10 -398 470 446 652 -37 -865 638 -827 57 779
1 -865
4 -733 -917 -385 -969
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return nil, fmt.Errorf("no test data")
	}
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[j+1])
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = v
		}
		res = append(res, testCase{n: n, arr: arr})
	}
	return res, nil
}

// solve mirrors 1502A.go.
func solve(tc testCase) int {
	sum := 0
	for _, v := range tc.arr {
		sum += v
	}
	return sum
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		gotStr, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(gotStr)
		if err != nil || got != expect {
			fmt.Printf("test %d failed: expected %d got %s\n", idx+1, expect, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

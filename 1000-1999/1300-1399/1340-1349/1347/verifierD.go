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

const testcasesRaw = `100
10
718 479 630 648 957 346 667 128 993 699
10
304 808 870 131 982 934 397 819 301 939
2
532 882
4
40 801 402 456
6
775 196 467 366 809 648
2
980 46
1
498
5
923 28 968 533 683
10
586 897 222 236 96 795 840 906 911 643
9
716 537 431 520 313 969 117 150 437
10
433 946 87 959 108 426 65 102 426 793
3
752 978 32
8
442 703 428 31 509 942 885 986
6
740 259 81 361 73 125
6
709 31 354 357 183 11
4
839 375 73 611
3
213 4 210
2
766 8
5
378 707 26 956 620
4
880 146 192 465
2
489 353
5
134 29 990 214 371
6
485 984 300 304 960 900
9
652 335 189 608 83 106 547 595 316
3
386 920 151
3
824 229 324
9
249 243 773 189 299 382 430 680 48
3
616 22 404
2
720 75
3
431 307 564
7
759 949 146 606 433 306 653
6
87 255 456 648 379 653
9
60 386 419 9 428 986 746 923 329
8
209 381 301 976 483 94 974 190
2
284 115
9
621 705 158 815 720 457 951 409 190
7
443 179 254 983 465 349 960
9
146 364 474 647 653 89 495 774 209
5
2 851 716 460 634
8
8 224 306 118 788 645 309 559
10
160 435 724 770 483 95 695 510 779 984
4
557 781 416 287
1
124
5
905 685 42 1 263
7
539 912 596 728 406 456 105
5
363 291 893 774 689
4
610 88 37 73
5
313 547 349 122 543
4
913 943 781 168
2
425 882
5
290 533 138 588 536
4
545 108 421 979
9
414 760 798 926 808 286 300 453 381
10
644 142 161 127 714 124 391 411 606 480
3
574 685 307
6
648 485 761 426 224 489
8
712 514 326 505 668 982 62 455
5
147 764 508 54 909
10
221 27 364 483 401 910 11 867 540 69
2
974 985
7
745 173 9 115 382 969 921
5
68 618 224 948 499
9
773 721 368 629 451 838 217 689 957
9
299 483 171 949 796 70 161 993 840
5
178 64 771 580 861
6
833 113 847 992 920 945
8
31 654 296 119 9 342 926 134
2
974 154
2
712 645
7
450 254 1 280 474 396 361
5
407 470 454 89 38
5
64 207 488 527 3
8
14 644 533 141 841 488 229 449
5
121 324 206 676 99
2
790 495
3
867 244 364
2
904 432
9
705 100 450 67 631 278 836 100 506
5
633 661 322 264 498
1
392
1
654
10
466 832 754 425 576 399 593 720 691 554
9
191 418 561 113 428 519 506 84 997
2
562 478
7
22 341 403 164 644 644 137
4
917 862 308 791
7
879 131 5 794 969 993 18
3
201 907 343
9
849 61 637 579 142 726 84 318 42
7
984 773 446 364 797 412 272
7
111 91 381 406 311 178 459
10
46 148 533 127 997 944 567 832 295 218
3
831 232 784
10
249 312 934 989 93 90 902 860 180 139
6
138 53 28 2 340 157
5
973 910 997 140 301
9
169 184 982 341 939 514 796 756 939
6
654 339 833 21 675 10
8
861 290 170 527 984 364 53 653
2
22 11
7
966 80 209 163 744 402 29
5
216 95 894 761 676
9
99 319 179 302 942 337 613 603 501
9
260 457 297 832 147 144 971 302 583
1
80
2
397 11
3
712 986 774
10
301 668 272 703 340 919 28 972 962 269
4
323 5 146 720
8
426 229 755 165 137 884 226 392
10
924 363 390 910 351 308 35 852 201 701
1
436
3
742 993 659
8
683 28 511 514 817 680 902 729
2
470 593
7
136 437 317 555 266 191 465
5
420 211 166 7 27
4
790 917 194 241
8
551 757 900 536 489 185 888 182
9
953 852 386 842 301 137 747 832 795
7
497 783 111 58 651 321 820
9
579 101 666 485 400 878 569 950 979
5
844 62 849 638 472
1
780
7
446 507 522 422 329 560 448
10
100 489 707 117 669 976 778 733 295 106
1
424
8
371 344 773 191 368 36 956 147
6
337 474 705 516 623 149
5
736 387 882 638 405
1
405
3
725 548 905
5
561 179 311 132 24
10
330 314 832 948 368 356 548 571 985 107
4
488 586 972 176
5
304 393 351 195 587
8
330 932 415 535 348 215 56 2
10
355 385 244 908 961 315 35 580 544 997
2
795 120
5
193 873 525 69 720
9
947 85 802 373 216 339 761 959 933
9
671 162 133 782 968 377 360 620 86
6
241 315 545 900 450 366
3
160 644 574
7
156 21 601 822 503 341 678
4
904 593 797 316
5
521 746 483 751 101
8
127 350 789 354 436 492 360 933
10
780 686 358 840 4 175 53 693 143 201
9
454 651 572 733 47 651 197 700 724
7
499 191 622 641 594 610 351
2
163 721
2
534 331
5
19 304 742 40 706
4
536 703 446 60
10
954 278 65 892 978 301 294 8 858 599
7
713 500 905 964 553 210 678
7
859 75 802 435 610 818 568
7
32 883 241 981 580 78 454
10
186 103 26 999 103 236 230 192 947 676
1
310
8
992 299 97 714 774 30 111 174
3
423 877 998
10
758 860 41 339 457 669 694 885 248 811
10
382 668 101 444 237 140 807 508 915 786
4
656 431 890 150
2
140 4
7
142 721 422 9 842 978 124
7
43 804 173 188 23 244 147
4
701 319 259 786
10
588 976 180 761 41 135 693 404 583 710
9
962 555 872 267 566 801 756 35 509
7
77 774 806 382 742 498 964
2
527 816
6
365 116 42 58 214 662
9
365 137 342 182 556 416 717 332 400
10
293 610 487 398 621 36 527 610 863 806
7
63 285 936 853 990 76 391
8
989 596 721 687 837 757 925 53
6
616 425 380 629 89 16
5
81 829 765 586 963
8
933 92 304 856 806 868 930 544
10
591 347 325 825 946 497 682 529 908 604
5
168 441 785 746 784
6
624 448 863 916 556 857
8
666 929 518 192 102 218 356 758
6
520 149 214 669 689 811
7
336 453 650 998 111 51 302
9
499 717 511 516 843 467 850 74 741
1
789
2
50 241
7
912 203 409 254 896 55 183
10
800 668 282 732 643 724 114 893 957 430
8
121 88 84 247 102 847 784 158
7
906 942 219 452 628 79 838
7
573 773 908 850 404 41 982
3
256 502 226
3
861 892 286
6
328 446 110 571 922 293
10
555 808 207 729 304 797 453 527 620 474
9
650 268 280 238 17 122 630 803 729
2
177 751
7
255 224 292 901 754 676 7`

type testCase struct {
	n   int
	arr []int
}

func solve(tc testCase) (int, int, int) {
	l, r := 0, tc.n-1
	moves := 0
	alice, bob := 0, 0
	prev := 0
	aliceTurn := true
	for l <= r {
		cur := 0
		if aliceTurn {
			for l <= r && cur <= prev {
				cur += tc.arr[l]
				l++
			}
			alice += cur
		} else {
			for l <= r && cur <= prev {
				cur += tc.arr[r]
				r--
			}
			bob += cur
		}
		prev = cur
		moves++
		aliceTurn = !aliceTurn
	}
	return moves, alice, bob
}

func parseTestcases(raw string) ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("empty test data")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid t")
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing n for case %d", i+1)
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid n for case %d", i+1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("missing value %d for case %d", j+1, i+1)
			}
			v, err := strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("invalid value %d for case %d", j+1, i+1)
			}
			arr[j] = v
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}
	return tests, nil
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse tests:", err)
		os.Exit(1)
	}

	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(tests)))
	expectedTokens := make([]string, 0, len(tests)*3)
	for _, tc := range tests {
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')

		m, a, b := solve(tc)
		expectedTokens = append(expectedTokens, strconv.Itoa(m), strconv.Itoa(a), strconv.Itoa(b))
	}

	got, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "candidate failed:", err)
		os.Exit(1)
	}
	fields := strings.Fields(got)
	if len(fields) != len(expectedTokens) {
		fmt.Fprintf(os.Stderr, "wrong number of outputs: expected %d got %d\n", len(expectedTokens), len(fields))
		os.Exit(1)
	}
	for i, exp := range expectedTokens {
		if fields[i] != exp {
			caseIdx := i / 3
			fmt.Fprintf(os.Stderr, "case %d mismatch at token %d: expected %s got %s\n", caseIdx+1, i%3, exp, fields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

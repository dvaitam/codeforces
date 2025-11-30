package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `3 228 283 802
9 -408 -973 225 -698 -281 365 -267 -249 699
5 27 314 409 -481 -727
1 -19
4 246 -188 -305 -31
5 -99 -715 -27 676 -394
6 -286 747 306 235 -221 -190
10 -179 424 833 27 -390 381 668 550 156 727
4 -614 -117 -557 -457
9 -658 43 -713 654 223 59 549 843 88
9 712 692 326 -26 878 -301 -470 -113 -383
8 856 332 78 -894 616 494 -603 982
5 -61 415 -860 -963 -588
10 -789 752 -475 625 931 -998 690 789 -557 924
3 265 -645 857
6 248 662 -452 -982 381 344
1 -958
10 876 901 -658 702 782 -365 970 683 -517 -196
8 512 -257 757 -274 -821 192 -116 601
2 432 -950
5 283 -601 961 6 -512
6 -365 -545 670 930 -560 971
2 -883 19
6 722 556 365 -840 -41 -127
5 -437 795 688 526 543
9 933 -425 -488 757 68 706 296 25 270
4 390 17 -809 76
1 -559
10 768 -944 -942 -968 -232 438 -126 -3 -144 358
10 251 82 170 437 157 -285 429 335 731 226
6 435 935 -710 -143 953 100
10 -77 -672 -530 707 -542 764 610 29 -926 991
8 -615 901 -594 927 167 474 -633 563
2 913 -645
6 -262 142 821 518 23 -961
8 46 11 866 -86 114 -242 50 807
1 490
4 -391 -267 110 -223
2 -600 102
6 97 592 -508 -816 811 101
5 -175 -115 777 799 -85
4 857 -566 288 889
8 704 760 293 -956 -428 723 -748 222
5 835 -122 987 233 121
1 688
2 -624 429
2 442 461
1 231
3 -660 765 983
6 610 389 -924 -328 704 587
3 425 970 17
6 824 47 -240 -400 265 959
6 728 -987 -578 -501 -833 -859
3 -523 -518 -459
6 -297 151 -941 763 20 319
8 -88 955 663 -712 194 614 766 710
10 -14 -448 -126 -471 797 937 620 20 331 -734
1 -59
8 53 198 496 14 578 632 -117 -390
8 -181 -78 886 375 988 -888 539 -619
3 -942 162 986
8 -631 996 876 -161 2 -436 235 163
7 942 -140 -410 427 488 802 174
7 167 465 320 421 424 -997 -823
5 -456 -880 -567 89 807
6 -581 -303 531 993 -586 -595
10 469 72 -72 -720 418 850 -523 329 544 58
6 -500 -830 502 372 -256 -528
6 -151 -530 -764 487 11 -866
10 -185 176 -308 -854 922 549 78 417 798 523
8 -705 -165 935 -982 493 -412 -479 -366
9 501 693 258 -616 -297 -265 -709 746 298
2 50 988
2 902 -77
9 388 -7 -304 246 -779 2 -426 -749 38
8 1000 645 -768 -533 211 418 767 618
10 204 -230 117 46 67 160 675 -621 685 -413
1 427
6 -751 -213 570 -523 131 994
10 -444 826 600 362 269 985 48 974 416 553
9 -698 71 451 560 209 -417 814 -450 -485
4 887 -462 281 -271
9 -810 7 835 -708 -754 459 -836 -302 163
9 86 638 -141 -414 532 -60 -878 -659 8
1 -457
5 449 746 729 402 -622
4 874 42 -371 101
5 -651 -548 187 -385 170
5 859 -1000 169 -363 -465
9 -95 -158 278 -533 525 -927 -116 798 -630
7 687 117 -643 -302 -25 -980 -131
1 -555
4 -754 -540 780 -129
3 -640 821 -452
2 -49 -43
7 825 912 340 415 180 129 541
3 -360 990 226
9 -296 857 -972 463 140 700 862 -816 -603
7 15 -384 739 -746 1000 -455 439
8 27 -152 538 302 -849 -335 -874 805
10 1000 711 822 693 -603 -717 354 -11 -6 917
7 598 696 -431 -925 -22 947 -378
8 -446 162 -699 490 566 -626 -882 -704
4 336 616 409 -443
3 -622 596 519
3 -182 85 -635
6 544 254 -616 -680 487 110
6 119 953 -13 -25 968 -510
7 -451 844 949 674 22 -548 444
2 -572 -668
3 -485 162 -403
9 -936 288 891 -997 80 89 929 594 -80`

func expected(arr []int64) string {
	if len(arr) == 0 {
		return "0"
	}
	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}
	return fmt.Sprintf("%d", min)
}

func loadCases() ([]string, []string) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var inputs []string
	var expects []string
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != n+1 {
			fmt.Printf("invalid test line %d\n", idx+1)
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				fmt.Printf("invalid number on line %d\n", idx+1)
				os.Exit(1)
			}
			arr[i] = v
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		inputs = append(inputs, sb.String())
		expects = append(expects, expected(arr))
	}
	return inputs, expects
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expects[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expects[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}

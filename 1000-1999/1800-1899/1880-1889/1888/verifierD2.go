package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod int64 = 1000000007

const testcasesRaw = `1 -339
8 619 134 527 785 423 15 -25 -836
1 96
7 619 605 758 895 -463 -944 327
9 -804 713 692 -836 -317 -266 956 -799 -35
1 -685
9 281 648 -413 965 -925 -999 -229 -309 -679
9 425 -700 -672 -643 596 885 -675 310 392
4 931 799 913 279
6 936 -952 -11 963 818 281
7 -910 -537 -508 301 -424 -328 -652
4 -279 -538 -665 798
7 -52 -256 795 157 -718 -209 156
1 -668
10 -989 406 -202 471 641 -648 -696 -961 -949 687
6 44 -993 -922 -904 591 895
2 173 252
3 594 -687 379
7 613 -946 -143 -107 163 871 405
6 456 -493 889 -718 -255 44
4 92 -181 -851 -735
7 159 348 -282 978 -803 -117 834
7 -497 -36 -219 -540 876 -192 -509
8 -187 885 190 -859 711 -484 -436 812
9 -239 110 -956 237 252 596 -29 795 -515
5 -918 251 -342 776 634
7 974 282 -782 94 797 746 -901
3 461 -197 -946
7 497 782 999 -195 -114 893 -786
8 237 -54 -670 -653 -303 -29 -159 -675
10 880 770 -419 560 31 -770 911 -246 -293 -708
6 569 -32 290 62 543 -890
4 796 960 -473 -638
10 -336 -394 -222 299 -916 -404 129 -117 -929 397
7 -449 -232 490 -574 -290 -726 -730
2 255 -269
3 -937 -119 181
7 -51 -845 314 455 448 394 -841
7 115 487 134 -717 -654 -687 -576
3 -533 -940 75
3 677 9 -269
10 507 -411 467 624 -314 401 -758 643 971 872
7 583 767 -464 841 -676 -299 286
9 983 -321 91 423 -704 868 -230 515 523
9 -375 -520 803 873 -225 -291 -201 942 -35
9 -373 975 -161 -168 664 -800 399 993 -685
3 -990 194 219
9 -782 427 595 325 429 -581 220 319 235
9 -773 -454 410 498 258 -655 -231 -826 624
1 -979
2 827 -259
8 -349 -779 384 -74 -245 196 426 -481
8 663 588 -532 889 -675 150 129 694
2 598 753
9 885 814 -656 -948 694 317 739 -659 799
9 -163 245 353 -574 -103 450 400 -171 -460
1 203
3 883 -205 -650
8 157 -891 876 -232 903 643 -823 319
10 -171 -315 -523 37 788 -68 -920 -18 269 -792
5 706 39 4 128 320
7 -32 -457 -630 -530 107 -254 -674
5 745 246 703 877 -706
8 -860 971 -860 -21 -195 159 951 -161
9 -812 -457 -13 -530 668 -773 -393 -712 -264
2 -717 348
1 685
3 201 144 -594
1 -933
7 143 433 546 232 4 462 413
2 907 -32
9 -289 791 597 -298 787 -801 404 -990 -506
4 19 717 -363 -438
4 -982 15 -274 939
9 897 -302 -813 -842 -376 174 -136 -537 513
6 973 -220 762 552 -703 -526
5 916 -592 520 589 -15
6 -421 -213 254 -736 702 594
2 -179 -273
9 778 -36 -532 343 452 -235 294 875 -266
7 -432 -265 992 -180 571 569 463
5 774 -785 -18 850 -404
2 -87 -685
6 910 -502 522 -618 903 957
6 22 -532 -775 394 -214 997
7 -49 51 -55 916 1000 158 696
10 -549 383 -183 28 -363 -4 -524 -356 60 400
1 -812
8 -351 -193 824 945 -535 680 981 -117
1 178
1 -168
2 -461 -590
6 -635 -769 -630 632 430 -255
1 -522
1 -984
7 94 -991 -735 -762 682 238 240
4 -833 517 -46 -596
1 65
10 732 -150 -864 99 -638 896 -521 -536 901 -150
7 -28 653 -1000 -108 -570 -220 998`

func expected(nums []int64) string {
	prod := int64(1)
	for _, v := range nums {
		val := v % mod
		if val < 0 {
			val += mod
		}
		prod = (prod * val) % mod
	}
	return fmt.Sprintf("%d", prod%mod)
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
		nums := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				fmt.Printf("invalid number on line %d\n", idx+1)
				os.Exit(1)
			}
			nums[i] = v
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		inputs = append(inputs, sb.String())
		expects = append(expects, expected(nums))
	}
	return inputs, expects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
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

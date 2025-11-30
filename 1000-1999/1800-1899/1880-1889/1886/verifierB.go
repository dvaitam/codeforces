package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
-725 165 735 643 564 -871
-478 -759 14 558 -80 -33
334 -223 615 -571 -808 -1
-942 829 711 -202 -114 244
561 571 -996 425 -88 -455
477 642 -532 210 935 -791
846 -350 -938 -955 -948 330
108 -982 923 804 -220 405
-557 984 -136 486 -941 80
-546 564 -104 923 15 132
-523 -293 -528 386 -552 558
-59 950 -407 897 -956 -148
715 876 139 888 315 -796
-620 288 482 761 -393 -753
521 -319 834 477 993 456
25 917 980 -136 39 699
864 372 -612 -379 -419 203
993 807 22 733 926 34
-195 206 747 -930 -17 -503
523 633 -173 -152 361 -646
-249 123 807 439 588 381
511 -233 -823 -102 359 41
-779 594 -665 66 720 -195
-242 2 500 -940 -39 -911
-369 440 737 259 214 184
-194 325 -652 -655 28 -536
-975 578 -592 105 884 761
122 -525 -172 52 -296 950
735 183 -277 -60 863 -449
350 122 247 960 493 -989
-215 604 755 680 955 814
921 516 49 657 -736 62
592 149 -580 -128 945 -886
-15 781 -254 167 135 -591
927 33 -154 -7 665 -270
-152 -292 -997 102 106 276
610 254 -322 -62 228 -943
647 -530 301 -638 127 196
-630 763 -813 635 128 632
743 672 906 -478 -934 723
932 378 -856 -830 777 -966
-73 -971 544 547 -425 -489
-450 -776 632 279 -622 -295
-406 -858 -658 -674 -478 80
949 -656 344 -442 327 457
-397 -69 438 -341 16 -30
-767 -952 -362 -209 -297 -138
630 -615 -471 -778 -481 842
495 44 1000 -572 977 240
-116 673 997 -958 -539 -964
-187 -701 -928 472 965 -672
-88 443 36 388 -127 115
704 -549 1000 998 291 633
423 57 -77 -543 72 328
-938 -192 382 179 645 -343
351 292 -127 -880 510 -389
-743 982 -566 793 -903 -373
-856 758 -844 -365 878 923
-390 523 -676 -148 156 -484
-733 -983 148 799 741 -923
209 678 -555 970 844 167
-57 -649 695 777 781 994
597 441 275 42 -924 -226
-590 -290 -798 -579 174 380
836 -114 211 -603 8 -787
920 363 -202 -394 32 23
-965 -334 253 785 -177 842
-424 -963 -679 -589 756 -329
661 153 602 -724 -306 -121
-564 -455 381 -803 715 -224
909 121 -296 872 807 715
407 94 -8 572 90 -520
-867 485 -918 -827 -728 -653
-659 865 102 -564 -452 554
-320 229 36 722 -478 -247
-307 -304 -767 -404 -519 776
933 236 596 954 464 817
1 -723 187 128 577 -787
-344 -920 -168 -851 -222 773
614 -699 696 -744 -302 -766
259 203 601 896 -226 -844
168 126 -542 159 -833 950
-454 -253 824 -395 155 94
894 -766 -63 836 -433 -780
611 -907 695 -395 -975 256
373 -971 -813 -154 -765 691
812 617 -919 -616 -510 608
201 -138 -669 -764 -77 -658
394 -506 -675 523 730 -790
-109 865 975 -226 651 987
111 862 675 -398 126 -482
457 -24 -356 -795 -575 335
-350 -919 -945 -979 611 895
-395 487 221 -345 -79 -199
-359 -184 -872 -869 870 -351
986 231 986 -67 -772 -488
-560 607 265 593 825 111
776 409 -40 355 -272 -470
-625 109 -575 -371 -593 -496
-262 -834 679 -425 -817 542`

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

func expected(px, py, ax, ay, bx, by float64) string {
	dPA := (ax-px)*(ax-px) + (ay-py)*(ay-py)
	dOA := ax*ax + ay*ay
	dPB := (bx-px)*(bx-px) + (by-py)*(by-py)
	dOB := bx*bx + by*by

	ans := math.Min(math.Max(dPA, dOA), math.Max(dPB, dOB))
	actual := math.Sqrt(ans)

	dAB := math.Hypot(ax-bx, ay-by)
	r1 := math.Max(dAB/2, math.Max(math.Sqrt(dOA), math.Sqrt(dPB)))
	r2 := math.Max(dAB/2, math.Max(math.Sqrt(dOB), math.Sqrt(dPA)))
	if r1 < actual {
		actual = r1
	}
	if r2 < actual {
		actual = r2
	}
	return fmt.Sprintf("%.20f", actual)
}

func loadCases() ([]string, []string) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Println("invalid testcase count")
		os.Exit(1)
	}
	expectedTokens := 1 + t*6
	if len(tokens) != expectedTokens {
		fmt.Println("embedded testcase count mismatch")
		os.Exit(1)
	}
	var inputs []string
	var expects []string
	pos := 1
	for i := 0; i < t; i++ {
		vals := make([]float64, 6)
		for j := 0; j < 6; j++ {
			v, err := strconv.ParseFloat(tokens[pos+j], 64)
			if err != nil {
				fmt.Printf("invalid value on case %d\n", i+1)
				os.Exit(1)
			}
			vals[j] = v
		}
		pos += 6
		inputs = append(inputs, fmt.Sprintf("1\n%s\n", strings.Join(tokens[pos-6:pos], " ")))
		expects = append(expects, expected(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]))
	}
	return inputs, expects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expects[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expects[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}

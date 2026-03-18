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

// Correct embedded solver for 2C: Apollonius point.
func solveC(x1, y1, r1, x2, y2, r2, x3, y3, r3 float64) string {
	if r1 == r2 && r2 == r3 {
		B1 := -2 * (x1 - x2)
		C1 := -2 * (y1 - y2)
		D1 := x1*x1 + y1*y1 - (x2*x2 + y2*y2)

		B2 := -2 * (x2 - x3)
		C2 := -2 * (y2 - y3)
		D2 := x2*x2 + y2*y2 - (x3*x3 + y3*y3)

		det := B1*C2 - B2*C1
		if math.Abs(det) > 1e-9 {
			X := (C1*D2 - C2*D1) / det
			Y := (D1*B2 - D2*B1) / det
			return fmt.Sprintf("%.5f %.5f", X, Y)
		}
		return ""
	}

	alpha := r3*r3 - r2*r2
	beta := r1*r1 - r3*r3
	gamma := r2*r2 - r1*r1

	U := 2.0 * (alpha*x1 + beta*x2 + gamma*x3)
	V := 2.0 * (alpha*y1 + beta*y2 + gamma*y3)
	W := -(alpha*(x1*x1 + y1*y1) + beta*(x2*x2 + y2*y2) + gamma*(x3*x3 + y3*y3))

	A12 := r2*r2 - r1*r1
	B12 := -2 * (r2*r2*x1 - r1*r1*x2)
	C12 := -2 * (r2*r2*y1 - r1*r1*y2)
	D12 := r2*r2*(x1*x1 + y1*y1) - r1*r1*(x2*x2 + y2*y2)

	A23 := r3*r3 - r2*r2
	B23 := -2 * (r3*r3*x2 - r2*r2*x3)
	C23 := -2 * (r3*r3*y2 - r2*r2*y3)
	D23 := r3*r3*(x2*x2 + y2*y2) - r2*r2*(x3*x3 + y3*y3)

	A31 := r1*r1 - r3*r3
	B31 := -2 * (r1*r1*x3 - r3*r3*x1)
	C31 := -2 * (r1*r1*y3 - r3*r3*y1)
	D31 := r1*r1*(x3*x3 + y3*y3) - r3*r3*(x1*x1 + y1*y1)

	var A, B, C, D float64
	if math.Abs(A12) >= math.Abs(A23) && math.Abs(A12) >= math.Abs(A31) {
		A, B, C, D = A12, B12, C12, D12
	} else if math.Abs(A23) >= math.Abs(A12) && math.Abs(A23) >= math.Abs(A31) {
		A, B, C, D = A23, B23, C23, D23
	} else {
		A, B, C, D = A31, B31, C31, D31
	}

	var qa, qb, qc float64
	var solveX bool
	if math.Abs(V) >= math.Abs(U) {
		solveX = true
		m := -U / V
		k := -W / V
		qa = A * (1 + m*m)
		qb = A*2*m*k + B + C*m
		qc = A*k*k + C*k + D
	} else {
		solveX = false
		m := -V / U
		k := -W / U
		qa = A * (1 + m*m)
		qb = A*2*m*k + B*m + C
		qc = A*k*k + B*k + D
	}

	delta := qb*qb - 4*qa*qc
	maxVal := math.Max(qb*qb, math.Abs(4*qa*qc))
	epsD := 1e-10 * maxVal
	if delta < -epsD {
		return ""
	}
	if delta < 0 {
		delta = 0
	}

	root1 := (-qb + math.Sqrt(delta)) / (2 * qa)
	root2 := (-qb - math.Sqrt(delta)) / (2 * qa)

	var cands [][2]float64
	if solveX {
		m := -U / V
		k := -W / V
		cands = append(cands, [2]float64{root1, m*root1 + k})
		cands = append(cands, [2]float64{root2, m*root2 + k})
	} else {
		m := -V / U
		k := -W / U
		cands = append(cands, [2]float64{m*root1 + k, root1})
		cands = append(cands, [2]float64{m*root2 + k, root2})
	}

	bestRatio := math.MaxFloat64
	var bestCand [2]float64
	found := false

	for _, cand := range cands {
		cx, cy := cand[0], cand[1]
		dist1 := math.Sqrt((cx-x1)*(cx-x1) + (cy-y1)*(cy-y1))
		dist2 := math.Sqrt((cx-x2)*(cx-x2) + (cy-y2)*(cy-y2))
		dist3 := math.Sqrt((cx-x3)*(cx-x3) + (cy-y3)*(cy-y3))

		ratio1 := dist1 / r1
		ratio2 := dist2 / r2
		ratio3 := dist3 / r3

		if math.Abs(ratio1-ratio2) > 1e-4*math.Max(1.0, ratio1) || math.Abs(ratio1-ratio3) > 1e-4*math.Max(1.0, ratio1) {
			continue
		}

		if ratio1 < bestRatio {
			bestRatio = ratio1
			bestCand = cand
			found = true
		}
	}

	if found {
		return fmt.Sprintf("%.5f %.5f", bestCand[0], bestCand[1])
	}
	return ""
}

// Embedded copy of testcasesC.txt to remove external dependency.
const testcasesC = `957 767 971 738 -885 94 -827 -261 856
-654 507 829 371 748 316 -485 240 218
242 -927 596 395 -676 442 307 -195 823
480 761 522 944 -239 558 916 -89 515
-451 845 37 783 -944 373 -48 908 327
-453 178 360 935 -968 830 304 341 34
422 -606 304 -418 -159 51 231 675 100
418 694 738 1 521 721 626 505 503
-758 -730 690 -793 -899 69 -989 627 146
-810 -315 672 -151 -620 273 898 676 564
-231 -443 515 -618 -471 808 -857 155 400
80 48 593 397 -709 690 262 263 167
28 807 459 808 525 676 -523 23 36
-258 691 851 -692 85 519 -86 702 502
920 38 477 848 -16 516 -806 254 949
-261 -736 50 250 496 86 -54 -317 66
692 -363 196 -415 305 924 -848 853 197
892 -72 56 958 0 605 -664 180 259
953 206 580 -82 -150 524 258 -922 943
-435 592 46 642 572 72 -240 279 667
505 -703 221 184 -927 563 -174 -844 932
869 543 749 69 858 638 -932 370 280
759 332 182 -816 -904 879 -875 281 420
481 617 189 -434 735 284 -980 419 327
224 495 65 805 159 547 -370 637 471
-698 -228 794 -724 131 768 -8 224 514
659 -555 162 -450 836 167 -131 -100 411
551 544 712 405 -94 731 457 286 385
295 -876 953 245 -753 969 230 -285 282
-614 -179 378 229 -541 775 -46 918 158
355 -804 599 -533 -87 730 -473 934 217
-134 511 962 483 -448 417 -221 386 535
587 275 75 1 -285 447 -259 179 956
388 71 540 654 243 304 584 -998 634
830 954 89 978 -850 147 550 -410 772
-632 907 977 -441 -262 638 192 -905 4
-706 -600 615 -698 987 329 686 -386 69
-850 -871 606 756 -224 194 -793 -629 153
67 558 238 -748 -382 293 -136 85 826
-114 235 683 -245 -910 262 -38 -897 8
911 236 629 324 -627 300 -462 -238 871
1000 663 71 -771 668 424 -385 -441 280
151 -603 137 996 -72 973 430 -377 124
-672 860 694 -634 -434 640 146 422 889
-340 -219 305 -999 -170 620 -116 965 311
875 548 121 28 -625 782 372 723 888
-368 188 183 -976 798 418 -954 368 661
321 925 970 -448 973 617 -438 -792 493
951 -488 925 499 203 792 689 -22 437
10 959 471 -612 688 581 377 -870 47
-25 -266 204 -461 508 282 -536 313 428
331 588 237 213 -546 197 -21 671 177
336 -800 541 802 -154 240 694 38 752
-779 -614 208 -136 632 272 -613 -783 79
-7 1000 891 -943 -268 8 4 392 269
-832 756 45 -767 -68 212 568 -759 498
-402 -970 279 -484 901 554 288 -863 494
-321 -375 609 706 573 519 346 164 35
-148 -309 634 768 997 444 742 64 413
-958 342 920 340 400 522 -242 -450 316
-835 42 156 48 927 704 932 -307 548
944 -772 344 694 -783 104 -75 -115 955
598 461 468 -274 893 467 44 874 319
-300 -772 114 -876 -969 276 289 473 378
832 771 437 -243 -2 235 30 381 394
852 -45 885 -596 297 40 958 -722 627
-902 -921 425 -780 508 593 205 666 405
-192 -99 606 -83 -97 25 138 943 953
-7 832 941 717 899 409 914 -636 35
-553 85 33 355 150 824 -416 -406 985
-77 -717 317 -398 432 11 534 63 493
422 807 794 -59 -136 277 776 -1000 857
878 -889 369 -814 19 491 -967 -889 470
590 -977 132 499 566 712 -413 122 682
115 -806 764 -602 -873 605 906 -175 946
797 498 291 67 -786 364 187 -535 510
659 -895 836 178 71 621 -789 -440 963
479 -158 136 166 465 411 207 265 305
360 769 656 404 -685 506 -501 546 999
702 -498 621 -954 391 741 463 835 510
92 -209 434 -12 698 81 504 939 211
-787 625 672 774 -670 726 -791 -724 722
-810 482 817 863 -624 249 -166 -807 479
983 66 10 -243 940 240 457 563 781
-711 874 139 760 363 132 262 935 468
-700 493 512 -192 -268 963 160 -947 472
443 -15 625 985 -845 103 -25 30 663
-923 602 22 557 -188 101 256 623 539
331 -859 613 -564 -62 657 77 865 354
-2 182 888 612 201 639 -570 885 320
971 -547 678 -137 516 970 -145 820 128
-259 -491 485 -772 -869 277 497 20 822
563 -820 379 -192 -28 509 63 529 478
742 -962 113 415 -644 342 92 -247 105
97 17 100 213 6 267 763 -196 168
-965 -210 220 -593 1 442 -231 -791 35
-303 -888 112 -135 -43 156 625 -879 522
430 873 928 -792 115 475 622 -50 453
29 -146 539 -520 -68 726 590 -163 55
919 226 253 487 210 825 445 25 251
-369 -538 787 -878 815 176 202 -738 528
789 -810 779 -686 778 371 258 488 409`

type circle struct{ x, y, r float64 }

type testCase struct {
	c [3]circle
}

func parseCases() ([]testCase, error) {
	data := strings.TrimSpace(testcasesC)
	lines := strings.Split(data, "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 9 {
			return nil, fmt.Errorf("line %d: expected 9 numbers, got %d", i+1, len(fields))
		}
		var tc testCase
		for j := 0; j < 3; j++ {
			x, err1 := strconv.ParseFloat(fields[j*3], 64)
			y, err2 := strconv.ParseFloat(fields[j*3+1], 64)
			r, err3 := strconv.ParseFloat(fields[j*3+2], 64)
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("line %d: parse error", i+1)
			}
			tc.c[j] = circle{x, y, r}
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func expected(tc testCase) string {
	return solveC(tc.c[0].x, tc.c[0].y, tc.c[0].r,
		tc.c[1].x, tc.c[1].y, tc.c[1].r,
		tc.c[2].x, tc.c[2].y, tc.c[2].r)
}

func runCandidate(bin, input string) (string, error) {
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
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Println("failed to load testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d %d\n%d %d %d\n%d %d %d\n",
			int(tc.c[0].x), int(tc.c[0].y), int(tc.c[0].r),
			int(tc.c[1].x), int(tc.c[1].y), int(tc.c[1].r),
			int(tc.c[2].x), int(tc.c[2].y), int(tc.c[2].r))
		expect := expected(tc)
		got, err := runCandidate(binary, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		// Compare numerically with tolerance rather than exact string match,
		// because different floating-point approaches may diverge in the last digit.
		if expect == "" && got == "" {
			continue
		}
		if (expect == "") != (got == "") {
			fmt.Printf("Test %d failed:\ninput:\n%s\nexpected:\n%q\ngot:\n%q\n", idx+1, input, expect, got)
			os.Exit(1)
		}
		var ex, ey, gx, gy float64
		if _, err := fmt.Sscanf(expect, "%f %f", &ex, &ey); err != nil {
			fmt.Printf("Test %d: failed to parse expected %q: %v\n", idx+1, expect, err)
			os.Exit(1)
		}
		if _, err := fmt.Sscanf(got, "%f %f", &gx, &gy); err != nil {
			fmt.Printf("Test %d: failed to parse got %q: %v\n", idx+1, got, err)
			os.Exit(1)
		}
		const tol = 1e-3
		if math.Abs(ex-gx) > tol || math.Abs(ey-gy) > tol {
			fmt.Printf("Test %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

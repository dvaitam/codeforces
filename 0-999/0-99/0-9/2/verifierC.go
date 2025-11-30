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

// Geometry primitives from 2C.go.
const eps = 1e-8

type point struct{ x, y float64 }

func (p point) sub(a point) point   { return point{p.x - a.x, p.y - a.y} }
func (p point) add(a point) point   { return point{p.x + a.x, p.y + a.y} }
func (p point) mul(k float64) point { return point{p.x * k, p.y * k} }
func (p point) len() float64        { return math.Hypot(p.x, p.y) }

type circle struct{ x, y, r float64 }

func (c circle) o() point { return point{c.x, c.y} }

type line struct{ a, b, c float64 } // ax + by + c = 0
type carrier struct {
	typ int
	c   circle
	l   line
} // 0 circle, 1 line

func crossLL(l1, l2 line) []point {
	det := l1.a*l2.b - l1.b*l2.a
	if math.Abs(det) < eps {
		return nil
	}
	det1 := -(l1.c*l2.b - l1.b*l2.c)
	det2 := -(l1.a*l2.c - l1.c*l2.a)
	return []point{{det1 / det, det2 / det}}
}

func crossCL(c circle, l line) []point {
	var res []point
	al, be := l.b, -l.a
	var x0, y0 float64
	if math.Abs(l.a) < math.Abs(l.b) {
		x0 = 0
		y0 = -l.c / l.b
	} else {
		y0 = 0
		x0 = -l.c / l.a
	}
	A := al*al + be*be
	B := 2*al*(x0-c.x) + 2*be*(y0-c.y)
	Cq := (x0-c.x)*(x0-c.x) + (y0-c.y)*(y0-c.y) - c.r*c.r
	D := B*B - 4*A*Cq
	if D < -eps {
		return nil
	}
	if D < 0 {
		D = 0
	}
	t1 := (-B + math.Sqrt(D)) / (2 * A)
	res = append(res, point{x0 + al*t1, y0 + be*t1})
	t2 := (-B - math.Sqrt(D)) / (2 * A)
	res = append(res, point{x0 + al*t2, y0 + be*t2})
	return res
}

func crossCLgen(cl1, cl2 carrier) []point {
	if cl1.typ == 0 && cl2.typ == 0 {
		c1, c2 := cl1.c, cl2.c
		a := 2 * (c2.x - c1.x)
		b := 2 * (c2.y - c1.y)
		c0 := c2.r*c2.r - c1.r*c1.r + c1.x*c1.x - c2.x*c2.x + c1.y*c1.y - c2.y*c2.y
		return crossCL(c1, line{a, b, c0})
	}
	if cl1.typ == 0 && cl2.typ == 1 {
		return crossCL(cl1.c, cl2.l)
	}
	if cl1.typ == 1 && cl2.typ == 0 {
		return crossCL(cl2.c, cl1.l)
	}
	if cl1.typ == 1 && cl2.typ == 1 {
		return crossLL(cl1.l, cl2.l)
	}
	return nil
}

func getL(c1, c2 circle) carrier {
	a := 2*c2.x - 2*c1.x
	b := 2*c2.y - 2*c1.y
	c0 := c1.x*c1.x - c2.x*c2.x + c1.y*c1.y - c2.y*c2.y
	return carrier{typ: 1, l: line{a, b, c0}}
}

func getC(c1, c2 circle) carrier {
	if c1.r > c2.r {
		return getC(c2, c1)
	}
	cr := c1.r / c2.r
	o1 := c1.o()
	o2 := c2.o()
	v := o2.sub(o1)
	p1 := o1.add(v.mul(cr / (1 + cr)))
	p2 := o1.add(v.mul(cr / (cr - 1)))
	o := p1.add(p2).mul(0.5)
	r := p1.sub(o).len()
	return carrier{typ: 0, c: circle{o.x, o.y, r}}
}

func getCL(c1, c2 circle) carrier {
	if math.Abs(c1.r-c2.r) < eps {
		return getL(c1, c2)
	}
	return getC(c1, c2)
}

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
	cl1 := getCL(tc.c[0], tc.c[1])
	cl2 := getCL(tc.c[1], tc.c[2])
	cr := crossCLgen(cl1, cl2)
	mi := 1e100
	var ans point
	for _, p := range cr {
		var q [3]float64
		ok := true
		for j := 0; j < 3; j++ {
			q[j] = p.sub(tc.c[j].o()).len() / tc.c[j].r
			if math.Abs(q[j]-q[0]) > eps {
				ok = false
			}
		}
		if q[0] < 1-eps {
			ok = false
		}
		if !ok {
			continue
		}
		if q[0] < mi {
			mi = q[0]
			ans = p
		}
	}
	if mi < 1e50 {
		return fmt.Sprintf("%.5f %.5f", ans.x, ans.y)
	}
	return ""
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
		if got != expect {
			fmt.Printf("Test %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

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

// Embedded testcases (previously in testcasesB.txt) to keep verifier self contained.
const rawTestcasesB = `
395 866 431 778
17 43 524 990
208 499 803 942
311 851 489 993
299 368 224 915
143 518 72 290
98 775 257 635
546 933 617 724
151 925 51 319
76 749 871 922
339 702 287 485
46 105 162 446
210 627 566 991
227 490 534 888
32 268 562 826
15 939 93 97
409 862 685 729
2 642 506 628
342 849 187 251
33 335 146 197
62 229 146 824
557 824 47 460
41 84 521 898
502 957 39 113
299 566 128 725
341 562 554 836
205 210 561 619
295 603 47 457
395 612 295 326
75 249 49 190
192 843 17 35
36 489 87 93
134 777 154 899
40 947 83 864
717 921 850 947
401 555 723 859
283 539 242 536
221 871 696 918
430 605 282 595
253 463 657 678
366 719 42 86
119 629 301 500
344 647 195 867
5 250 278 751
91 121 96 227
175 815 219 342
64 837 101 105
57 151 37 48
548 651 76 618
4 29 194 652
590 622 51 124
48 95 119 855
2 39 48 201
127 737 108 492
63 746 696 961
18 25 318 437
34 105 29 73
39 75 224 360
16 186 479 517
39 42 90 105
103 402 184 268
750 928 430 483
584 925 173 175
197 210 51 61
163 694 166 867
272 352 61 258
453 613 180 683
8 15 420 699
583 923 521 897
319 941 366 666
337 399 79 258
13 576 380 470
43 82 47 758
288 559 62 140
494 782 313 362
184 296 136 606
318 734 384 399
425 426 83 668
1 3 343 717
62 165 164 230
194 460 690 729
425 583 26 34
719 894 429 583
679 792 48 728
115 171 34 67
162 720 271 459
499 908 575 931
1 620 40 907
167 508 240 321
52 53 832 844
97 427 86 563
743 859 4 135
348 413 162 429
2 5 12 16
3 774 692 844
101 543 31 197
204 624 310 896
94 288 61 104
407 876 84 644
9 24 464 937
`

type testCase struct {
	a, b, c, d int
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(rawTestcasesB)
	if len(fields)%4 != 0 {
		return nil, fmt.Errorf("unexpected token count %d (want multiple of 4)", len(fields))
	}
	cases := make([]testCase, 0, len(fields)/4)
	for i := 0; i < len(fields); i += 4 {
		a, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("parse a at token %d (%q): %w", i+1, fields[i], err)
		}
		b, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("parse b at token %d (%q): %w", i+2, fields[i+1], err)
		}
		c, err := strconv.Atoi(fields[i+2])
		if err != nil {
			return nil, fmt.Errorf("parse c at token %d (%q): %w", i+3, fields[i+2], err)
		}
		d, err := strconv.Atoi(fields[i+3])
		if err != nil {
			return nil, fmt.Errorf("parse d at token %d (%q): %w", i+4, fields[i+3], err)
		}
		cases = append(cases, testCase{a: a, b: b, c: c, d: d})
	}
	return cases, nil
}

// expected probability calculation mirrors 312B.go.
func expected(a, b, c, d int) float64 {
	p1 := float64(a) / float64(b)
	p2 := float64(c) / float64(d)
	return p1 / (1 - (1-p1)*(1-p2))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		want := expected(tc.a, tc.b, tc.c, tc.d)
		input := fmt.Sprintf("%d %d %d %d\n", tc.a, tc.b, tc.c, tc.d)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
			fmt.Printf("case %d: failed to parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if math.Abs(got-want) > 1e-6 {
			fmt.Printf("case %d failed: expected %.10f got %.10f\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

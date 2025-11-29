package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
100
865 395 388 912 430
42 266 261 498 207
941 803 310 992 488
367 598 223 517 142
289 144 24 634 256
932 546 150 318 50
748 76 42 484 286
104 363 222 324 312
656 935 209 990 565
489 454 443 534 266
64 825 561 938 14
96 737 408 728 684
641 2 1 848 341
250 748 333 721 64
196 940 581 228 61
823 991 145 823 556
459 94 10 328 260
956 502 55 309 282
299 724 127 561 340
835 945 553 209 204
618 561 294 456 46
611 818 394 325 294
248 298 94 194 47
34 628 266 488 35
92 696 133 898 153
946 40 5 920 716
946 850 553 700 400
858 723 537 283 267
832 242 217 221 173
604 846 429 594 281
462 505 338 657 365
85 333 313 119 62
602 646 343 866 194
249 17 8 120 90
226 381 87 341 218
836 64 12 802 149
876 715 224 47 36
650 932 547 617 75
28 128 48 621 589
123 401 46 380 59
38 621 22 200 47
736 127 61 216 186
820 63 59 696 23
558 436 317 104 33
72 227 18 663 308
359 447 92 63 32
479 41 38 104 89
401 205 66 368 240
859 924 583 174 172
209 990 785 60 50
693 163 41 351 271
257 121 76 944 452
682 180 3 483 348
420 922 582 896 520
940 319 182 398 336
257 158 143 708 12
469 760 80 344 23
558 288 69 246 195
977 494 180 625 294
690 368 302 970 913
649 875 635 136 79
398 767 424 849 666
83 2 0 716 342
164 246 57 653 458
388 728 689 582 424
33 412 359 582 428
791 679 47 170 114
66 266 80 458 270
907 499 464 575 0
906 40 31 334 159
858 479 25 829 425
193 562 85 858 742
134 16 12 973 694
428 324 1 219 3
735 773 2 843 691
542 627 100 196 30
623 665 203 895 309
287 706 186 103 60
875 945 406 643 83
23 282 231 820 811
119 883 262 137 133
837 667 660 356 58
893 159 71 872 19
44 42 13 698 265
572 323 187 961 581
932 870 43 867 767
985 719 622 672 506
730 660 469 656 445
382 893 550 183 53
385 602 298 10 2
155 278 170 346 188
736 96 43 799 635
37 43 17 168 38
598 297 184 405 280
133 301 58 490 374
246 957 49 316 91
878 536 72 310 206
856 337 153 425 55
102 575 492 486 172
862 818 351 128 122
`

type testCase struct {
	n int64
	a int64
	b int64
	c int64
	d int64
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func solve(tc testCase) string {
	minTotal := tc.n * (tc.a - tc.b)
	maxTotal := tc.n * (tc.a + tc.b)
	if minTotal <= tc.c+tc.d && maxTotal >= tc.c-tc.d {
		return "Yes"
	}
	return "No"
}

func parseTestcases() ([]testCase, error) {
	rawLines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, ln := range rawLines {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			lines = append(lines, ln)
		}
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	start := 0
	count := 0
	if fields := strings.Fields(lines[0]); len(fields) == 1 {
		if v, err := strconv.Atoi(fields[0]); err == nil {
			count = v
			start = 1
		}
	}
	if count == 0 {
		count = len(lines)
		start = 0
	}
	if start+count != len(lines) {
		return nil, fmt.Errorf("testcase count mismatch: declared %d actual %d", count, len(lines)-start)
	}
	cases := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		line := strings.TrimSpace(lines[start+i])
		if line == "" {
			return nil, fmt.Errorf("empty line at case %d", i+1)
		}
		fields := strings.Fields(line)
		if len(fields) != 5 {
			return nil, fmt.Errorf("case %d expected 5 numbers got %d", i+1, len(fields))
		}
		vals := make([]int64, 5)
		for j := 0; j < 5; j++ {
			v, err := strconv.ParseInt(fields[j], 10, 64)
			if err != nil {
				return nil, err
			}
			vals[j] = v
		}
		cases = append(cases, testCase{n: vals[0], a: vals[1], b: vals[2], c: vals[3], d: vals[4]})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("1\n%d %d %d %d %d\n", tc.n, tc.a, tc.b, tc.c, tc.d)
		want := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

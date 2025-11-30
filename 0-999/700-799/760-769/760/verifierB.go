package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `100
138 611178140 17
262 126614504 254
780 482638132 484
668 407609409 215
97 523832193 4
915 897396863 400
444 652232025 391
786 2262139 713
457 285970713 370
822 245632386 606
968 109766543 924
326 32846077 12
27 697444882 18
10 946217664 7
703 232572534 433
744 31183049 541
228 820017927 113
962 532375303 567
239 371193231 60
694 234915040 471
976 311151610 949
23 446869829 18
945 689659269 103
191 675762725 186
881 318247645 124
761 357229494 739
997 763637346 513
959 453234901 520
850 977304617 687
195 325739658 73
602 947555211 512
867 542545236 403
604 916211566 36
492 260640548 381
817 434101856 425
681 185765967 376
562 947826855 384
89 471331550 85
521 115890830 168
534 901891637 403
380 525804795 376
31 503928697 2
316 755251083 315
608 620812259 404
663 182911722 173
515 243672628 13
790 214229858 553
943 923730056 562
238 434280342 132
353 909954663 296
362 492989300 138
676 588407230 624
981 783188468 6
393 841443791 380
525 868807879 133
532 834724398 211
437 60262371 247
891 391633238 584
568 214576506 517
424 520684794 417
366 444985305 178
2 578187202 2
470 644090071 15
824 246537347 651
182 591370218 150
186 924501412 24
818 591663846 817
872 876643821 262
34 903816658 5
86 932091843 3
464 15634118 387
774 301933408 256
276 117562791 95
353 311690776 36
172 171396774 66
541 180544801 280
664 764064535 302
466 754438907 165
509 508708263 59
25 335012767 13
352 451958337 97
265 116782245 130
922 783995907 523
215 650310492 111
837 22354111 231
19 426614150 5
37 771843743 11
457 756564988 260
695 458128775 558
853 236868027 646
817 746306032 529
462 239655102 269
665 32964834 405
692 618310583 329
676 677475783 437
61 791832311 20
129 227774800 13
314 75942714 40
318 984810881 153
762 169874654 427`

func calc(x, length int64) int64 {
	if length <= x-1 {
		return (x - 1 + x - length) * length / 2
	}
	return (x-1)*x/2 + (length - (x - 1))
}

func feasible(n, m, k, x int64) bool {
	left := calc(x, k-1)
	right := calc(x, n-k)
	return x+left+right <= m
}

func expected(n, m, k int64) string {
	l, r := int64(1), m
	var ans int64
	for l <= r {
		mid := (l + r) / 2
		if feasible(n, m, k, mid) {
			ans = mid
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return fmt.Sprintf("%d", ans)
}

func runCase(exe, input, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

type testCase struct {
	n int64
	m int64
	k int64
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	if len(lines)-1 < t {
		return nil, fmt.Errorf("expected %d cases, found %d", t, len(lines)-1)
	}
	out := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		fields := strings.Fields(lines[i+1])
		if len(fields) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields got %d", i+2, len(fields))
		}
		n, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", i+2, err)
		}
		m, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %w", i+2, err)
		}
		k, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %w", i+2, err)
		}
		out = append(out, testCase{n: n, m: m, k: k})
	}
	return out, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to load embedded testcases:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k)
		exp := expected(tc.n, tc.m, tc.k) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

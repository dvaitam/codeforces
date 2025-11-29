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

const testcasesData = `
884 451
685 743
676 673
232 266
353 687
167 312
18 364
588 552
60 749
646 154
362 22
503 642
63 25
248 46
13 231
969 669
335 68
841 63
354 995
680 432
140 913
945 221
460 445
146 366
320 181
665 336
746 804
766 418
392 9
420 912
271 546
545 823
754 703
722 473
778 42
579 943
126 418
977 399
176 3
513 141
637 889
678 526
861 750
717 151
82 337
244 860
849 840
181 252
969 22
980 825
918 172
761 806
859 700
575 172
736 80
438 949
886 612
107 635
644 468
728 153
630 616
42 258
349 825
756 750
386 27
641 954
911 37
509 91
367 299
688 154
470 241
520 364
167 753
772 414
346 276
823 504
968 402
16 318
544 922
296 563
481 1000
36 791
545 584
566 268
957 703
40 466
405 738
123 413
355 507
53 21
280 757
36 260
698 696
596 719
795 951
297 702
781 212
781 541
529 348
396 854
`

type testCase struct {
	n int
	k int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func positions(n, k, seg int) []int {
	rem := seg - (k+1)*2
	start := 1 + min(rem, k)
	step := 2*k + 1
	var pos []int
	for p := start; p <= n; p += step {
		pos = append(pos, p)
	}
	return pos
}

func expected(n, k int) string {
	pos := solve(n, k)
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(pos))
	for j, v := range pos {
		if j > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	return b.String()
}

func solve(n, k int) []int {
	if n < (k+1)*2 {
		return []int{min(n, k+1)}
	}
	step := 2*k + 1
	baseMin := (k + 1) * 2
	for i := baseMin; i <= baseMin+step; i++ {
		if (n-i)%step == 0 {
			return positions(n, k, i)
		}
	}
	pos := make([]int, n)
	for i := 0; i < n; i++ {
		pos[i] = i + 1
	}
	return pos
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expect) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expect, got)
	}
	return nil
}

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	var cases []testCase
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{n: n, k: k})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
		expect := expected(tc.n, tc.k)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

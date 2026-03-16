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

func solveA(a, b int64) string {
	if a < b {
		return "-1"
	}
	k := (a + b) / (2 * b)
	ans := float64(a+b) / (2 * float64(k))
	return fmt.Sprintf("%.12f", ans)
}

func run(bin string, input string) (string, error) {
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

const testcasesARaw = `138 583
868 822
783 65
262 121
508 780
461 484
668 389
808 215
97 500
30 915
856 400
444 623
781 786
3 713
457 273
739 822
235 606
968 105
924 326
32 23
27 666
555 10
962 903
391 703
222 993
433 744
30 541
228 783
449 962
508 567
239 354
237 694
225 780
471 976
297 949
23 427
858 939
570 945
658 103
191 645
742 881
304 124
761 341
918 739
997 729
513 959
991 433
520 850
933 687
195 311
291 602
997 904
512 867
964 518
403 604
874 36
492 249
762 817
414 425
681 178
376 562
904 720
795 691
756 384
89 450
680 521
111 798
168 534
861 403
380 502
751 31
481 45
316 721
869 630
608 593
404 663
175 173
515 233
13 790
205 553
943 881
562 238
415 527
353 976
868 592
362 471
932 276
676 562
624 981
747 6
393 803
878 841
978 908
961 759
525 829
133 532
797 575
211 437
973 58
493 891
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "invalid test case format on line %d\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.ParseInt(parts[0], 10, 64)
		b, _ := strconv.ParseInt(parts[1], 10, 64)
		expected := solveA(a, b)
		input := fmt.Sprintf("%d %d\n", a, b)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		// allow small floating point error
		if expected == "-1" {
			if got != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\n", idx, got)
				os.Exit(1)
			}
			continue
		}
		valExp, _ := strconv.ParseFloat(expected, 64)
		valGot, err := strconv.ParseFloat(got, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: non-float output %s\n", idx, got)
			os.Exit(1)
		}
		if diff := valExp - valGot; diff > 1e-6 || diff < -1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

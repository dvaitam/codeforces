package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases from testcasesA.txt.
const embeddedTestcasesA = `137 582
867 821
782 64
261 120
507 779
460 483
667 388
807 214
96 499
29 914
855 399
443 622
780 785
2 712
456 272
738 821
234 605
967 104
923 325
31 22
26 665
554 9
961 902
390 702
221 992
432 743
29 540
227 782
448 961
507 566
238 353
236 693
224 779
470 975
296 948
22 426
857 938
569 944
657 102
190 644
741 880
303 123
760 340
917 738
996 728
512 958
990 432
519 849
932 686
194 310
290 601
996 903
511 866
963 517
402 603
873 35
491 248
761 816
413 424
680 177
375 561
903 719
794 690
755 383
88 449
679 520
110 797
167 533
860 402
379 501
750 30
480 44
315 720
868 629
607 592
403 662
174 172
514 232
12 789
204 552
942 880
561 237
414 526
352 975
867 591
361 470
931 275
675 561
623 980
746 5
392 802
877 840
977 907
960 758
524 828
132 531
796 574
210 436
972 57
492 890`

func solve(n, m int) int {
	count := 0
	for a := 0; a*a <= n; a++ {
		b := n - a*a
		if b >= 0 && a+b*b == m {
			count++
		}
	}
	return count
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcasesA))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, m int
		if _, err := fmt.Sscan(line, &n, &m); err != nil {
			fmt.Fprintf(os.Stderr, "invalid line %d: %v\n", idx, err)
			os.Exit(1)
		}
		expect := solve(n, m)
		input := fmt.Sprintf("%d %d\n", n, m)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(outStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse output %q\n", idx, outStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d (n=%d m=%d)\n", idx, expect, got, n, m)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

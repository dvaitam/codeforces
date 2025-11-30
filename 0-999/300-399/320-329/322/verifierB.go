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

const testcasesB = `100
137 582 867
821 782 64
261 120 507
779 460 483
667 388 807
214 96 499
29 914 855
399 443 622
780 785 2
712 456 272
738 821 234
605 967 104
923 325 31
22 26 665
554 9 961
902 390 702
221 992 432
743 29 540
227 782 448
961 507 566
238 353 236
693 224 779
470 975 296
948 22 426
857 938 569
944 657 102
190 644 741
880 303 123
760 340 917
738 996 728
512 958 990
432 519 849
932 686 194
310 290 601
996 903 511
866 963 517
402 603 873
35 491 248
761 816 413
424 680 177
375 561 903
719 794 690
755 383 88
449 679 520
110 797 167
533 860 402
379 501 750
30 480 44
315 720 868
629 607 592
403 662 174
172 514 232
12 789 204
552 942 880
561 237 414
526 352 975
867 591 361
470 931 275
675 561 623
980 746 5
392 802 877
840 977 907
960 758 524
828 132 531
796 574 210
436 972 57
492 890 373
583 567 204
963 516 423
496 832 365
424 354 1
551 553 638
805 627 339
469 614 28
823 235 650
181 563 598
185 881 93
817 564 816
871 836 953
261 33 861
966 689 72
85 888 17
463 14 772
773 287 255
275 112 816
639 189 352
297 71 171
163 261 540
974 172 672
279 663 728
301 465 719
329 508 485
116 24 319
395 351 431
815 192 264
111 259 921
747 522 1000
214 988 620
442 836 998
21 230 18`

func expected(r, g, b int64) int64 {
	base := r/3 + g/3 + b/3
	best := base
	for k := int64(1); k <= 2; k++ {
		if r < k || g < k || b < k {
			break
		}
		cur := k + (r-k)/3 + (g-k)/3 + (b-k)/3
		if cur > best {
			best = cur
		}
	}
	return best
}

func runCase(bin string, r, g, b int64) error {
	input := fmt.Sprintf("%d %d %d\n", r, g, b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	want := expected(r, g, b)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scan := bufio.NewScanner(strings.NewReader(testcasesB))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "empty test data")
		os.Exit(1)
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad test count: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing r\n", i+1)
			os.Exit(1)
		}
		rVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing g\n", i+1)
			os.Exit(1)
		}
		gVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing b\n", i+1)
			os.Exit(1)
		}
		bVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if err := runCase(bin, rVal, gVal, bVal); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}

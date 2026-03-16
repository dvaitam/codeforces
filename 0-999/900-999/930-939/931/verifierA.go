package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(a, b int) int {
	if a > b {
		a, b = b, a
	}
	diff := b - a
	x := diff / 2
	y := diff - x
	return x*(x+1)/2 + y*(y+1)/2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `100
865 395
777 912
431 42
266 989
524 498
415 941
803 850
311 992
489 367
598 914
930 224
517 143
289 144
774 98
634 819
257 932
546 723
830 617
924 151
318 102
748 76
921 871
701 339
484 574
104 363
445 324
626 656
935 210
990 566
489 454
887 534
267 64
825 941
562 938
15 96
737 861
409 728
845 804
685 641
2 627
506 848
889 342
250 748
334 721
892 65
196 940
582 228
245 823
991 146
823 557
459 94
83 328
897 521
956 502
112 309
565 299
724 128
561 341
835 945
554 209
987 819
618 561
602 295
456 94
611 818
395 325
590 248
298 189
194 842
192 34
628 673
267 488
71 92
696 776
134 898
154 946
40 863
83 920
717 946
850 554
700 401
858 723
538 283
535 832
242 870
221 917
696 604
846 973
430 594
282 462
505 677
657 718
939 813
366 85
333 628
119 499
602 646
344 866
195 249
17 750`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Printf("test %d missing a\n", i+1)
			os.Exit(1)
		}
		var a int
		fmt.Sscan(scan.Text(), &a)
		if !scan.Scan() {
			fmt.Printf("test %d missing b\n", i+1)
			os.Exit(1)
		}
		var b int
		fmt.Sscan(scan.Text(), &b)
		exp := expected(a, b)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		var got int
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Printf("case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

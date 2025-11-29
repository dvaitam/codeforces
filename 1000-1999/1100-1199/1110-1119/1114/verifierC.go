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

const testcasesRaw = `905036 972
890299 59
96034 88
378597 857
177298 755
848259 687
895311 317
263805 622
222528 623
37471 597
714339 164
451590 655
412649 824
758134 882
533796 974
390134 559
981165 457
526456 276
944985 38
913345 30
381697 478
977112 328
951845 391
444189 915
927009 540
172479 575
186056 243
241805 26
185305 334
182022 141
534949 524
377164 528
707244 575
190677 917
467285 817
434815 754
550886 930
952609 783
381950 810
622315 364
379466 881
467423 167
790631 411
749891 758
483820 672
556119 257
513817 287
969757 512
525170 529
871917 817
371117 679
925737 467
943405 926
483407 361
595282 745
965037 572
758931 469
510247 676
232586 965
340439 836
733555 857
174137 899
952045 633
281163 793
955649 493
324601 312
838084 725
871438 518
589498 532
532003 669
645722 604
426425 321
766436 214
512656 526
384407 958
717470 640
924917 79
822378 842
358046 745
8831 931
854652 197
780962 110
61614 590
684626 52
286366 607
237625 700
918020 939
111428 774
547737 141
895425 274
256725 846
220704 968
923357 63
443462 921
751789 779
33422 60
379959 370
180231 257
705446 26
86931 119`

type testCase struct {
	n int64
	b int64
}

// calc returns exponent of prime k in n!
func calc(n, k int64) int64 {
	var res int64
	for n > 0 {
		res += n / k
		n /= k
	}
	return res
}

// solve embeds logic from 1114C.go.
func solve(n, b int64) int64 {
	bb := b
	ans := int64(math.MaxInt64)
	for i := int64(2); i*i <= bb; i++ {
		if bb%i == 0 {
			cnt := int64(0)
			for bb%i == 0 {
				cnt++
				bb /= i
			}
			val := calc(n, i) / cnt
			if val < ans {
				ans = val
			}
		}
	}
	if bb > 1 {
		val := calc(n, bb)
		if val < ans {
			ans = val
		}
	}
	return ans
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.b)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	want := strconv.FormatInt(solve(tc.n, tc.b), 10)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		n, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{n: n, b: b})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

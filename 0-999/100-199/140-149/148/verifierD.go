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

// solve mirrors 148D.go.
func solve(w, b int) float64 {
	var ans, v float64
	cnt := 0
	total := w + b
	for i := 0; cnt < total; i++ {
		if i&1 == 1 {
			cnt++
		} else {
			ans += (1 - v) * float64(w) / float64(total-i)
		}
		cnt++
		v += (1 - v) * float64(w) / float64(total-i)
	}
	return ans
}

type testCase struct {
	w int
	b int
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
243 606
557 133
378 937
618 485
640 594
67 620
13 930
857 480
265 564
239 196
734 481
553 856
562 487
406 654
881 154
237 650
155 888
948 535
399 759
15 687
795 65
163 776
980 605
43 308
798 31
843 886
275 484
609 736
942 899
396 731
807 943
437 404
745 820
590 455
987 958
137 899
374 99
36 139
506 222
264 988
688 446
797 641
875 308
431 519
853 395
587 359
546 599
417 598
237 925
344 698
937 951
29 876
286 620
687 712
167 715
881 334
987 554
926 585
582 106
730 671
216 648
851 587
273 291
127 64
493 874
654 495
90 352
819 68
420 918
154 20
300 437
787 425
893 121
45 619
629 779
46 386
735 600
338 564
902 944
285 517
241 36
317 7
78 110
614 548
32 971
202 994
417 298
625 269
159 706
43 888
347 321
368 981
141 918
882 386
385 471
890 532
395 659
887 609
697 572
105 635
996 963
830 519
277 441
649 737
732 243
958 308
447 264
533 310
561 347
11 807
425 593
322 20
385 630
603 647
136 61
648 642
340 477
361 695
939 361
623 723
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields)%2 != 0 {
		return nil, fmt.Errorf("malformed test data")
	}
	res := make([]testCase, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		w, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("bad w at pair %d: %v", i/2+1, err)
		}
		b, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("bad b at pair %d: %v", i/2+1, err)
		}
		res = append(res, testCase{w: w, b: b})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func almostEqual(aStr, bStr string) bool {
	a, err1 := strconv.ParseFloat(aStr, 64)
	b, err2 := strconv.ParseFloat(bStr, 64)
	if err1 != nil || err2 != nil {
		return aStr == bStr
	}
	return math.Abs(a-b) <= 1e-9
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.w, tc.b)
		expected := fmt.Sprintf("%.15f", solve(tc.w, tc.b))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if !almostEqual(strings.TrimSpace(got), strings.TrimSpace(expected)) {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

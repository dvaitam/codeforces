package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors the intended logic for 1449A (sum of two integers).
func solve(a, b int) string {
	return strconv.Itoa(a + b)
}

type testCase struct {
	a int
	b int
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
100
137 582
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
492 890
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	pos := 1
	for i := 0; i < t; i++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d missing values", i+1)
		}
		a, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad a: %v", i+1, err)
		}
		b, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad b: %v", i+1, err)
		}
		res = append(res, testCase{a: a, b: b})
		pos += 2
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("testcase count mismatch: read %d tokens, have %d", pos, len(fields))
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

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
		expected := solve(tc.a, tc.b)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

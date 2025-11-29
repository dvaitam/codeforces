package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
138 583
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
493 891`

// gcd returns greatest common divisor.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

type testCase struct {
	a int
	b int
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	if len(fields) != 1+2*t {
		return nil, fmt.Errorf("expected %d numbers got %d", 1+2*t, len(fields))
	}
	tests := make([]testCase, t)
	idx := 1
	for i := 0; i < t; i++ {
		a, _ := strconv.Atoi(fields[idx])
		b, _ := strconv.Atoi(fields[idx+1])
		tests[i] = testCase{a: a, b: b}
		idx += 2
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d\n", tc.a, tc.b)
}

func runBinary(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	expected := strconv.Itoa(gcd(tc.a, tc.b))
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
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

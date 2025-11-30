package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (previously in testcasesA.txt) to keep verifier self contained.
const rawTestcasesA = `
100
865
395
777
912
431
42
266
989
524
498
415
941
803
850
311
992
489
367
598
914
930
224
517
143
289
144
774
98
634
819
257
932
546
723
830
617
924
151
318
102
748
76
921
871
701
339
484
574
104
363
445
324
626
656
935
210
990
566
489
454
887
534
267
64
825
941
562
938
15
96
737
861
409
728
845
804
685
641
2
627
506
848
889
342
250
748
334
721
892
65
196
940
582
228
245
823
991
146
823
557
`

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

// solve235A mirrors 235A.go so the verifier has no external oracle.
func solve235A(n int64) int64 {
	switch {
	case n <= 2:
		return n
	case n == 3:
		return 6
	default:
		if n%2 != 0 {
			return n * (n - 1) * (n - 2)
		}
		if n%3 != 0 {
			return lcm(lcm(n, n-1), n-3)
		}
		return (n - 1) * (n - 2) * (n - 3)
	}
}

func loadTestcases() ([]int64, error) {
	lines := strings.Fields(rawTestcasesA)
	nums := make([]int64, 0, len(lines))
	for idx, s := range lines {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d (%q): %w", idx+1, s, err)
		}
		nums = append(nums, v)
	}
	return nums, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, n := range testcases {
		expect := solve235A(n)
		input := fmt.Sprintf("%d\n", n)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("case %d: failed to parse output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

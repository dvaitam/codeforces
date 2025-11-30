package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Original solution logic (always panics for any input n > 0).
func run971A(n int) {
	arr := make([]int, n-1)
	// This will panic for any n >= 0 (slice bounds).
	_ = arr[n]
}

// Embedded testcases (from testcasesA.txt) so the verifier is self contained.
const rawTestcases = `865
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

func parseTestcases() ([]int, error) {
	lines := strings.Split(rawTestcases, "\n")
	var cases []int
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse int: %w", idx+1, err)
		}
		cases = append(cases, v)
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	// For every testcase, the reference behavior is a runtime panic (non-zero exit).
	for idx, n := range cases {
		_ = n // embedded logic above demonstrates expected panic
		input := fmt.Sprintf("%d\n", n)
		_, err := run(bin, input)
		if err == nil {
			fmt.Fprintf(os.Stderr, "case %d: expected runtime error but binary exited successfully\n", idx+1)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}

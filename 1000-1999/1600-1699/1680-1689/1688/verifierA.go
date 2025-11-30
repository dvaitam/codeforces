package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve is the embedded logic from 1688A.go.
func solve(x int) int {
	if x == 1 {
		return 3
	}
	if x&1 == 1 {
		return 1
	}
	if x&(x-1) == 0 {
		return x + 1
	}
	return x & -x
}

func runCase(exe string, x int, exp int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("1\n%d\n", x))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	if gotStr != fmt.Sprintf("%d", exp) {
		return fmt.Errorf("expected %d got %s", exp, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Println("failed to load embedded testcases:", err)
		os.Exit(1)
	}
	for i, x := range tests {
		exp := solve(x)
		if err := runCase(exe, x, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

// Embedded copy of testcasesA.txt (first line was the count).
const testcaseData = `
100
1
2
4
8
16
32
64
128
256
512
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
`

func loadTestcases() ([]int, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	cnt, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad count: %v", err)
	}
	if len(fields)-1 != cnt {
		return nil, fmt.Errorf("testcase count mismatch: declared %d have %d", cnt, len(fields)-1)
	}
	tests := make([]int, 0, cnt)
	for i, f := range fields[1:] {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("bad value at %d: %v", i+1, err)
		}
		tests = append(tests, v)
	}
	return tests, nil
}

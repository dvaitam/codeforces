package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `865
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
557`

func solveCase(n int64) string {
	if n%2 == 1 || n < 4 {
		return "-1"
	}
	min := (n + 5) / 6
	max := n / 4
	return fmt.Sprintf("%d %d", min, max)
}

func parseTestcases() ([]int64, error) {
	fields := strings.Fields(strings.TrimSpace(testcasesRaw))
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	cases := make([]int64, 0, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse case %d: %v", i+1, err)
		}
		cases = append(cases, v)
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, n := range cases {
		expected := solveCase(n)
		input := fmt.Sprintf("1\n%d\n", n)

		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

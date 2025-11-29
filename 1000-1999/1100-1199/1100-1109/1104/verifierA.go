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
81
890
52
702
690
119
367
315
806
136
890
559
36
805
69
644
103
620
780
604
15
467
64
665
118
402
128
845
63
93
329
636
294
101
215
358
901
98
326
464
571
981
794
792
677
281
654
451
277
278
640
154
347
235
294
862
108
6
150
485
72
881
434
546
896
140
738
317
447
282
100
873
320
451
845
266
917
847
40
518`

// solve mirrors 1104A.go: find highest divisor of n in [9..1], output count and digit.
func solve(n int) (int, int) {
	div := 1
	q := n
	for i := 9; i >= 1; i-- {
		if n%i == 0 {
			div = i
			q = n / i
			break
		}
	}
	return q, div
}

type testCase struct {
	n int
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	var cases []testCase
	for scan.Scan() {
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{n: v})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d\n", tc.n)
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
	gotStr, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	tokens := strings.Fields(gotStr)
	q, div := solve(tc.n)
	if len(tokens) != q+1 {
		return fmt.Errorf("expected %d tokens got %d", q+1, len(tokens))
	}
	first, err := strconv.Atoi(tokens[0])
	if err != nil || first != q {
		return fmt.Errorf("expected count %d got %q", q, tokens[0])
	}
	for i := 1; i < len(tokens); i++ {
		v, err := strconv.Atoi(tokens[i])
		if err != nil || v != div {
			return fmt.Errorf("expected digit %d got %q", div, tokens[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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

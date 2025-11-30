package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `138
583
868
822
783
65
262
121
508
780
461
484
668
389
808
215
97
500
30
915
856
400
444
623
781
786
3
713
457
273
739
822
235
606
968
105
924
326
32
23
27
666
555
10
962
903
391
703
222
993
433
744
30
541
228
783
449
962
508
567
239
354
237
694
225
780
471
976
297
949
23
427
858
939
570
945
658
103
191
645
742
881
304
124
761
341
918
739
997
729
513
959
991
433
520
850
933
687
195
311`

func solveCase(n int64) int64 {
	if n%2 == 0 {
		k := n / 2
		return (k + 1) * (k + 1)
	}
	return (n + 1) * (n + 3) / 2
}

func parseTestcases() ([]int64, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	nums := make([]int64, len(fields))
	for i := 0; i < len(fields); i++ {
		v, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse n%d: %v", i+1, err)
		}
		nums[i] = v
	}
	return nums, nil
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		input := fmt.Sprintf("%d\n", n)
		expected := strconv.FormatInt(solveCase(n), 10)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

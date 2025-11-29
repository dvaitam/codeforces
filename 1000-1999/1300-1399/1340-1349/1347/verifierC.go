package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesC.txt to avoid external files.
const testcasesRaw = `100
582
8604
2431
4209
9873
2555
6210
9552
4824
7708
1088
1387
8463
646
1088
3688
2139
666
4923
251
7350
5417
2633
2439
7549
6084
8274
6263
8681
8232
551
9403
1486
8487
9833
1252
6986
3377
4746
8773
9808
6846
7901
6367
9953
9608
3827
336
4
2982
4956
8305
9344
4170
5451
1076
8086
4293
4962
6686
6295
6287
1021
2684
2087
3915
4704
5473
910
589
7887
6847
2309
8059
9863
1338
2481
5780
6737
577
7640
6336
7519
771
1663
7716
2481
332
532
9803
2175
5307
1726
8998
5680
3195
6282
8034
1819
986`

type testCase struct {
	n int
}

// referenceSolution embeds logic from 1347C.go (decompose into round numbers).
func referenceSolution(n int) []int {
	res := make([]int, 0)
	power := 1
	for n > 0 {
		digit := n % 10
		if digit != 0 {
			res = append(res, digit*power)
		}
		n /= 10
		power *= 10
	}
	return res
}

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil
	}
	// first line is count
	count, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	tests := make([]testCase, 0, count)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			panic("invalid testcase value")
		}
		tests = append(tests, testCase{n: n})
	}
	return tests
}

func runExe(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := parseTestcases()

	for i, tc := range tests {
		parts := referenceSolution(tc.n)
		expected := fmt.Sprintf("%d", len(parts))
		if len(parts) > 0 {
			expected += "\n"
			for idx, v := range parts {
				if idx > 0 {
					expected += " "
				}
				expected += strconv.Itoa(v)
			}
		} else {
			expected += "\n"
		}
		input := fmt.Sprintf("1\n%d\n", tc.n)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("Test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

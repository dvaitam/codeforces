package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesA.txt to avoid external dependency.
const testcasesRaw = `100
6312
6891
664
4243
8377
7962
6635
4970
7809
5867
9559
3579
8269
2282
4618
2290
1554
4105
8726
9862
2408
5082
1619
1209
5410
7736
9172
1650
5797
7114
5181
3351
9053
7816
7254
8542
4268
1021
8990
231
1529
6535
19
8087
5459
3997
5329
1032
3131
9299
3633
3910
2335
8897
7340
1495
1319
5244
8323
8017
1787
4939
9032
4770
2045
8970
5452
8853
3330
9883
8966
9628
4713
7291
1502
9770
6307
5195
9432
3967
4757
3013
3103
3060
541
4261
7808
1132
1472
2134
2451
634
1315
8858
6411
8595
4516
8550
3859
3526`

type testCase struct {
	n int
}

// referenceSolution matches 1352A.go: decompose n into non-zero place values.
func referenceSolution(n int) []int {
	res := make([]int, 0)
	base := 1
	for n > 0 {
		d := n % 10
		if d != 0 {
			res = append(res, d*base)
		}
		n /= 10
		base *= 10
	}
	return res
}

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil
	}
	// first line is count; remaining lines are n values
	tests := make([]testCase, 0, len(lines)-1)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			panic("invalid testcase number")
		}
		tests = append(tests, testCase{n: n})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	out, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("cannot parse k: %v", err)
	}
	nums := make([]int, 0, len(fields)-1)
	for _, f := range fields[1:] {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number %q", f)
		}
		nums = append(nums, v)
	}
	exp := referenceSolution(n)
	if k != len(exp) || len(nums) != k {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(nums))
	}
	for i := 0; i < k; i++ {
		if nums[i] != exp[i] {
			return fmt.Errorf("mismatch at position %d: expected %d got %d", i, exp[i], nums[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases()
	for i, tc := range tests {
		if err := runCase(bin, tc.n); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
